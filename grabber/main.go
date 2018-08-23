package main

import (
	"context"
	"os"

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
		// go listener(rpcUrl, importer)
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

func backfill(url string, importer *backend.Backend, startFrom int64) {
	client := getClient(url)
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("backfill - HeaderByNumber")
	}
	log.Info().Msg(header.Number.String())
	blockNumber := header.Number
	if startFrom > 0 {
		blockNumber = big.NewInt(startFrom)
	}
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
			contractDataArray, err := client.CodeAt(context.Background(), common.HexToAddress(address.Address), nil)
			contractData := string(contractDataArray[:])
			var tokenName, tokenSymbol string
			go20 := false
			contract := false
			if contractData != "" {
				go20 = true
				contract = true
				txs := importer.GetTransactionList(address.Address)
				for _, tx := range txs {
					res, err := importer.GetTokenBalance(address.Address, tx.From)
					if err != nil {
						log.Info().Err(err).Msg("Cannot GetTokenBalance, seems like not ERC20 compatible contract")
						go20 = false
						continue
					}
					tokenName = res.Name
					tokenSymbol = res.Symbol
					importer.ImportTokenHolder(address.Address, tx.From, res.Balance, tokenName, tokenSymbol)
					log.Info().Str("Balance", res.Balance.String()).Str("Transaction", tx.From).Msg("Contract data is not empty")
				}
			}
			log.Info().Str("Balance of the address:", address.Address).Str("Balance", balance.String()).Str("Contract data", contractData).Msg("updateAddresses")
			importer.ImportAddress(address.Address, balance, tokenName, tokenSymbol, contract, go20)
		}
		lastUpdatedAt = time.Now()
		time.Sleep(120 * time.Second) //sleep for 2 minutes
	}
}
