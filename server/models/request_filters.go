package models

import (
	"github.com/gochain-io/explorer/server/utils"
	"time"
)

type DefaultFilter struct {
	Skip  int `schema:"skip,omitempty"`
	Limit int `schema:"limit,omitempty"`
}

type SortFilter struct {
	SortBy string `schema:"sortby,omitempty"`
	Asc    bool   `schema:"asc,omitempty"`
}

type ContractsFilter struct {
	DefaultFilter
	SortFilter
	ContractName string        `schema:"contract_name,omitempty"`
	TokenName    string        `schema:"token_name,omitempty"`
	TokenSymbol  string        `schema:"token_symbol,omitempty"`
	ErcType      utils.ErcName `schema:"erc_type,omitempty"`
}

type InternalTxFilter struct {
	DefaultFilter
	TokenTransactions bool `schema:"token_transactions,omitempty"`
}

type TxsFilter struct {
	DefaultFilter
	InputDataEmpty *bool     `schema:"input_data_empty,omitempty"`
	FromTime       time.Time `schema:"input_data_empty,omitempty"`
	ToTime         time.Time `schema:"to_time,omitempty"`
}
