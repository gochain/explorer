package main

import (
	"os"

	"math/big"
	"time"

	"encoding/hex"

	"github.com/gochain-io/explorer/server/backend"
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
	}

	app.Action = func(c *cli.Context) error {
		level, _ := zerolog.ParseLevel(loglevel)
		zerolog.SetGlobalLevel(level)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		importer := backend.NewBackend(mongoUrl, rpcUrl, dbName)
		go listener(rpcUrl, importer)
		go updateStats(importer)
		go backfill(rpcUrl, importer, startFrom)
		go updateAddresses(rpcUrl, false, importer) // update only addresses
		updateAddresses(rpcUrl, true, importer)     // update contracts
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

func listener(url string, importer *backend.Backend) {
	var prevHeader int64
	ticker := time.NewTicker(time.Second * 1).C
	for {
		select {
		case <-ticker:
			latestBlocknumber := getFirstBlockNumber(importer)
			log.Debug().Int64("Block", latestBlocknumber).Msg("Gettting block in listener")
			if prevHeader != latestBlocknumber {
				log.Info().Int64("Listener is downloading the block:", latestBlocknumber).Msg("Gettting block in listener")
				block, err := importer.BlockByNumber(latestBlocknumber)
				if block != nil {
					importer.ImportBlock(block)
					if err != nil {
						log.Fatal().Err(err).Msg("listener")
					}
					checkParentForBlock(importer, block.Number().Int64(), 100)
					prevHeader = latestBlocknumber
				}
			}
		}
	}
}

func getFirstBlockNumber(importer *backend.Backend) int64 {
	number, err := importer.GetFirstBlockNumber()
	if err != nil {
		log.Fatal().Err(err).Msg("getFirstBlockNumber")
	}
	return number
}
func backfill(url string, importer *backend.Backend, startFrom int64) {
	blockNumber := getFirstBlockNumber(importer)
	if startFrom > 0 {
		blockNumber = startFrom
	}
	for {
		if (blockNumber % 1000) == 0 {
			log.Info().Int64("Block", blockNumber).Msg("Checking block in backfill")
		}
		blocksFromDB := importer.GetBlockByNumber(blockNumber)
		if blocksFromDB == nil {
			log.Info().Int64("Backfilling the block:", blockNumber).Msg("Gettting block in backfill")
			block, err := importer.BlockByNumber(blockNumber)
			if block != nil {
				importer.ImportBlock(block)
				if err != nil {
					log.Fatal().Err(err).Msg("importBlock - backfill")
				}
			}
		}
		checkParentForBlock(importer, blockNumber, 5)
		checkTransactionsConsistency(importer, blockNumber)
		if blockNumber > 0 {
			blockNumber = blockNumber - 1
		} else {
			blockNumber = getFirstBlockNumber(importer)
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
		log.Info().Int64("Redownloading the block because it's corrupted or missing:", blockNumber).Msg("checkParentForBlock")
		block, err := importer.BlockByNumber(blockNumber)
		if block != nil {
			importer.ImportBlock(block)
			if err != nil {
				log.Fatal().Err(err).Msg("importBlock - checkParentForBlock")
			}
		}
		if err != nil {
			log.Info().Err(err).Msg("BlockByNumber - checkParentForBlock")
			checkParentForBlock(importer, blockNumber+1, numBlocksToCheck)
		}
		if numBlocksToCheck > 0 && block != nil {
			checkParentForBlock(importer, block.Number().Int64(), numBlocksToCheck)
		}
	}
}

func checkTransactionsConsistency(importer *backend.Backend, blockNumber int64) {
	if !importer.TransactionsConsistent(blockNumber) {
		log.Info().Int64("Redownloading the block because number of transactions are wrong", blockNumber).Msg("checkTransactionsConsistency")
		block, err := importer.BlockByNumber(blockNumber)
		if err != nil {
			log.Fatal().Err(err).Msg("checkTransactionsConsistency")
		}
		if block != nil {
			importer.ImportBlock(block)
		}
	}
}

func updateAddresses(url string, updateContracts bool, importer *backend.Backend) {
	lastUpdatedAt := time.Unix(0, 0)
	lastBlockUpdatedAt := int64(0)
	for {
		start := time.Now()
		currentTime := time.Now()
		currentBlock := getFirstBlockNumber(importer)
		addresses := importer.GetActiveAdresses(lastUpdatedAt, updateContracts)
		log.Info().Int("Addresses in db", len(addresses)).Time("lastUpdatedAt", lastUpdatedAt).Msg("updateAddresses")
		for index, address := range addresses {
			normalizedAddress := common.HexToAddress(address.Address).Hex()
			balance, err := importer.BalanceAt(normalizedAddress, "pending")
			if err != nil {
				log.Fatal().Err(err).Msg("updateAddresses")
			}
			contractDataArray, err := importer.CodeAt(normalizedAddress)
			contractData := string(contractDataArray[:])
			go20 := false
			var tokenDetails = &backend.TokenDetails{TotalSupply: big.NewInt(0)}
			contract := false
			if contractData != "" {
				contract = true
				byteCode := hex.EncodeToString(contractDataArray)
				importer.ImportContract(normalizedAddress, byteCode)
				tokenDetails, err = importer.GetTokenDetails(normalizedAddress, byteCode)
				if err != nil {
					log.Info().Err(err).Str("Address", normalizedAddress).Msg("Cannot GetTokenDetails")
					go20 = false
					// continue
				} else {
					go20 = true
					var fromBlock int64
					contractFromDB := importer.GetAddressByHash(normalizedAddress)
					if contractFromDB != nil && contractFromDB.UpdatedAtBlock > 0 {
						fromBlock = contractFromDB.UpdatedAtBlock
					} else {
						fromBlock = importer.GetContractBlock(normalizedAddress)
					}
					internalTxs := importer.GetInternalTransactions(normalizedAddress, fromBlock)
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
							tokenHolder, err := importer.GetTokenBalance(normalizedAddress, tokenHolderAddress)
							log.Debug().Int("Index", index).Int("Total number", len(tokenHoldersList)).Msg("Importing token holder")
							if err != nil {
								log.Info().Err(err).Str("Address", tokenHolderAddress).Msg("Cannot GetTokenBalance, in internal transaction")
								go20 = false
								continue
							}
							importer.ImportTokenHolder(normalizedAddress, tokenHolderAddress, tokenHolder)
						}
					}
				}
			}
			log.Info().Str("Balance of the address:", normalizedAddress).Int("Index", index).Int("Total number", len(addresses)).Str("Balance", balance.String()).Msg("updateAddresses")
			importer.ImportAddress(normalizedAddress, balance, tokenDetails, contract, go20, currentBlock)
		}
		elapsed := time.Since(start)
		log.Info().Bool("updateContracts", updateContracts).Str("Updating all addresses took", elapsed.String()).Int64("Current block", lastBlockUpdatedAt).Msg("Performance measurement")
		lastBlockUpdatedAt = currentBlock
		lastUpdatedAt = currentTime
		time.Sleep(180 * time.Second) //sleep for 3 minutes
	}
}
func updateStats(importer *backend.Backend) {
	for {
		log.Info().Msg("Updating stats")
		importer.UpdateStats()
		log.Info().Msg("Updating stats finished")
		time.Sleep(300 * time.Second) //sleep for 5 minutes
	}
}
