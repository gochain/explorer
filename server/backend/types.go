package backend

import (
	"context"
	"math/big"

	"github.com/gochain-io/gochain/common"
)

type Config struct {
	GethLocation string
	Logs         bool
}

type TokenDetails struct {
	Contract    common.Address
	Wallet      common.Address
	Name        string
	Symbol      string
	Balance     *big.Int
	TotalSupply *big.Int
	ETH         *big.Int
	Decimals    int64
	Block       int64
	ctx         context.Context
}
