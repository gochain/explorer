package models

import (
	"time"
)

type Address struct {
	Address        string    `json:"address" bson:"address"`
	BalanceFloat   float64   `json:"-" bson:"balance_float"`        //low precise balance for sorting purposes
	BalanceString  string    `json:"balance" bson:"balance_string"` //high precise balance for API
	BalanceWei     string    `json:"balance_wei" bson:"balance_wei"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedAtBlock int64     `json:"-" bson:"updated_at_block"`

	TokenName   string `json:"token_name,omitempty" bson:"token_name"`
	TokenSymbol string `json:"token_symbol,omitempty" bson:"token_symbol"`
	Decimals    int64  `json:"decimals,omitempty" bson:"decimals"`
	TotalSupply string `json:"total_supply" bson:"total_supply"`

	Contract   bool     `json:"contract" bson:"contract"`
	ErcTypes   []string `json:"erc_types" bson:"erc_types"`
	Interfaces []string `json:"interfaces" bson:"interfaces"`

	NumberOfTransactions         int `json:"number_of_transactions" bson:"number_of_transactions"`
	NumberOfTokenHolders         int `json:"number_of_token_holders,omitempty" bson:"number_of_token_holders"`
	NumberOfInternalTransactions int `json:"number_of_internal_transactions,omitempty" bson:"number_of_internal_transactions"`
	NumberOfTokenTransactions    int `json:"number_of_token_transactions,omitempty" bson:"number_of_token_transactions"`

	Target string `json:"target,omitempty" bson:"target"`
	Owner  string `json:"owner,omitempty" bson:"owner"`

	AttachedContract Contract `json:"attached_contract,omitempty" bson:"attached_contract,omitempty"`
}
