package models

import "math/big"

type SupplyStats struct {
	Total       *big.Int `json:"total"`
	Circulating *big.Int `json:"circulating"`
	Locked      *big.Int `json:"locked"`
	FeesBurned  *big.Int `json:"fees_burned"`
}
