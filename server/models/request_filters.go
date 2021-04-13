package models

import (
	"time"
)

const (
	defaultLimit = 50
	maximumLimit = 500
)

type PaginationFilter struct {
	Skip  int `schema:"skip,omitempty"`
	Limit int `schema:"limit,omitempty"`
}

func (f *PaginationFilter) Sanitize() {
	if f.Skip < 0 {
		f.Skip = 0
	}
	if f.Limit <= 0 {
		f.Limit = defaultLimit
	} else if f.Limit > maximumLimit {
		f.Limit = maximumLimit
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

func (f *TimeFilter) Sanitize() {
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
	TokenTransactions bool   `schema:"token_transactions,omitempty"`
	InternalAddress   string `schema:"internal_address,omitempty"`
	TokenID           string `schema:"token_id,omitempty"`
}

type TxsFilter struct {
	PaginationFilter
	TimeFilter
}

func (f *TxsFilter) Sanitize() {
	f.PaginationFilter.Sanitize()
	f.TimeFilter.Sanitize()
}
