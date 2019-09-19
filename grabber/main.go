package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/explorer/server/models"

	"github.com/gochain-io/gochain/v3/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
)

func main() {
	var rpcUrl string
	var mongoUrl string
	var dbName string
	var loglevel string
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
		cli.StringFlag{
			Name:        "log, l",
			Value:       "info",
			Usage:       "loglevel debug/info/warn/fatal",
			Destination: &loglevel,
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
		level, _ := zerolog.ParseLevel(loglevel)
		zerolog.SetGlobalLevel(level)
		lockedAccounts := c.StringSlice("locked-accounts")
		for i, l := range lockedAccounts {
			if !common.IsHexAddress(l) {
				return fmt.Errorf("invalid hex address: %s", l)
			}
			// Ensure canonical form, since queries are case-sensitive.
			lockedAccounts[i] = common.HexToAddress(l).Hex()
		}
		importer := backend.NewBackend(mongoUrl, rpcUrl, dbName, lockedAccounts, nil)
		go listener(importer)
		go updateStats(importer)
		go backfill(importer, startFrom)
		go updateAddresses(3*time.Minute, false, blockRangeLimit, workersCount, importer) // update only addresses
		updateAddresses(5*time.Second, true, blockRangeLimit, workersCount, importer)     // update contracts
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Run")
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
			log.Debug().Int64("Block", latestBlocknumber).Msg("Getting block in listener")
			if prevHeader != latestBlocknumber {
				log.Info().Int64("Listener is downloading the block:", latestBlocknumber).Msg("Getting block in listener")
				block, err := importer.BlockByNumber(latestBlocknumber)
				if block != nil {
					importer.ImportBlock(block)
					if err != nil {
						log.Fatal().Err(err).Msg("listener")
					}
					checkAncestors(importer, block.Number().Int64(), 100)
					prevHeader = latestBlocknumber
				}
			}
		}
	}
}

func getLatestBlockNumber(importer *backend.Backend) int64 {
	number, err := importer.GetLatestBlockNumber()
	if err != nil {
		log.Fatal().Err(err).Msg("getLatestBlockNumber")
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
		logger := log.With().Int64("blockNumber", blockNumber).Logger()
		if (blockNumber % 1000) == 0 {
			logger.Info().Msg("Backfill progress")
		}
		var backfill bool
		dbBlock := importer.GetBlockByNumber(blockNumber)
		if dbBlock == nil {
			// Missing.
			backfill = true
			logger = logger.With().Str("reason", "missing").Logger()
		} else if !common.EmptyHash(hash) && dbBlock.BlockHash != hash.Hex() {
			// Mismatch with expected hash from parent of previous block.
			backfill = true
			logger = logger.With().Str("reason", "hash").Logger()
		} else {
			// Note parent as next expected hash.
			hash = common.HexToHash(dbBlock.ParentHash)
		}
		if backfill {
			logger.Info().Msg("Backfilling block")
			rpcBlock, err := importer.BlockByNumber(blockNumber)
			if rpcBlock != nil {
				importer.ImportBlock(rpcBlock)
				if err != nil {
					logger.Fatal().Err(err).Msg("importBlock - backfill")
				}
				// Note parent as next expected hash.
				hash = rpcBlock.ParentHash()
			} else {
				hash = common.Hash{}
			}
		} else {
			hash = checkTransactionsConsistency(importer, blockNumber)
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
func checkAncestors(importer *backend.Backend, blockNumber int64, ancestorsToCheck int) {
	if blockNumber == 0 {
		return
	}
	oldest := blockNumber + int64(ancestorsToCheck)
	for ; blockNumber > oldest && importer.NeedReloadParent(blockNumber); blockNumber-- {
		log.Info().Int64("blockNumber", blockNumber).
			Msg("Redownloading the block because it's corrupted or missing")
		block, err := importer.BlockByNumber(blockNumber)
		if err != nil {
			log.Fatal().Int64("blockNumber", blockNumber).Err(err).Msg("blockByNumber - checkAncestors")
		}
		if block != nil {
			importer.ImportBlock(block)
			if err != nil {
				log.Fatal().Err(err).Msg("importBlock - checkAncestors")
			}
		}
	}
}

func checkTransactionsConsistency(importer *backend.Backend, blockNumber int64) common.Hash {
	if !importer.TransactionsConsistent(blockNumber) {
		log.Info().Int64("blockNumber", blockNumber).
			Msg("Redownloading block because number of transactions are wrong")
		block, err := importer.BlockByNumber(blockNumber)
		if err != nil {
			log.Fatal().Int64("blockNumber", blockNumber).Err(err).Msg("checkTransactionsConsistency")
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
		log.Info().Int("Addresses in db", len(addresses)).Time("lastUpdatedAt", lastUpdatedAt).Msg("updateAddresses")
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
		log.Info().Bool("updateContracts", updateContracts).Str("Updating all addresses took", elapsed.String()).Int64("Current block", lastBlockUpdatedAt).Msg("Performance measurement")
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
		log.Fatal().Err(err).Msg("updateAddresses")
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
			log.Info().Err(err).Str("Address", normalizedAddress).Msg("Cannot GetTokenDetails")
			// continue
		} else {
			var fromBlock int64
			contractFromDB, err := importer.GetAddressByHash(normalizedAddress)
			if err != nil {
				log.Fatal().Err(err).Msg("updateAddresses")
			}
			if contractFromDB != nil && contractFromDB.UpdatedAtBlock > 0 {
				fromBlock = contractFromDB.UpdatedAtBlock
			} else {
				fromBlock = importer.GetContractBlock(normalizedAddress)
			}
			if contractFromDB.TokenName == "" || contractFromDB.TokenSymbol == "" {
				log.Info().Str("Address", normalizedAddress).Int64("Contract block", fromBlock).Str("Name", tokenDetails.Name).Msg("TokenName and TokenSymbols are empty using from token details")
				contractFromDB.TokenName = tokenDetails.Name
				contractFromDB.TokenSymbol = tokenDetails.Symbol
			}
			internalTxs := importer.GetInternalTransactions(normalizedAddress, fromBlock, blockRangeLimit)
			internalTxsFromDb := importer.CountInternalTransactions(normalizedAddress)
			log.Info().Str("Address", normalizedAddress).Int64("Contract block", fromBlock).Int("In the gochain", len(internalTxs)).Int("In the db", internalTxsFromDb).Msg("Comparing number of internal txs in the db and in the gochain")
			if len(internalTxs) != internalTxsFromDb {
				var tokenHoldersList []string
				for _, itx := range internalTxs {
					log.Debug().Str("From", itx.From.String()).Str("To", itx.To.String()).Int64("Value", itx.Value.Int64()).Msg("Internal Transaction")
					importer.ImportInternalTransaction(normalizedAddress, itx)
					// if itx.BlockNumber > lastBlockUpdatedAt.Int64() {
					log.Debug().Str("addr 1", itx.From.String()).Str("addr 2", itx.To.String()).Int64("Value", itx.Value.Int64()).Msg("Updating following token holder addresses")
					tokenHoldersList = appendIfMissing(tokenHoldersList, itx.To.String())
					tokenHoldersList = appendIfMissing(tokenHoldersList, itx.From.String())
					// }
				}
				for index, tokenHolderAddress := range tokenHoldersList {
					if tokenHolderAddress == "0x0000000000000000000000000000000000000000" {
						continue
					}
					log.Info().Int("Index", index).Int("Total number", len(tokenHoldersList)).Msg("Importing token holder")
					tokenHolder, err := importer.GetTokenBalance(normalizedAddress, tokenHolderAddress)
					if err != nil {
						log.Info().Err(err).Str("Address", tokenHolderAddress).Msg("Cannot GetTokenBalance, in internal transaction")
						continue
					}
					if contractFromDB == nil {
						log.Info().Err(err).Str("Address", tokenHolderAddress).Msg("Cannot find contract in DB")
						continue
					}
					importer.ImportTokenHolder(normalizedAddress, tokenHolderAddress, tokenHolder, contractFromDB)
				}
			}
		}
	}
	log.Info().Str("Balance of the address:", normalizedAddress).Str("Balance", balance.String()).Msg("updateAddresses")
	importer.ImportAddress(normalizedAddress, balance, tokenDetails, contract, currentBlock)
}
func updateStats(importer *backend.Backend) {
	for {
		log.Info().Msg("Updating stats")
		importer.UpdateStats()
		log.Info().Msg("Updating stats finished")
		time.Sleep(300 * time.Second) //sleep for 5 minutes
	}
}
