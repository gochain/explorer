package tokens

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/gochain-io/explorer/server/utils"

	"github.com/gochain/gochain/v4"
	"github.com/gochain/gochain/v4/accounts/abi"
	"github.com/gochain/gochain/v4/accounts/abi/bind"
	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/core/types"
	"github.com/gochain/gochain/v4/goclient"
	"go.uber.org/zap"
)

var ABIs struct {
	ERC20, ERC721 abi.ABI
}

func init() {
	var err error
	ABIs.ERC20, err = abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		fmt.Println("Failed to parse token ERC20 ABI:", err)
		os.Exit(1)
	}
	ABIs.ERC721, err = abi.JSON(strings.NewReader(ERC721ABI))
	if err != nil {
		fmt.Println("Failed to parse token ERC721 ABI:", err)
		os.Exit(1)
	}
}

type TokenDetails struct {
	Contract common.Address
	Block    int64

	ErcTypes  map[utils.EVMInterface]struct{}
	Functions map[utils.EVMFunction]struct{}

	Name        string
	Symbol      string
	TotalSupply *big.Int
	Decimals    int64

	Target string
	Owner  string
}

func (t *TokenDetails) ERCTypesSlice() []string {
	ercTypes := make([]string, 0, len(t.ErcTypes))
	for et := range t.ErcTypes {
		ercTypes = append(ercTypes, et.String())
	}
	return ercTypes
}

func (t *TokenDetails) FunctionsSlice() []string {
	functions := make([]string, 0, len(t.Functions))
	for et := range t.Functions {
		functions = append(functions, et.String())
	}
	return functions
}

type TokenHolderDetails struct {
	Contract    common.Address
	TokenHolder common.Address
	Balance     *big.Int
	Block       int64
}

type TokenClient struct {
	url                string
	conn               *goclient.Client
	initialBlockNumber int64
	Lgr                *zap.Logger
}

// TransferEvent represents a Transfer event emitted from an ERC20 or ERC721.
type TransferEvent struct {
	From            common.Address
	To              common.Address
	Value           *big.Int // ERC20:value/ERC721:tokenId
	BlockNumber     int64
	TransactionHash string
}

func NewERC20Balance(ctx context.Context, client *goclient.Client, lgr *zap.Logger) (*TokenClient, error) {
	num, err := client.LatestBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %v", err)
	}
	return &TokenClient{
		conn:               client,
		initialBlockNumber: num.Int64(),
		Lgr:                lgr,
	}, nil
}

func (rpc *TokenClient) GetTokenHolderDetails(contract, wallet string) (*TokenHolderDetails, error) {
	if !common.IsHexAddress(contract) {
		return nil, fmt.Errorf("invalid hex address: %s", contract)
	}
	if !common.IsHexAddress(wallet) {
		return nil, fmt.Errorf("invalid hex address: %s", wallet)
	}
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

func (rpc *TokenClient) GetTokenDetails(contractAddress string, byteCode string) (*TokenDetails, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", contractAddress)
	}
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
	token, err := NewERC20Caller(th.Contract, conn)
	if err != nil {
		return err
	}
	th.Balance, err = token.BalanceOf(nil, th.TokenHolder)
	if err != nil {
		lgr.Info("Failed to get balance from contract", zap.Error(err), zap.Stringer("wallet", th.Contract))
		th.Balance = big.NewInt(0)
	}
	return err
}

func (tb *TokenDetails) queryTokenDetails(conn *goclient.Client, byteCode string, lgr *zap.Logger) error {
	var count int
	tb.ErcTypes, tb.Functions, count = utils.ScanContract(byteCode)
	lgr.Info("Analyzing contract byte code", zap.Stringer("address", tb.Contract),
		zap.Strings("ercTypes", tb.ERCTypesSlice()), zap.Strings("functions", tb.FunctionsSlice()),
		zap.Int("cachedContracts", count))
	if _, ok := tb.ErcTypes[utils.Upgradeable]; ok {
		if err := tb.queryTarget(conn, lgr); err != nil {
			return err
		}
		//TODO ensure target loaded in db, and merge its types and funcs in to this set
	}
	if _, ok := tb.Functions[utils.Owner]; ok {
		if err := tb.queryOwner(conn, lgr); err != nil {
			return err
		}
	}
	if _, ok := tb.ErcTypes[utils.Go20]; ok {
		return tb.queryERC20Details(conn, lgr)
	}
	if _, ok := tb.ErcTypes[utils.Go721]; ok {
		return tb.queryERC721Details(conn, lgr)
	}
	return nil
}

func (tb *TokenDetails) queryTarget(conn *goclient.Client, lgr *zap.Logger) error {
	token, err := NewUpgradeableCaller(tb.Contract, conn)
	if err != nil {
		return err
	}
	target, err := token.Target(nil)
	if err != nil {
		lgr.Error("Failed to get target from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
	} else {
		tb.Target = target.Hex()
	}
	return nil
}

func (tb *TokenDetails) queryOwner(conn *goclient.Client, lgr *zap.Logger) error {
	token, err := NewUpgradeableCaller(tb.Contract, conn)
	if err != nil {
		return err
	}
	owner, err := token.Owner(nil)
	if err != nil {
		lgr.Error("Failed to get owner", zap.Error(err), zap.Stringer("address", tb.Contract))
	} else {
		tb.Owner = owner.Hex()
	}
	return nil
}

func (tb *TokenDetails) queryERC20Details(conn *goclient.Client, lgr *zap.Logger) error {
	token, err := NewERC20Caller(tb.Contract, conn)
	if err != nil {
		return err
	}

	if _, ok := tb.Functions[utils.Decimals]; ok {
		decimals, err := token.Decimals(nil)
		if err != nil {
			if noCode, emptyResp := contractErr(err); noCode || emptyResp {
				lgr.Warn("Failed to get decimals from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
			} else {
				lgr.Error("Failed to get decimals from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
			}
		} else {
			tb.Decimals = int64(decimals)
		}
	}
	if _, ok := tb.Functions[utils.TotalSupply]; ok {
		totalSupply, err := token.TotalSupply(nil)
		if err != nil {
			if noCode, emptyResp := contractErr(err); noCode || emptyResp {
				lgr.Warn("Failed to get total supply", zap.Error(err), zap.Stringer("address", tb.Contract))
			} else {
				lgr.Error("Failed to get total supply", zap.Error(err), zap.Stringer("address", tb.Contract))
			}
			tb.TotalSupply = big.NewInt(0)
		} else {
			tb.TotalSupply = totalSupply
		}
	}
	if _, ok := tb.Functions[utils.Symbol]; ok {
		tb.Symbol, err = token.Symbol(nil)
		if err != nil {
			if noCode, emptyResp := contractErr(err); noCode || emptyResp {
				lgr.Warn("Failed to get symbol from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
			} else {
				lgr.Error("Failed to get symbol from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
			}
			tb.Symbol = "MISSING"
		}
	}
	if _, ok := tb.Functions[utils.Name]; ok {
		tb.Name, err = token.Name(nil)
		if err != nil {
			if noCode, emptyResp := contractErr(err); noCode || emptyResp {
				lgr.Error("Failed to retrieve token name from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
			} else {
				lgr.Error("Failed to retrieve token name from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
			}
			tb.Name = "MISSING"
		}
	}

	return nil
}

func contractErr(err error) (noCode, emptyResp bool) {
	noCode = err == bind.ErrNoCode
	if noCode {
		return
	}
	emptyResp = strings.Contains(err.Error(),
		"abi: attempting to unmarshall an empty string while arguments are expected")
	return
}

func (tb *TokenDetails) queryERC721Details(conn *goclient.Client, lgr *zap.Logger) error {
	token, err := NewERC721Caller(tb.Contract, conn)
	if err != nil {
		return err
	}

	if _, ok := tb.Functions[utils.TotalSupply]; ok {
		totalSupply, err := token.TotalSupply(nil)
		if err != nil {
			if noCode, emptyResp := contractErr(err); noCode || emptyResp {
				lgr.Warn("Failed to get total supply", zap.Error(err), zap.Stringer("address", tb.Contract))
			} else {
				lgr.Error("Failed to get total supply", zap.Error(err), zap.Stringer("address", tb.Contract))
			}
			tb.TotalSupply = big.NewInt(0)
		} else {
			tb.TotalSupply = totalSupply
		}
	}
	if _, ok := tb.Functions[utils.Symbol]; ok {
		tb.Symbol, err = token.Symbol(nil)
		if err != nil {
			if noCode, emptyResp := contractErr(err); noCode || emptyResp {
				lgr.Warn("Failed to get symbol from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
			} else {
				lgr.Error("Failed to get symbol from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
			}
			tb.Symbol = "MISSING"
		}
	}
	if _, ok := tb.Functions[utils.Name]; ok {
		tb.Name, err = token.Name(nil)
		if err != nil {
			if noCode, emptyResp := contractErr(err); noCode || emptyResp {
				lgr.Warn("Failed to retrieve token name from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
			} else {
				lgr.Error("Failed to retrieve token name from contract", zap.Error(err), zap.Stringer("address", tb.Contract))
			}
			tb.Name = "MISSING"
		}
	}

	return nil
}

var transferEventID = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

func (rpc *TokenClient) GetTransferEvents(ctx context.Context, tokenDetails *TokenDetails, contractBlock int64, blockRangeLimit uint64) ([]*TransferEvent, error) {
	lgr := rpc.Lgr.With(zap.Stringer("address", tokenDetails.Contract))
	var unpackTransferEvent func(log types.Log) (*TransferEvent, error)
	if _, ok := tokenDetails.ErcTypes[utils.Go20]; ok {
		unpackTransferEvent = unpackERC20TransferEvent
	} else if _, ok := tokenDetails.ErcTypes[utils.Go721]; ok {
		unpackTransferEvent = unpackERC721TransferEvent
	} else {
		return nil, nil
	}

	numOfBlocksPerRequest := int64(blockRangeLimit)
	latestBlockNumber := rpc.initialBlockNumber
	if num, err := rpc.conn.LatestBlockNumber(ctx); err != nil {
		lgr.Warn("Failed to get latest block number", zap.Error(err))
	} else {
		latestBlockNumber = num.Int64()
	}
	contractBlock -= numOfBlocksPerRequest
	if contractBlock < 0 {
		contractBlock = 0
	}
	numOfCycles := int((latestBlockNumber - contractBlock) / numOfBlocksPerRequest)
	var transferEvents []*TransferEvent
	for i := 0; i <= numOfCycles; i++ {
		fromBlock := contractBlock + int64(i)*numOfBlocksPerRequest
		lgr.Debug("Querying for token transfer events", zap.Int64("from", fromBlock), zap.Int64("to", fromBlock+numOfBlocksPerRequest),
			zap.Int64("block", contractBlock), zap.Int64("latest", latestBlockNumber), zap.Int("events", len(transferEvents)))
		query := gochain.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(fromBlock + numOfBlocksPerRequest),
			Addresses: []common.Address{tokenDetails.Contract},
			Topics:    [][]common.Hash{{transferEventID}},
		}

		var logs []types.Log
		err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
			logs, err = rpc.conn.FilterLogs(ctx, query)
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query RPC for logs: %v", err)
		}
		for _, log := range logs {
			event, err := unpackTransferEvent(log)
			if err != nil {
				lgr.Error("Failed to unpack event", zap.Error(err), zap.Uint64("block", log.BlockNumber),
					zap.Uint("index", log.Index), zap.Stringer("contract", log.Address))
				continue
			}
			transferEvents = append(transferEvents, event)
		}
	}
	return transferEvents, nil
}

var addrTopicPrefix = make([]byte, 12)

func unpackERC20TransferEvent(event types.Log) (*TransferEvent, error) {
	if l := len(event.Topics); l != 3 {
		return nil, fmt.Errorf("incorrect number of topics: %d", l)
	}
	from := event.Topics[1].Bytes()
	to := event.Topics[2].Bytes()
	if len(bytes.TrimPrefix(from, addrTopicPrefix)) != 20 {
		return nil, fmt.Errorf("from topic longer than address: %s", string(from))
	}
	if len(bytes.TrimPrefix(to, addrTopicPrefix)) != 20 {
		return nil, fmt.Errorf("to topic longer than address: %s", string(from))
	}
	var transferEvent TransferEvent
	err := ABIs.ERC20.UnpackIntoInterface(&transferEvent, "Transfer", event.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack log data: %v", err)
	}
	transferEvent.From = common.BytesToAddress(from)
	transferEvent.To = common.BytesToAddress(to)
	transferEvent.BlockNumber = int64(event.BlockNumber)
	transferEvent.TransactionHash = event.TxHash.String()
	return &transferEvent, nil
}

func unpackERC721TransferEvent(event types.Log) (*TransferEvent, error) {
	if l := len(event.Topics); l != 4 {
		return nil, fmt.Errorf("incorrect number of topics: %d", l)
	}
	if len(event.Data) > 0 {
		return nil, fmt.Errorf("invalid extra data")
	}
	from := event.Topics[1].Bytes()
	to := event.Topics[2].Bytes()
	if len(bytes.TrimPrefix(from, addrTopicPrefix)) != 20 {
		return nil, fmt.Errorf("from topic longer than address: %s", string(from))
	}
	if len(bytes.TrimPrefix(to, addrTopicPrefix)) != 20 {
		return nil, fmt.Errorf("to topic longer than address: %s", string(from))
	}
	var transferEvent TransferEvent
	transferEvent.From = common.BytesToAddress(from)
	transferEvent.To = common.BytesToAddress(to)
	transferEvent.Value = event.Topics[3].Big()
	transferEvent.BlockNumber = int64(event.BlockNumber)
	transferEvent.TransactionHash = event.TxHash.String()
	return &transferEvent, nil
}
