package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/gochain-io/explorer/server/models"
	"go.uber.org/zap"

	"github.com/gochain-io/explorer/server/backend"
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
		cfg := zapdriver.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		logger, err := cfg.Build()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
			os.Exit(1)
		}

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
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start app: %v\n", err)
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
					checkParentForBlock(importer, block.Number().Int64(), 100)
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
func backfill(importer *backend.Backend, startFrom int64) {
	blockNumber := startFrom
	if startFrom <= 0 {
		blockNumber = getLatestBlockNumber(importer)
	}
	for {
		if (blockNumber % 1000) == 0 {
			importer.Lgr.Info("Checking block in backfill", zap.Int64("Block", blockNumber))
		}
		blocksFromDB := importer.GetBlockByNumber(blockNumber)
		if blocksFromDB == nil {
			importer.Lgr.Info("Backfilling the block:", zap.Int64("Block", blockNumber))
			block, err := importer.BlockByNumber(blockNumber)
			if block != nil {
				importer.ImportBlock(block)
				if err != nil {
					importer.Lgr.Fatal("importBlock - backfill", zap.Error(err))
				}
			}
		}
		checkParentForBlock(importer, blockNumber, 5)
		checkTransactionsConsistency(importer, blockNumber)
		if blockNumber > 0 {
			blockNumber = blockNumber - 1
		} else {
			blockNumber = getLatestBlockNumber(importer)
		}
	}
}

func checkParentForBlock(importer *backend.Backend, blockNumber int64, numBlocksToCheck int) {
	numBlocksToCheck--
	if blockNumber == 0 {
		return
	}
	if importer.NeedReloadBlock(blockNumber) {
		blockNumber--
		importer.Lgr.Info("Redownloading the block because it's corrupted or missing:", zap.Int64("Block", blockNumber))
		block, err := importer.BlockByNumber(blockNumber)
		if block != nil {
			importer.ImportBlock(block)
			if err != nil {
				importer.Lgr.Fatal("importBlock - checkParentForBlock", zap.Error(err))
			}
		}
		if err != nil {
			importer.Lgr.Info("BlockByNumber - checkParentForBlock", zap.Error(err))
			checkParentForBlock(importer, blockNumber+1, numBlocksToCheck)
		}
		if numBlocksToCheck > 0 && block != nil {
			checkParentForBlock(importer, block.Number().Int64(), numBlocksToCheck)
		}
	}
}

func checkTransactionsConsistency(importer *backend.Backend, blockNumber int64) {
	if !importer.TransactionsConsistent(blockNumber) {
		importer.Lgr.Info("Redownloading the block because number of transactions are wrong:", zap.Int64("Block", blockNumber))
		block, err := importer.BlockByNumber(blockNumber)
		if err != nil {
			importer.Lgr.Fatal("checkTransactionsConsistency", zap.Error(err))
		}
		if block != nil {
			importer.ImportBlock(block)
		}
	}
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
