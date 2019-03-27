package models

import (
	"github.com/gochain-io/explorer/server/utils"
	"time"
)

type Address struct {
	Address        string    `json:"address" bson:"address"`
	Owner          string    `json:"owner,omitempty" bson:"owner"`
	BalanceFloat   float64   `json:"-" bson:"balance_float"`        //low precise balance for sorting purposes
	BalanceString  string    `json:"balance" bson:"balance_string"` //high precise balance for API
	BalanceWei     string    `json:"balance_wei" bson:"balance_wei"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedAtBlock int64     `json:"-" bson:"updated_at_block"`
	TokenName      string    `json:"token_name,omitempty" bson:"token_name"`
	TokenSymbol    string    `json:"token_symbol,omitempty" bson:"token_symbol"`
	Decimals       int64     `json:"decimals,omitempty" bson:"decimals"`
	TotalSupply    string    `json:"total_supply" bson:"total_supply"`
	Contract       bool      `json:"contract" bson:"contract"`
	// Depreciated and gonna be deleted
	// Should use ErcTypes
	GO20                         bool            `json:"go20" bson:"go20"`
	ErcTypes                     []utils.ErcName `json:"erc_types" bson:"erc_types"`
	NumberOfTransactions         int             `json:"number_of_transactions" bson:"number_of_transactions"`
	NumberOfTokenHolders         int             `json:"number_of_token_holders,omitempty" bson:"number_of_token_holders"`
	NumberOfInternalTransactions int             `json:"number_of_internal_transactions,omitempty" bson:"number_of_internal_transactions"`
}

type AddressesList struct {
	Adresses []*Address `json:"adresses"`
}
