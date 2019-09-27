package backend

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/gochain-io/explorer/server/utils"
	"go.uber.org/zap"

	"github.com/gochain-io/gochain/v3"
	"github.com/gochain-io/gochain/v3/common"
	"github.com/gochain-io/gochain/v3/core/types"
	"github.com/gochain-io/gochain/v3/goclient"
)

type TokenDetails struct {
	Contract    common.Address
	Name        string
	Symbol      string
	TotalSupply *big.Int
	Decimals    int64
	Block       int64
	Types       []utils.ErcName
	Interfaces  []utils.FunctionName
}

type TokenHolderDetails struct {
	Contract    common.Address
	TokenHolder common.Address
	Balance     *big.Int
	Block       int64
}

type TokenBalance struct {
	url                string
	conn               *goclient.Client
	initialBlockNumber int64
	Lgr                *zap.Logger
}

type TransferEvent struct {
	From            common.Address
	To              common.Address
	Value           *big.Int
	BlockNumber     int64
	TransactionHash string
}

func NewTokenBalanceClient(ctx context.Context, client *goclient.Client, lgr *zap.Logger) (*TokenBalance, error) {
	num, err := client.LatestBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %v", err)
	}
	return &TokenBalance{
		conn:               client,
		initialBlockNumber: num.Int64(),
		Lgr:                lgr,
	}, nil
}

func (rpc *TokenBalance) GetTokenHolderDetails(contract, wallet string) (*TokenHolderDetails, error) {
	if rpc.conn == nil {
		return nil, errors.New("geth server connection has not been created")
	}
	th := &TokenHolderDetails{
		Contract:    common.HexToAddress(contract),
		TokenHolder: common.HexToAddress(wallet),
		Balance:     big.NewInt(0),
	}
	err := th.queryTokenHolderDetails(rpc.conn, rpc.Lgr)
	return th, err
}

func (rpc *TokenBalance) GetTokenDetails(contractAddress string, byteCode string) (*TokenDetails, error) {
	if rpc.conn == nil {
		return nil, errors.New("geth server connection has not been created")
	}
	tb := &TokenDetails{
		Contract:    common.HexToAddress(contractAddress),
		Decimals:    0,
		TotalSupply: big.NewInt(0),
	}
	err := tb.queryTokenDetails(rpc.conn, byteCode, rpc.Lgr)
	return tb, err
}

func (th *TokenHolderDetails) queryTokenHolderDetails(conn *goclient.Client, lgr *zap.Logger) error {
	token := NewTokenCaller(th.Contract, conn)
	var err error
	th.Balance, err = token.BalanceOf(nil, th.TokenHolder)
	if err != nil {
		lgr.Info("Failed to get balance from contract", zap.Error(err), zap.Stringer("wallet", th.Contract))
		th.Balance = big.NewInt(0)
	}
	return err
}

func (tb *TokenDetails) queryTokenDetails(conn *goclient.Client, byteCode string, lgr *zap.Logger) error {
	token := NewTokenCaller(tb.Contract, conn)
	var err error
	tb.Types, tb.Interfaces = token.GetInfo(byteCode)
	for _, interfaceName := range tb.Interfaces {
		if utils.InterfaceIdentifiers[interfaceName].Callable {
			switch interfaceName {
			case utils.Decimals:
				decimals, err := token.Decimals(nil)
				if err != nil {
					lgr.Error("Failed to get decimals from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
					continue
				}
				tb.Decimals = decimals.Int64()
			case utils.TotalSupply:
				totalSupply, err := token.TotalSupply(nil)
				if err != nil {
					lgr.Error("Failed to get total supply", zap.Error(err), zap.Stringer("address", tb.Contract))
					tb.TotalSupply = big.NewInt(0)
					continue
				}
				tb.TotalSupply = totalSupply
			case utils.Symbol:
				tb.Symbol, err = token.Symbol(nil)
				if err != nil {
					lgr.Error("Failed to get symbol from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
					tb.Symbol = "MISSING"
				}
			case utils.Name:
				tb.Name, err = token.Name(nil)
				if err != nil {
					lgr.Error("Failed to retrieve token name from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
					tb.Name = "MISSING"
				}
			}
		}
	}

	return err
}

func (rpc *TokenBalance) getInternalTransactions(ctx context.Context, address string, contractBlock int64, blockRangeLimit uint64) []TransferEvent {
	numOfBlocksPerRequest := int64(blockRangeLimit)
	latestBlockNumber := rpc.initialBlockNumber

	if num, err := rpc.conn.LatestBlockNumber(ctx); err != nil {
		rpc.Lgr.Error("Failed to get latest block number", zap.Error(err))
	} else {
		latestBlockNumber = num.Int64()
	}
	contractBlock -= numOfBlocksPerRequest
	numOfCycles := int((latestBlockNumber - contractBlock) / numOfBlocksPerRequest)
	contractAddress := common.HexToAddress(address)
	transferOperation := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	var transferEvents []TransferEvent
	for i := 0; i <= numOfCycles; i++ {
		fromBlock := contractBlock + int64(i)*numOfBlocksPerRequest
		rpc.Lgr.Debug("Querying for internal txs", zap.Int64("from", fromBlock), zap.Int64("to", fromBlock+numOfBlocksPerRequest),
			zap.Int64("block", contractBlock), zap.Int64("latest", latestBlockNumber), zap.Int("events", len(transferEvents)))
		query := gochain.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(fromBlock + numOfBlocksPerRequest),
			Addresses: []common.Address{contractAddress},
			Topics:    [][]common.Hash{[]common.Hash{transferOperation}},
		}

		var events []types.Log
		err := retry(ctx, 5, 2*time.Second, func() (err error) {
			events, err = rpc.conn.FilterLogs(ctx, query)
			return err
		})
		if err != nil {
			rpc.Lgr.Error("Failed to query for internal txs", zap.Error(err))
		}
		for _, event := range events {

			var transferEvent TransferEvent
			err = TokenABI.Unpack(&transferEvent, "Transfer", event.Data)
			if err != nil {
				rpc.Lgr.Warn("Failed to unpack event", zap.Error(err), zap.Uint64("block", event.BlockNumber),
					zap.Uint("index", event.Index), zap.String("contract", address))
				continue
			}
			if l := len(event.Topics); l < 3 {
				rpc.Lgr.Warn("Failed to parse event - too few topics. Expected 3.", zap.Error(err),
					zap.Uint64("block", event.BlockNumber), zap.Uint("index", event.Index), zap.String("contract", address))
				continue
			}
			transferEvent.From = common.BytesToAddress(event.Topics[1].Bytes())
			transferEvent.To = common.BytesToAddress(event.Topics[2].Bytes())
			transferEvent.BlockNumber = int64(event.BlockNumber)
			transferEvent.TransactionHash = event.TxHash.String()
			rpc.Lgr.Debug("event", zap.Uint64("block", event.BlockNumber), zap.Uint("index", event.Index),
				zap.String("contract", address), zap.String("from", transferEvent.From.Hex()),
				zap.Stringer("to", transferEvent.To), zap.Stringer("value", transferEvent.Value))
			transferEvents = append(transferEvents, transferEvent)
		}
	}
	return transferEvents
}
