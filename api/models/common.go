package models

type Stats struct {
	NumberOfTransactions int64 `json:"total_transactions_count"`
	NumberOfBlocks       int64 `json:"total_blocks_count"`
}
type Richlist struct {
	TotalSupply       int64      `json:"total_supply"`
	CirculatingSupply int64      `json:"circulating_supply"`
	Rankings          []*Address `json:"rankings"`
}
