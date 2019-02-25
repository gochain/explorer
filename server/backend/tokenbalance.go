package backend

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/gochain-io/gochain/v3"
	"github.com/gochain-io/gochain/v3/accounts/abi"
	"github.com/gochain-io/gochain/v3/common"
	"github.com/gochain-io/gochain/v3/goclient"
	"github.com/rs/zerolog/log"
)

type TokenDetails struct {
	Contract    common.Address
	Name        string
	Symbol      string
	TotalSupply *big.Int
	Decimals    int64
	Block       int64
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
}

type TransferEvent struct {
	From            common.Address
	To              common.Address
	Value           *big.Int
	BlockNumber     int64
	TransactionHash string
}

func NewTokenBalanceClient(rpcUrl string) *TokenBalance {
	var err error

	if rpcUrl == "" {
		log.Fatal().Msg("geth endpoint has not been set")
	}
	ethConn, err := goclient.Dial(rpcUrl)

	if err != nil {
		log.Fatal().Err(err).Msg("NewTokenBalanceClient")
	}
	block, err := ethConn.BlockByNumber(context.TODO(), nil)
	if block == nil {
		log.Fatal().Err(err).Msg("NewTokenBalanceClient")
	}
	log.Info().Str("url", rpcUrl).Msg("Connected to Geth at")

	return &TokenBalance{
		url:                rpcUrl,
		conn:               ethConn,
		initialBlockNumber: block.Number().Int64(),
	}
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
	err := th.queryTokenHolderDetails(rpc.conn)
	return th, err
}

func (rpc *TokenBalance) GetTokenDetails(contract string) (*TokenDetails, error) {
	if rpc.conn == nil {
		return nil, errors.New("geth server connection has not been created")
	}
	tb := &TokenDetails{
		Contract:    common.HexToAddress(contract),
		Decimals:    0,
		TotalSupply: big.NewInt(0),
	}
	err := tb.queryTokenDetails(rpc.conn)
	return tb, err
}

func (th *TokenHolderDetails) queryTokenHolderDetails(conn *goclient.Client) error {
	var err error

	token, err := NewTokenCaller(th.Contract, conn)
	if err != nil {
		log.Info().Err(err).Msg("Failed to instantiate a Token contract")
		return err
	}
	th.Balance, err = token.BalanceOf(nil, th.TokenHolder)
	if err != nil {
		log.Info().Err(err).Str("Wallet", th.Contract.String()).Msg("Failed to get balance from contract")
		th.Balance = big.NewInt(0)
	}
	return err
}

func (tb *TokenDetails) queryTokenDetails(conn *goclient.Client) error {
	var err error

	token, err := NewTokenCaller(tb.Contract, conn)
	if err != nil {
		log.Info().Err(err).Msg("Failed to instantiate a Token contract")
		return err
	}

	decimals, err := token.Decimals(nil)
	if err != nil {
		log.Info().Err(err).Str("Contract", tb.Contract.String()).Msg("Failed to get decimals from contract")
		return err
	}
	tb.Decimals = decimals.Int64()

	totalSupply, err := token.TotalSupply(nil)
	if err != nil {
		log.Info().Err(err).Str("Contract", tb.Contract.String()).Msg("Failed to get total supply")
		tb.TotalSupply = big.NewInt(0)
		return err
	}
	tb.TotalSupply = totalSupply

	tb.Symbol, err = token.Symbol(nil)
	if err != nil {
		log.Info().Err(err).Str("Wallet", tb.Contract.String()).Msg("Failed to get symbol from contract")
		tb.Symbol = "MISSING"
	}

	tb.Name, err = token.Name(nil)
	if err != nil {
		log.Info().Err(err).Str("Wallet", tb.Contract.String()).Msg("Failed to retrieve token name from contract")
		tb.Name = "MISSING"
	}

	return err
}

func (rpc *TokenBalance) getInternalTransactions(address string, contractBlock int64) []TransferEvent {
	const numOfBlocksPerRequest = 100000
	latestBlockNumber := rpc.initialBlockNumber
	block, err := rpc.conn.BlockByNumber(context.Background(), nil)
	if block == nil {
		log.Error().Err(err).Msg("getInternalTransactions")
	} else {
		latestBlockNumber = block.Number().Int64()
	}
	contractBlock -= numOfBlocksPerRequest
	numOfCycles := int((latestBlockNumber - contractBlock) / numOfBlocksPerRequest)
	contractAddress := common.HexToAddress(address)
	transferOperation := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	var transferEvents []TransferEvent
	for i := 0; i <= numOfCycles; i++ {
		fromBlock := contractBlock + int64(i)*numOfBlocksPerRequest
		log.Info().Int64("From block", fromBlock).Int64("To Block", fromBlock+numOfBlocksPerRequest).Int64("Block from request", contractBlock).Int64("Latest block", latestBlockNumber).Int("Number of the events", len(transferEvents)).Msg("list of transactions")
		query := gochain.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(fromBlock + numOfBlocksPerRequest),
			Addresses: []common.Address{contractAddress},
			Topics:    [][]common.Hash{[]common.Hash{transferOperation}},
		}
		ctx := context.Background()

		events, err := rpc.conn.FilterLogs(ctx, query)

		if err != nil {
			log.Info().Err(err)
		}
		tokenAbi, err := abi.JSON(strings.NewReader(string(TokenABI)))

		if err != nil {
			log.Info().Err(err)
		}
		for _, event := range events {
			var transferEvent TransferEvent
			err = tokenAbi.Unpack(&transferEvent, "Transfer", event.Data)
			if err != nil {
				log.Info().Str("Address", address).Msg("Failed to unpack event")
				continue
			}
			transferEvent.From = common.BytesToAddress(event.Topics[1].Bytes())
			transferEvent.To = common.BytesToAddress(event.Topics[2].Bytes())
			transferEvent.BlockNumber = int64(event.BlockNumber)
			transferEvent.TransactionHash = event.TxHash.String()
			log.Debug().Uint64("Block", event.BlockNumber).Str("Contract", address).Str("From", transferEvent.From.Hex()).Str("To", transferEvent.To.Hex()).Str("Value", transferEvent.Value.String()).Msg("event")
			transferEvents = append(transferEvents, transferEvent)
		}
	}
	return transferEvents
}
