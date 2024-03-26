package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/blendle/zapdriver"
	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain-io/explorer/server/tokens"
	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/web3"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

const (
	pkVarName       = "WEB3_PRIVATE_KEY"
	addrVarName     = "WEB3_ADDRESS"
	networkVarName  = "WEB3_NETWORK"
	rpcURLVarName   = "WEB3_RPC_URL"
	blockRangeLimit = 10000
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
	var backendInstance *backend.Backend
	app.Before = func(*cli.Context) error {
		var err error
		cfg := zapdriver.NewDevelopmentConfig()
		logger, err = cfg.Build()
		if err != nil {
			fatalExit(err)
		}
		defer logger.Sync()
		defer func() {
			if rerr := recover(); rerr != nil {
				fatalExit(fmt.Errorf("%+v", rerr))
			}
		}()
		network := getNetwork(netName, rpcUrl, testnet)
		backendInstance, err = backend.NewBackend(ctx, mongoUrl, network.URL, dbName, nil, nil, nil, logger, nil)
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
		fatalExit(err)
	}
	logger.Info("Done")
}

func fatalExit(err error) {
	logger.Error("ERROR:", zap.Error(err))
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
	}
	logger.Info("Network", zap.String("network", name))
	return network
}

//ReloadBlock reloads a block - deletes from DB and reimport it
func ReloadBlock(ctx context.Context, backendInstance *backend.Backend, blockID string) {
	bnum, err := strconv.ParseInt(blockID, 10, 0)
	var block *models.Block
	if err != nil {
		logger.Info("Deleting block ", zap.String("number", blockID))
		err = backendInstance.DeleteBlockByHash(blockID)
		if err != nil {
			fatalExit(err)
		}
		logger.Info("Reimporting block ", zap.String("number", blockID))
		block, err = backendInstance.GetBlockByHash(ctx, blockID)
	} else {
		logger.Info("Deleting block ", zap.String("number", blockID))
		err = backendInstance.DeleteBlockByNumber(bnum)
		if err != nil {
			fatalExit(err)
		}
		logger.Info("Reimporting block ", zap.String("number", blockID))
		block, err = backendInstance.GetBlockByNumber(ctx, bnum, false)
	}
	if err != nil {
		fatalExit(fmt.Errorf("Failed to get block:%s, %v", blockID, err))
	} else if block == nil {
		fatalExit(fmt.Errorf("Block not found"))
	}
	logger.Info("Block reimported", zap.String("number", blockID))
}

//ReloadTransaction removes a transaction, cleans DB and reimport TX
func ReloadTransaction(ctx context.Context, backendInstance *backend.Backend, txHash string) {
	tx, err := backendInstance.GetTransactionByHash(ctx, txHash)
	if err != nil || tx == nil {
		//it's impossible to reload transaction that has not been imported since Goclient doesn't show a block number in a transaction object (but should)
		//https://github.com/gochain/gochain/blob/83b0633553802c066e5073f969ac740c7c6cb3bb/core/types/transaction.go#L47-L73
		//but API has it - https://github.com/ethereum/wiki/wiki/JavaScript-API#web3ethgettransaction (BlockNumber)
		fatalExit(fmt.Errorf("Failed to get transaction:%s, %v", txHash, err))
	}
	logger.Info("Deleting block ", zap.Int64("number", tx.BlockNumber))
	err = backendInstance.DeleteBlockByNumber(tx.BlockNumber)
	if err != nil {
		fatalExit(err)
	}
	_, err = backendInstance.GetBlockByNumber(ctx, tx.BlockNumber, false)
	if err != nil {
		fatalExit(fmt.Errorf("Failed to get block:%d, %v", tx.BlockNumber, err))
	}
}

//ReloadContract removes contract from db, delete all tokens and transactions and reimports everything
func ReloadContract(ctx context.Context, backendInstance *backend.Backend, address string) {
	if !common.IsHexAddress(address) {
		fatalExit(fmt.Errorf("invalid hex address: %s", address))
	}
	addr := common.HexToAddress(address)
	normalizedAddress := addr.Hex()
	err := backendInstance.DeleteContract(normalizedAddress)
	if err != nil {
		fatalExit(fmt.Errorf("failed to delete contract %v", err))
	}
	balance, err := backendInstance.Balance(ctx, addr)
	if err != nil {
		fatalExit(fmt.Errorf("failed to get balance"))
	}
	currentBlock, err := backendInstance.GetLatestBlockNumber(ctx)
	if err != nil {
		logger.Error("Update Addresses: Failed to get latest block number", zap.Error(err))
	}
	contractDataArray, err := backendInstance.CodeAt(ctx, normalizedAddress)
	if err != nil {
		fatalExit(fmt.Errorf("failed to get code %v", err))
	}
	contractData := string(contractDataArray[:])
	var tokenDetails = &tokens.TokenDetails{TotalSupply: big.NewInt(0)}
	contract := false
	if contractData != "" {
		contract = true
		byteCode := hex.EncodeToString(contractDataArray)
		if err := backendInstance.ImportContract(normalizedAddress, byteCode); err != nil {
			fatalExit(fmt.Errorf("failed to import contract: %v", err))
		}
		contractFromDB, err := backendInstance.GetAddressByHash(ctx, normalizedAddress)
		if err != nil {
			fatalExit(fmt.Errorf("failed to get contract from DB: %v", err))
		}
		fromBlock, err := backendInstance.GetContractBlock(normalizedAddress)
		if err != nil {
			fatalExit(fmt.Errorf("failed to get contract block: %v", err))
		}
		tokenDetails, err = backendInstance.GetTokenDetails(normalizedAddress, byteCode)
		if err != nil {
			fatalExit(fmt.Errorf("failed to get token details: %v", err))
		}
		contractFromDB.TokenName = tokenDetails.Name
		contractFromDB.TokenSymbol = tokenDetails.Symbol
		tokenTransfers, err := backendInstance.GetTransferEvents(ctx, tokenDetails, fromBlock, blockRangeLimit)
		if err != nil {
			fatalExit(fmt.Errorf("failed to get internal txs: %v", err))
		}
		tokenHoldersList := make(map[string]struct{})
		for _, itx := range tokenTransfers {
			logger.Debug("Internal Transaction", zap.Stringer("from", itx.From),
				zap.Stringer("to", itx.To), zap.Stringer("value", itx.Value))
			if _, err := backendInstance.ImportTransferEvent(ctx, normalizedAddress, itx); err != nil {
				fatalExit(fmt.Errorf("failed to import internal tx: %v", err))
			}
			logger.Debug("Updating following token holder addresses", zap.Stringer("from", itx.From),
				zap.Stringer("to", itx.To), zap.Stringer("value", itx.Value))
			tokenHoldersList[itx.To.String()] = struct{}{}
			tokenHoldersList[itx.From.String()] = struct{}{}
		}
		for tokenHolderAddress := range tokenHoldersList {
			if tokenHolderAddress == "0x0000000000000000000000000000000000000000" {
				continue
			}
			logger.Info("Importing token holder", zap.String("holder", tokenHolderAddress),
				zap.Int("total", len(tokenHoldersList)))
			tokenHolder, err := backendInstance.GetTokenBalance(normalizedAddress, tokenHolderAddress)
			if err != nil {
				logger.Error("Failed to get token balance", zap.Error(err), zap.String("holder", tokenHolderAddress))
				fatalExit(fmt.Errorf("failed to get balance"))
			}
			if _, err := backendInstance.ImportTokenHolder(normalizedAddress, tokenHolderAddress, tokenHolder, contractFromDB); err != nil {
				fatalExit(fmt.Errorf("failed to import token holder: %v", err))
			}
		}
	}
	logger.Info("Update Addresses: updated address", zap.String("Address", normalizedAddress), zap.Stringer("balance", balance))
	_, err = backendInstance.ImportAddress(normalizedAddress, balance, tokenDetails, contract, currentBlock.Int64())
	if err != nil {
		fatalExit(fmt.Errorf("failed to import address %v", err))
	}
	logger.Info("Address successfully updated", zap.String("address", normalizedAddress))
}
