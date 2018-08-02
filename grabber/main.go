package main

import (
	"context"
	"os"

	"math/big"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gochain-io/explorer/api/backend"
	"github.com/gochain-io/gochain/common"
	"github.com/gochain-io/gochain/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var rpcUrl string
	var mongoUrl string
	var loglevel string
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
	}

	app.Action = func(c *cli.Context) error {
		level, _ := zerolog.ParseLevel(loglevel)
		zerolog.SetGlobalLevel(level)
		importer := backend.NewBackend(mongoUrl, rpcUrl)
		go listener(rpcUrl, importer)
		go backfill(rpcUrl, importer)
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
			log.Info().Int64("Block", header.Number.Int64()).Msg("Gettting block in listener")
			if prevHeader != header.Number.String() {
				log.Info().Str("Listener is downloading the block:", header.Number.String()).Msg("Gettting block in listener")
				block, err := client.BlockByNumber(context.Background(), header.Number)
				importer.ImportBlock(block)
				if err != nil {
					log.Fatal().Err(err).Msg("listener")
				}
				checkParentForBlock(&client, importer, block.Number().Int64(), 5)
				prevHeader = header.Number.String()
			}
		}
	}
}

func backfill(url string, importer *backend.Backend) {
	client := getClient(url)
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("backfill - HeaderByNumber")
	}
	log.Info().Msg(header.Number.String())
	blockNumber := header.Number
	for {
		log.Info().Int64("Block", blockNumber.Int64()).Msg("Checking block in backfill")
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
		blockNumber = big.NewInt(0).Sub(blockNumber, big.NewInt(1))
	}
}

func checkParentForBlock(client *ethclient.Client, importer *backend.Backend, blockNumber int64, numBlocksToCheck int) {
	numBlocksToCheck--
	log.Info().Int64("Checking the block for it's parent:", blockNumber)
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
	log.Info().Int64("Checking a transaction consistency for the block :", blockNumber)
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

func updateAddresses(url string, importer *backend.Backend) {
	client := getClient(url)
	lastUpdatedAt := time.Unix(0, 0)
	for {
		addresses := importer.GetActiveAdresses(lastUpdatedAt)
		log.Info().Int("Addresses in db", len(addresses)).Time("lastUpdatedAt", lastUpdatedAt).Msg("updateAddresses")
		for _, address := range addresses {
			balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address.Address), nil)
			if err != nil {
				log.Fatal().Err(err).Msg("updateAddresses")
			}
			log.Info().Str("Balance of the address:", address.Address).Str("Balance", balance.String()).Msg("updateAddresses")
			importer.ImportAddress(address.Address, balance)
		}
		lastUpdatedAt = time.Now()
		time.Sleep(120 * time.Second) //sleep for 2 minutes
	}
}
