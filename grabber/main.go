package main

import (
	"context"
	"os"
	"strings"

	"math/big"
	"time"

	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/gochain/common"
	"github.com/gochain-io/gochain/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
)

func main() {
	var rpcUrl string
	var mongoUrl string
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
		importer := backend.NewBackend(mongoUrl, rpcUrl)
		go listener(rpcUrl, importer)
		go backfill(rpcUrl, importer, startFrom)
		updateAddresses(rpcUrl, importer)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Run")
	}
}

func getClient(url string) ethclient.Client {
	client, err := ethclient.Dial(url)
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

func getFirstBlockNumber(client ethclient.Client) *big.Int {
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

func checkParentForBlock(client *ethclient.Client, importer *backend.Backend, blockNumber int64, numBlocksToCheck int) {
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

func checkTransactionsConsistency(client *ethclient.Client, importer *backend.Backend, blockNumber int64) {
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
func updateAddresses(url string, importer *backend.Backend) {
	client := getClient(url)
	lastUpdatedAt := time.Unix(0, 0)
	_, genesisAddressList, err := importer.GenesisAlloc()
	if err != nil {
		log.Fatal().Err(err).Msg("failed response from GenesisAlloc")
	}
	log.Info().Str("Genesis addresses", strings.Join(genesisAddressList[:], ",")).Msg("updateAddresses")
	for {
		addresses := importer.GetActiveAdresses(lastUpdatedAt)
		log.Info().Int("Addresses in db", len(addresses)).Time("lastUpdatedAt", lastUpdatedAt).Msg("updateAddresses")
		for _, address := range addresses {
			if stringInSlice(address.Address, genesisAddressList) {
				log.Info().Str("Following address is in the list of genesis addresses", address.Address).Msg("updateAddresses")
				continue
			}
			balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address.Address), nil)
			if err != nil {
				log.Fatal().Err(err).Msg("updateAddresses")
			}
			contractDataArray, err := client.CodeAt(context.Background(), common.HexToAddress(address.Address), nil)
			contractData := string(contractDataArray[:])
			var tokenName, tokenSymbol string
			go20 := false
			contract := false
			if contractData != "" {
				go20 = true
				contract = true
				internalTxs := importer.GetInternalTransactions(address.Address)
				for _, itx := range internalTxs {
					log.Info().Str("From", itx.From.String()).Str("To", itx.To.String()).Int64("Value", itx.Value.Int64()).Msg("Internal Transaction")
					importer.ImportInternalTransaction(address.Address, itx)
					res, err := importer.GetTokenBalance(address.Address, itx.To.String())
					tokenName = res.Name
					tokenSymbol = res.Symbol
					if err != nil {
						log.Info().Err(err).Str("Address", itx.To.String()).Msg("Cannot GetTokenBalance, in internal transaction")
						go20 = false
						continue
					}
					importer.ImportTokenHolder(address.Address, itx.To.String(), res.Balance, tokenName, tokenSymbol)
					res, err = importer.GetTokenBalance(address.Address, itx.From.String())
					if err != nil {
						log.Info().Err(err).Str("Address", itx.From.String()).Msg("Cannot GetTokenBalance, in internal transaction")
						go20 = false
						continue
					}
					importer.ImportTokenHolder(address.Address, itx.From.String(), res.Balance, tokenName, tokenSymbol)
				}
			}
			log.Info().Str("Balance of the address:", address.Address).Str("Balance", balance.String()).Msg("updateAddresses")
			importer.ImportAddress(address.Address, balance, tokenName, tokenSymbol, contract, go20)
		}
		lastUpdatedAt = time.Now()
		time.Sleep(300 * time.Second) //sleep for 5 minutes
	}
}
