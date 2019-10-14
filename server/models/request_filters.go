package models

import (
	"github.com/gochain-io/explorer/server/utils"
	"time"
)

type PaginationFilter struct {
	Skip  int `schema:"skip,omitempty"`
	Limit int `schema:"limit,omitempty"`
}

func (f *PaginationFilter) ProcessPagination() {
	if f.Limit == 0 || f.Limit > utils.MaxFetchLimit {
		f.Limit = utils.MaxFetchLimit
	}
}

type SortFilter struct {
	SortBy string `schema:"sortby,omitempty"`
	Asc    bool   `schema:"asc,omitempty"`
}

type TimeFilter struct {
	FromTime time.Time `schema:"from_time,omitempty"`
	ToTime   time.Time `schema:"to_time,omitempty"`
}

func (f *TimeFilter) ProcessTime() {
	if f.FromTime.IsZero() {
		f.FromTime = time.Unix(0, 0)
	}
	if f.ToTime.IsZero() {
		f.ToTime = time.Now()
	}
}

type ContractsFilter struct {
	PaginationFilter
	SortFilter
	ContractName string `schema:"contract_name,omitempty"`
	TokenName    string `schema:"token_name,omitempty"`
	TokenSymbol  string `schema:"token_symbol,omitempty"`
	ErcType      string `schema:"erc_type,omitempty"`
}

type InternalTxFilter struct {
	PaginationFilter
	TokenTransactions bool `schema:"token_transactions,omitempty"`
}

type TxsFilter struct {
	PaginationFilter
	TimeFilter
	InputDataEmpty *bool `schema:"input_data_empty,omitempty"`
}
