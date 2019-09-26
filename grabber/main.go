package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/explorer/server/models"

	"github.com/gochain-io/gochain/v3/common"

	"github.com/blendle/zapdriver"
	"github.com/urfave/cli"
)

func main() {
	var rpcUrl string
	var mongoUrl string
	var dbName string
	var startFrom int64
	var blockRangeLimit uint64
	var workersCount uint
	cfg := zapdriver.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	logger, err := cfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	app := cli.NewApp()
	app.Usage = "Grabber populates a mongo database with explorer data."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "rpc-url, u",
			Value:       "https://rpc.gochain.io",
			Usage:       "rpc api url",
			Destination: &rpcUrl,
		},
		cli.StringFlag{
			Name:        "mongo-url, m",
			Value:       "127.0.0.1:27017",
			Usage:       "mongo connection url",
			Destination: &mongoUrl,
		},
		cli.StringFlag{
			Name:        "mongo-dbname, db",
			Value:       "blocks",
			Usage:       "mongo database name",
			Destination: &dbName,
		},
		cli.Int64Flag{
			Name:        "start-from, s",
			Value:       0, //1365000
			Usage:       "refill from this block",
			Destination: &startFrom,
		},
		cli.Uint64Flag{
			Name:        "block-range-limit, b",
			Value:       10000,
			Usage:       "block range limit",
			Destination: &blockRangeLimit,
		},
		cli.UintFlag{
			Name:        "workers-amount, w",
			Value:       10,
			Usage:       "parallel workers amount",
			Destination: &workersCount,
		},
		cli.StringSliceFlag{
			Name:  "locked-accounts",
			Usage: "accounts with locked funds to exclude from rich list and circulating supply",
		},
	}

	app.Action = func(c *cli.Context) error {
		lockedAccounts := c.StringSlice("locked-accounts")
		for i, l := range lockedAccounts {
			if !common.IsHexAddress(l) {
				return fmt.Errorf("invalid hex address: %s", l)
			}
			// Ensure canonical form, since queries are case-sensitive.
			lockedAccounts[i] = common.HexToAddress(l).Hex()
		}
		importer := backend.NewBackend(mongoUrl, rpcUrl, dbName, lockedAccounts, nil, logger)
		go listener(importer)
		go updateStats(importer)
		go backfill(importer, startFrom)
		go updateAddresses(3*time.Minute, false, blockRangeLimit, workersCount, importer) // update only addresses
		updateAddresses(5*time.Second, true, blockRangeLimit, workersCount, importer)     // update contracts
		return nil
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Fatal("Failed to start app", zap.Error(err))
	}
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func listener(importer *backend.Backend) {
	var prevHeader int64
	ticker := time.NewTicker(time.Second * 1).C
	for {
		select {
		case <-ticker:
			latestBlocknumber := getLatestBlockNumber(importer)
			importer.Lgr.Debug("Getting block in listener", zap.Int64("Block Number", latestBlocknumber))
			if prevHeader != latestBlocknumber {
				importer.Lgr.Info("Getting block in listener", zap.Int64("Block Number", latestBlocknumber))
				block, err := importer.BlockByNumber(latestBlocknumber)
				if block != nil {
					importer.ImportBlock(block)
					if err != nil {
						importer.Lgr.Fatal("Listener", zap.Int64("Block Number", latestBlocknumber))
					}
					checkAncestors(importer, block.Number().Int64(), 100)
					prevHeader = latestBlocknumber
				}
			}
		}
	}
}

func getLatestBlockNumber(importer *backend.Backend) int64 {
	number, err := importer.GetFirstBlockNumber()
	if err != nil {
		importer.Lgr.Fatal("getLatestBlockNumber", zap.Error(err))
	}
	return number
}

// backfill continuously loops over all blocks in reverse order, verifying that all
// data is loaded and consistent, and reloading as necessary.
func backfill(importer *backend.Backend, startFrom int64) {
	blockNumber := startFrom
	if startFrom <= 0 {
		blockNumber = getLatestBlockNumber(importer)
	}
	var hash common.Hash // Expected hash.
	for {
		logger := importer.Lgr.With(zap.Int64("blockNumber", blockNumber))
		if (blockNumber % 1000) == 0 {
			logger.Info("Backfill progress")
		}
		var backfill bool
		dbBlock := importer.GetBlockByNumber(blockNumber)
		if dbBlock == nil {
			// Missing.
			backfill = true
			logger = logger.With(zap.String("reason", "missing"))
		} else if !common.EmptyHash(hash) && dbBlock.BlockHash != hash.Hex() {
			// Mismatch with expected hash from parent of previous block.
			backfill = true
			logger = logger.With(zap.String("reason", "hash"))
		} else {
			// Note parent as next expected hash.
			hash = common.HexToHash(dbBlock.ParentHash)
		}
		if backfill {
			logger.Info("Backfilling block")
			rpcBlock, err := importer.BlockByNumber(blockNumber)
			if rpcBlock != nil {
				importer.ImportBlock(rpcBlock)
				if err != nil {
					logger.Fatal("importBlock - backfill", zap.Error(err))
				}
				// Note parent as next expected hash.
				hash = rpcBlock.ParentHash()
			} else {
				hash = common.Hash{}
			}
		} else {
			newParent := checkTransactionsConsistency(importer, blockNumber)
			if !common.EmptyHash(newParent) {
				hash = newParent
			}
		}
		if blockNumber > 0 {
			blockNumber = blockNumber - 1
		} else {
			blockNumber = getLatestBlockNumber(importer)
			// No expected hash for latest.
			hash = common.Hash{}
		}
	}
}

// checkAncestors scans the ancestors of a block to verify that they exists and their hashes match.
func checkAncestors(importer *backend.Backend, blockNumber int64, generations int) {
	if blockNumber == 0 || generations <= 0 {
		return
	}
	oldest := blockNumber - int64(generations)
	if oldest < 0 {
		oldest = 0
	}
	for ; blockNumber > oldest && importer.NeedReloadParent(blockNumber); blockNumber-- {
		importer.Lgr.Info("Redownloading the block because it's corrupted or missing", zap.Int64("blockNumber", blockNumber))
		block, err := importer.BlockByNumber(blockNumber)
		if err != nil {
			importer.Lgr.Fatal("blockByNumber - checkAncestors", zap.Int64("blockNumber", blockNumber), zap.Error(err))
		}
		if block != nil {
			importer.ImportBlock(block)
			if err != nil {
				importer.Lgr.Fatal("importBlock - checkAncestors", zap.Error(err))
			}
		}
	}
}

// checkTransactionsConsistency checks if the block tx count matches the number of txs, and reimports the block if not.
// If a block is reimported, then the parent hash is returned.
func checkTransactionsConsistency(importer *backend.Backend, blockNumber int64) common.Hash {
	if !importer.TransactionsConsistent(blockNumber) {
		importer.Lgr.Info("Redownloading block because number of transactions are wrong", zap.Int64("blockNumber", blockNumber))
		block, err := importer.BlockByNumber(blockNumber)
		if err != nil {
			importer.Lgr.Fatal("checkTransactionsConsistency", zap.Int64("blockNumber", blockNumber), zap.Error(err))
		}
		if block != nil {
			importer.ImportBlock(block)
			return block.ParentHash()
		}
	}
	return common.Hash{}
}

func updateAddresses(sleep time.Duration, updateContracts bool, blockRangeLimit uint64, workersCount uint, importer *backend.Backend) {
	lastUpdatedAt := time.Unix(0, 0)
	lastBlockUpdatedAt := int64(0)
	for {
		start := time.Now()
		currentTime := time.Now()
		currentBlock := getLatestBlockNumber(importer)
		addresses := importer.GetActiveAdresses(lastUpdatedAt, updateContracts)
		importer.Lgr.Info("updateAddresses:", zap.Time("lastUpdatedAt", lastUpdatedAt), zap.Int("Addresses in db", len(addresses)))
		var jobs = make(chan *models.ActiveAddress, workersCount)
		var wg sync.WaitGroup
		for i := 0; i < int(workersCount); i++ {
			wg.Add(1)
			go worker(&wg, jobs, currentBlock, blockRangeLimit, importer)
		}
		for _, address := range addresses {
			jobs <- address
		}
		close(jobs)
		wg.Wait()
		elapsed := time.Since(start)
		importer.Lgr.Info("Performance measurement:", zap.Bool("updateContracts", updateContracts), zap.String("Updating all addresses took", elapsed.String()), zap.Int64("Current block", lastBlockUpdatedAt))
		lastBlockUpdatedAt = currentBlock
		lastUpdatedAt = currentTime
		time.Sleep(sleep)
	}
}

func worker(wg *sync.WaitGroup, jobs chan *models.ActiveAddress, currentBlock int64, blockRangeLimit uint64, importer *backend.Backend) {
	for address := range jobs {
		updateAddress(address, currentBlock, blockRangeLimit, importer)
	}
	wg.Done()
}
func updateAddress(address *models.ActiveAddress, currentBlock int64, blockRangeLimit uint64, importer *backend.Backend) {
	normalizedAddress := common.HexToAddress(address.Address).Hex()
	balance, err := importer.BalanceAt(normalizedAddress, "latest")
	if err != nil {
		importer.Lgr.Fatal("updateAddresses", zap.Error(err))
	}
	contractDataArray, err := importer.CodeAt(normalizedAddress)
	contractData := string(contractDataArray[:])
	var tokenDetails = &backend.TokenDetails{TotalSupply: big.NewInt(0)}
	contract := false
	if contractData != "" {
		contract = true
		byteCode := hex.EncodeToString(contractDataArray)
		importer.ImportContract(normalizedAddress, byteCode)
		tokenDetails, err = importer.GetTokenDetails(normalizedAddress, byteCode)
		if err != nil {
			importer.Lgr.Info("Cannot GetTokenDetails", zap.Error(err), zap.String("Address", normalizedAddress))
			// continue
		} else {
			var fromBlock int64
			contractFromDB, err := importer.GetAddressByHash(normalizedAddress)
			if err != nil {
				importer.Lgr.Fatal("updateAddresses", zap.Error(err))
			}
			if contractFromDB != nil && contractFromDB.UpdatedAtBlock > 0 {
				fromBlock = contractFromDB.UpdatedAtBlock
			} else {
				fromBlock = importer.GetContractBlock(normalizedAddress)
			}
			if contractFromDB.TokenName == "" || contractFromDB.TokenSymbol == "" {
				importer.Lgr.Info("TokenName and TokenSymbols are empty using from token details", zap.String("Address", normalizedAddress), zap.Int64("Contract block", fromBlock), zap.String("Name", tokenDetails.Name))
				contractFromDB.TokenName = tokenDetails.Name
				contractFromDB.TokenSymbol = tokenDetails.Symbol
			}
			internalTxs := importer.GetInternalTransactions(normalizedAddress, fromBlock, blockRangeLimit)
			internalTxsFromDb := importer.CountInternalTransactions(normalizedAddress)
			importer.Lgr.Info("Comparing number of internal txs in the db and in the gochain", zap.String("Address", normalizedAddress), zap.Int("In the gochain", len(internalTxs)), zap.Int("In the db", internalTxsFromDb))
			if len(internalTxs) != internalTxsFromDb {
				var tokenHoldersList []string
				for _, itx := range internalTxs {
					importer.Lgr.Debug("Internal Transaction", zap.String("From", itx.From.String()), zap.String("To", itx.To.String()), zap.Int64("Value", itx.Value.Int64()))
					importer.ImportInternalTransaction(normalizedAddress, itx)
					// if itx.BlockNumber > lastBlockUpdatedAt.Int64() {
					importer.Lgr.Debug("Updating following token holder addresses", zap.String("addr 1", itx.From.String()), zap.String("addr 2", itx.To.String()), zap.Int64("Value", itx.Value.Int64()))
					tokenHoldersList = appendIfMissing(tokenHoldersList, itx.To.String())
					tokenHoldersList = appendIfMissing(tokenHoldersList, itx.From.String())
					// }
				}
				for index, tokenHolderAddress := range tokenHoldersList {
					if tokenHolderAddress == "0x0000000000000000000000000000000000000000" {
						continue
					}
					importer.Lgr.Info("Importing token holder", zap.Int("Index", index), zap.Int("Total number", len(tokenHoldersList)))
					tokenHolder, err := importer.GetTokenBalance(normalizedAddress, tokenHolderAddress)
					if err != nil {
						importer.Lgr.Info("Cannot GetTokenBalance, in internal transaction", zap.Error(err), zap.String("Address", tokenHolderAddress))
						continue
					}
					if contractFromDB == nil {
						importer.Lgr.Info("Cannot find contract in DB", zap.Error(err), zap.String("Address", tokenHolderAddress))
						continue
					}
					importer.ImportTokenHolder(normalizedAddress, tokenHolderAddress, tokenHolder, contractFromDB)
				}
			}
		}
	}
	importer.Lgr.Info("updateAddresses", zap.String("Balance of the address:", normalizedAddress), zap.String("Balance", balance.String()))
	importer.ImportAddress(normalizedAddress, balance, tokenDetails, contract, currentBlock)
}
func updateStats(importer *backend.Backend) {
	for {
		importer.Lgr.Info("Updating stats")
		importer.UpdateStats()
		importer.Lgr.Info("Updating stats finished")
		time.Sleep(300 * time.Second) //sleep for 5 minutes
	}
}
