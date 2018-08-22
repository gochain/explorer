package backend

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/gochain-io/gochain/common"
	"github.com/gochain-io/gochain/ethclient"
	"github.com/rs/zerolog/log"
)

type TokenBalance struct {
	url string
}

var (
	conn    *ethclient.Client
	config  *Config
	VERSION string
)

func NewTokenBalanceClient(rpcUrl string) *TokenBalance {
	return &TokenBalance{
		url: rpcUrl,
	}
}
func (rpc *TokenBalance) GetTokenBalance(contract, wallet string) (*tokenBalance, error) {
	configs := &Config{
		GethLocation: rpc.url,
		Logs:         true,
	}
	configs.Connect()
	return New(contract, wallet)

}

func New(contract, wallet string) (*tokenBalance, error) {
	var err error
	if config == nil || conn == nil {
		return nil, errors.New("geth server connection has not been created")
	}
	tb := &tokenBalance{
		Contract: common.HexToAddress(contract),
		Wallet:   common.HexToAddress(wallet),
		Decimals: 0,
		Balance:  big.NewInt(0),
		ctx:      context.TODO(),
	}
	err = tb.query()
	return tb, err
}

func (c *Config) Connect() error {
	var err error
	if c.GethLocation == "" {
		return errors.New("geth endpoint has not been set")
	}
	ethConn, err := ethclient.Dial(c.GethLocation)
	if err != nil {
		return err
	}
	block, err := ethConn.BlockByNumber(context.TODO(), nil)
	if block == nil {
		return err
	}
	config = c
	conn = ethConn
	log.Info().Str("url", c.GethLocation).Msg("Connected to Geth at")
	return err
}

func (tb *tokenBalance) ETHString() string {
	return BigIntString(tb.ETH, 18)
}

func (tb *tokenBalance) BalanceString() string {
	if tb.Decimals == 0 {
		return "0"
	}
	return BigIntString(tb.Balance, tb.Decimals)
}

func (tb *tokenBalance) query() error {
	var err error

	token, err := NewTokenCaller(tb.Contract, conn)
	if err != nil {
		log.Info().Err(err).Msg("Failed to instantiate a Token contract")
		return err
	}

	block, err := conn.BlockByNumber(tb.ctx, nil)
	if err != nil {
		log.Info().Err(err).Msg("Failed to get current block number")
	}
	tb.Block = block.Number().Int64()

	decimals, err := token.Decimals(nil)
	if err != nil {
		log.Info().Err(err).Str("Contract", tb.Contract.String()).Msg("Failed to get decimals from contract")
		return err
	}
	tb.Decimals = decimals.Int64()

	tb.ETH, err = conn.BalanceAt(tb.ctx, tb.Wallet, nil)
	if err != nil {
		log.Info().Err(err).Str("Wallet", tb.Wallet.String()).Msg("Failed to get ethereum balance from address")
	}

	tb.Balance, err = token.BalanceOf(nil, tb.Wallet)
	if err != nil {
		log.Info().Err(err).Str("Wallet", tb.Contract.String()).Msg("Failed to get balance from contract")
		tb.Balance = big.NewInt(0)
	}

	tb.Symbol, err = token.Symbol(nil)
	if err != nil {
		log.Info().Err(err).Str("Wallet", tb.Contract.String()).Msg("Failed to get symbol from contract")
		tb.Symbol = symbolFix(tb.Contract.String())
	}

	tb.Name, err = token.Name(nil)
	if err != nil {
		log.Info().Err(err).Str("Wallet", tb.Contract.String()).Msg("Failed to retrieve token name from contract")
		tb.Name = "MISSING"
	}

	return err
}

func symbolFix(contract string) string {
	switch common.HexToAddress(contract).String() {
	case "0x86Fa049857E0209aa7D9e616F7eb3b3B78ECfdb0":
		return "EOS"
	}
	return "MISSING"
}

func (tb *tokenBalance) ToJSON() tokenBalanceJson {
	jsonData := tokenBalanceJson{
		Contract: tb.Contract.String(),
		Wallet:   tb.Wallet.String(),
		Name:     tb.Name,
		Symbol:   tb.Symbol,
		Balance:  tb.BalanceString(),
		ETH:      tb.ETHString(),
		Decimals: tb.Decimals,
		Block:    tb.Block,
	}
	return jsonData
}

func BigIntString(balance *big.Int, decimals int64) string {
	amount := BigIntFloat(balance, decimals)
	deci := fmt.Sprintf("%%0.%vf", decimals)
	return clean(fmt.Sprintf(deci, amount))
}

func BigIntFloat(balance *big.Int, decimals int64) *big.Float {
	if balance.Sign() == 0 {
		return big.NewFloat(0)
	}
	bal := big.NewFloat(0)
	bal.SetInt(balance)
	pow := bigPow(10, decimals)
	p := big.NewFloat(0)
	p.SetInt(pow)
	bal.Quo(bal, p)
	return bal
}

func bigPow(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

func clean(newNum string) string {
	stringBytes := bytes.TrimRight([]byte(newNum), "0")
	newNum = string(stringBytes)
	if stringBytes[len(stringBytes)-1] == 46 {
		newNum += "0"
	}
	if stringBytes[0] == 46 {
		newNum = "0" + newNum
	}
	return newNum
}
