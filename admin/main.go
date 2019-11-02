package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/blendle/zapdriver"
	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain/web3"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var (
	verbose bool
)

const (
	pkVarName      = "WEB3_PRIVATE_KEY"
	addrVarName    = "WEB3_ADDRESS"
	networkVarName = "WEB3_NETWORK"
	rpcURLVarName  = "WEB3_RPC_URL"
)

func main() {
	var netName, rpcUrl, mongoUrl, dbName string
	var testnet bool

	ctx, cancelFn := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigCh {
			cancelFn()
		}
	}()

	app := cli.NewApp()
	app.Name = "explorer admin"
	app.Version = "0.0.1"
	app.Usage = "explorer admin tool"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "testnet",
			Usage:       "Shorthand for '-network testnet'.",
			Destination: &testnet,
			Hidden:      false},
		cli.StringFlag{
			Name:        "rpc-url",
			Usage:       "The network RPC URL",
			Destination: &rpcUrl,
			EnvVar:      rpcURLVarName,
			Hidden:      false},
		cli.BoolFlag{
			Name:        "verbose",
			Usage:       "Enable verbose logging",
			Destination: &verbose,
			Hidden:      false},
		cli.StringFlag{
			Name:        "mongo-url, m",
			Value:       "127.0.0.1:27017",
			Usage:       "mongo connection url",
			EnvVar:      "MONGO_URL",
			Destination: &mongoUrl,
		},
		cli.StringFlag{
			Name:        "mongo-dbname, db",
			Value:       "blocks",
			Usage:       "mongo database name",
			EnvVar:      "MONGO_DBNAME",
			Destination: &dbName,
		},
	}
	var network web3.Network
	var backendInstance *backend.Backend
	app.Before = func(*cli.Context) error {
		cfg := zapdriver.NewProductionConfig()
		logger, err := cfg.Build()
		if err != nil {
			fatalExit(err)
		}
		defer logger.Sync()
		defer func() {
			if rerr := recover(); rerr != nil {
				fatalExit(fmt.Errorf("%+v", rerr))
			}
		}()
		network = getNetwork(netName, rpcUrl, testnet)
		backendInstance, err = backend.NewBackend(ctx, mongoUrl, network.URL, dbName, nil, nil, logger)
		if err != nil {
			fatalExit(err)
		}
		return nil
	}
	app.Commands = []cli.Command{
		{Name: "reload",
			Usage:   "Reloads block/transaction/contract",
			Aliases: []string{"r"},
			Subcommands: []cli.Command{
				{
					Name:    "block",
					Usage:   "Reloads a block",
					Aliases: []string{"b"},
					Action: func(c *cli.Context) {
						ReloadBlock(ctx, backendInstance, c.Args().First())
					},
				},
				{
					Name:    "transaction",
					Usage:   "Reloads a transaction",
					Aliases: []string{"tx"},
					Action: func(c *cli.Context) {
						ReloadTransaction(ctx, backendInstance, c.Args().First())
					},
				},
				{
					Name:    "contract",
					Usage:   "Reloads a contract and all tokens transactonis if any(ERC20)",
					Aliases: []string{"c"},
					Action: func(c *cli.Context) {
						ReloadContract(ctx, backendInstance, c.Args().First())
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Fatal error", zap.Error(err))
	}
	log.Println("Stopping")
}

func fatalExit(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	os.Exit(1)
}

// getNetwork resolves the rpcUrl from the user specified options, or quits if an illegal combination or value is found.
func getNetwork(name, rpcURL string, testnet bool) web3.Network {
	var network web3.Network
	if rpcURL != "" {
		if name != "" {
			fatalExit(fmt.Errorf("Cannot set both rpcURL %q and network %q", rpcURL, network))
		}
		if testnet {
			fatalExit(fmt.Errorf("Cannot set both rpcURL %q and testnet", rpcURL))
		}
		network.URL = rpcURL
		network.Unit = "GO"
	} else {
		if testnet {
			if name != "" {
				fatalExit(fmt.Errorf("Cannot set both network %q and testnet", name))
			}
			name = "testnet"
		} else if name == "" {
			name = "gochain"
		}
		var ok bool
		network, ok = web3.Networks[name]
		if !ok {
			fatalExit(fmt.Errorf("Unrecognized network %q", name))
		}
		if verbose {
			log.Printf("Network: %v", name)
		}
	}
	if verbose {
		log.Println("Network Info:", network)
	}
	return network
}

//ReloadBlock reloads a block - deletes from DB and reimport it
func ReloadBlock(ctx context.Context, backendInstance *backend.Backend, blockID string) {
	bnum, err := strconv.ParseInt(blockID, 10, 0)
	var block *models.Block
	if err != nil {
		fmt.Printf("Deleting block %s\n", blockID)
		err = backendInstance.DeleteBlockByHash(blockID)
		if err != nil {
			fatalExit(err)
		}
		fmt.Printf("Reimporting block %s\n", blockID)
		block, err = backendInstance.GetBlockByHash(ctx, blockID)
	} else {
		fmt.Printf("Deleting block %d\n", bnum)
		err = backendInstance.DeleteBlockByNumber(bnum)
		if err != nil {
			fatalExit(err)
		}
		fmt.Printf("Reimporting block %d\n", bnum)
		block, err = backendInstance.GetBlockByNumber(ctx, bnum)
	}
	if err != nil {
		fatalExit(fmt.Errorf("Failed to get block:%s, %v", blockID, err))
	} else if block == nil {
		fatalExit(fmt.Errorf("Block not found"))
	}
	log.Printf("Block:%v", block)
}

//ReloadTransaction removes a transaction, cleans DB and reimport TX
func ReloadTransaction(ctx context.Context, backendInstance *backend.Backend, txHash string) {
	tx, err := backendInstance.GetTransactionByHash(ctx, txHash)
	if err != nil || tx == nil {
		fatalExit(fmt.Errorf("Failed to get transaction:%s, %v", txHash, err))
	}
	fmt.Printf("Deleting block %d\n", tx.BlockNumber)
	err = backendInstance.DeleteBlockByNumber(tx.BlockNumber)
	if err != nil {
		fatalExit(err)
	}
	_, err = backendInstance.GetBlockByNumber(ctx, tx.BlockNumber)
	if err != nil {
		fatalExit(fmt.Errorf("Failed to get block:%d, %v", tx.BlockNumber, err))
	}
}

//ReloadContract removes contract from db, delete all tokens and transactions and reimports everything
func ReloadContract(context.Context, *backend.Backend, string) {

}
