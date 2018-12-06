package main

import (
	"context"
	"os"
	"strings"

	"math/big"
	"time"

	"encoding/hex"

	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/gochain/common"
	"github.com/gochain-io/gochain/goclient"
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

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "rpc-url, u",
			Value:       "https://rpc.gochain.io",
			Usage:       "rpc api url, 'https://rpc.gochain.io'",
			Destination: &rpcUrl,
		},
		cli.StringFlag{
			Name:        "mongo-url, m",
			Value:       "127.0.0.1:27017",
			Usage:       "mongo connection url, '127.0.0.1:27017'",
			Destination: &mongoUrl,
		},
		cli.StringFlag{
			Name:        "mongo-dbname, db",
			Value:       "blocks",
			Usage:       "mongo database name, 'blocks'",
			Destination: &dbName,
		},
		cli.StringFlag{
			Name:        "log, l",
			Value:       "info",
			Usage:       "loglevel debug/info/warn/fatal, default is Info",
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
		go updateAddresses(rpcUrl, true, importer) // update contracts
		updateAddresses(rpcUrl, false, importer)   // update only addresses
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

func getClient(url string) goclient.Client {
	client, err := goclient.Dial(url)
	if err != nil {
		log.Fatal().Err(err).Msg("main")
	}
	return *client
}
func listener(url string, importer *backend.Backend) {
	client := getClient(url)
	var prevHeader string
	ticker := time.NewTicker(time.Second * 1).C
	for {
		select {
		case <-ticker:
			header, err := client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Fatal().Err(err).Msg("HeaderByNumber")
			}
			log.Debug().Int64("Block", header.Number.Int64()).Msg("Gettting block in listener")
			if prevHeader != header.Number.String() {
				log.Info().Str("Listener is downloading the block:", header.Number.String()).Msg("Gettting block in listener")
				block, err := client.BlockByNumber(context.Background(), header.Number)
				if block != nil {
					importer.ImportBlock(block)
					if err != nil {
						log.Fatal().Err(err).Msg("listener")
					}
					checkParentForBlock(&client, importer, block.Number().Int64(), 100)
					prevHeader = header.Number.String()
				}
			}
		}
	}
}

func getFirstBlockNumber(client goclient.Client) *big.Int {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("backfill - HeaderByNumber")
	}
	log.Info().Msg(header.Number.String())
	return header.Number
}
func backfill(url string, importer *backend.Backend, startFrom int64) {
	client := getClient(url)
	blockNumber := getFirstBlockNumber(client)
	if startFrom > 0 {
		blockNumber = big.NewInt(startFrom)
	}
	for {
		if (blockNumber.Int64() % 1000) == 0 {
			log.Info().Int64("Block", blockNumber.Int64()).Msg("Checking block in backfill")
		}
		blocksFromDB := importer.GetBlockByNumber(blockNumber.Int64())
		if blocksFromDB == nil {
			log.Info().Str("Backfilling the block:", blockNumber.String()).Msg("Gettting block in backfill")
			block, err := client.BlockByNumber(context.Background(), blockNumber)
			if block != nil {
				importer.ImportBlock(block)
				if err != nil {
					log.Fatal().Err(err).Msg("importBlock - backfill")
				}
			}
		}
		checkParentForBlock(&client, importer, blockNumber.Int64(), 5)
		checkTransactionsConsistency(&client, importer, blockNumber.Int64())
		if blockNumber.Int64() > 0 {
			blockNumber = big.NewInt(0).Sub(blockNumber, big.NewInt(1))
		} else {
			blockNumber = getFirstBlockNumber(client)
		}
	}
}

func checkParentForBlock(client *goclient.Client, importer *backend.Backend, blockNumber int64, numBlocksToCheck int) {
	numBlocksToCheck--
	if blockNumber == 0 {
		return
	}
	if importer.NeedReloadBlock(blockNumber) {
		blockNumber--
		log.Info().Int64("Redownloading the block because it's corrupted or missing:", blockNumber).Msg("checkParentForBlock")
		block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
		if block != nil {
			importer.ImportBlock(block)
			if err != nil {
				log.Fatal().Err(err).Msg("importBlock - checkParentForBlock")
			}
		}
		if err != nil {
			log.Info().Err(err).Msg("BlockByNumber - checkParentForBlock")
			checkParentForBlock(client, importer, blockNumber+1, numBlocksToCheck)
		}
		if numBlocksToCheck > 0 && block != nil {
			checkParentForBlock(client, importer, block.Number().Int64(), numBlocksToCheck)
		}
	}
}

func checkTransactionsConsistency(client *goclient.Client, importer *backend.Backend, blockNumber int64) {
	if !importer.TransactionsConsistent(blockNumber) {
		log.Info().Int64("Redownloading the block because number of transactions are wrong", blockNumber).Msg("checkTransactionsConsistency")
		block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
		if block != nil {
			importer.ImportBlock(block)
			if err != nil {
				log.Fatal().Err(err).Msg("importBlock - checkParentForBlock")
			}
		}
		if err != nil {
			log.Fatal().Err(err).Msg("BlockByNumber - checkParentForBlock")
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}
func updateAddresses(url string, updateContracts bool, importer *backend.Backend) {
	client := getClient(url)
	lastUpdatedAt := time.Unix(0, 0)
	lastBlockUpdatedAt := big.NewInt(0)
	_, genesisAddressList, err := importer.GenesisAlloc()
	if err != nil {
		log.Fatal().Err(err).Msg("failed response from GenesisAlloc")
	}
	log.Info().Str("Genesis addresses", strings.Join(genesisAddressList[:], ",")).Msg("updateAddresses")
	for {
		start := time.Now()
		currentTime := time.Now()
		currentBlock := getFirstBlockNumber(client)
		addresses := importer.GetActiveAdresses(lastUpdatedAt, updateContracts)
		log.Info().Int("Addresses in db", len(addresses)).Time("lastUpdatedAt", lastUpdatedAt).Msg("updateAddresses")
		for index, address := range addresses {
			normalizedAddress := common.HexToAddress(address.Address).Hex()
			if stringInSlice(normalizedAddress, genesisAddressList) {
				log.Info().Str("Following address is in the list of genesis addresses", normalizedAddress).Msg("updateAddresses")
				continue
			}
			balance, err := client.BalanceAt(context.Background(), common.HexToAddress(normalizedAddress), nil)
			if err != nil {
				log.Fatal().Err(err).Msg("updateAddresses")
			}
			contractDataArray, err := client.CodeAt(context.Background(), common.HexToAddress(normalizedAddress), nil)
			contractData := string(contractDataArray[:])
			go20 := false
			var tokenDetails = &backend.TokenDetails{TotalSupply: big.NewInt(0)}
			contract := false
			if contractData != "" {
				contract = true
				byteCode := hex.EncodeToString(contractDataArray)
				importer.ImportContract(normalizedAddress, byteCode)
				tokenDetails, err = importer.GetTokenDetails(normalizedAddress)
				if err != nil {
					log.Info().Err(err).Str("Address", normalizedAddress).Msg("Cannot GetTokenDetails")
					go20 = false
					// continue
				} else {
					go20 = true
					internalTxs := importer.GetInternalTransactions(normalizedAddress)
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
			log.Info().Str("Balance of the address:", normalizedAddress).Int("Index", index).Int("Total number", len(addresses)).Str("Balance", balance.String()).Msg("updateAddresses")
			importer.ImportAddress(normalizedAddress, balance, tokenDetails, contract, go20)
		}
		elapsed := time.Since(start)
		log.Info().Bool("updateContracts", updateContracts).Str("Updating all addresses took", elapsed.String()).Int64("Current block", lastBlockUpdatedAt.Int64()).Msg("Performance measurement")
		lastBlockUpdatedAt = currentBlock
		lastUpdatedAt = currentTime
		time.Sleep(300 * time.Second) //sleep for 5 minutes
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
