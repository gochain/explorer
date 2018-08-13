package models

type Stats struct {
	NumberOfTransactions int64 `json:"total_transactions_count"`
	NumberOfBlocks       int64 `json:"total_blocks_count"`
}
type Richlist struct {
	TotalSupply       string     `json:"total_supply"`
	CirculatingSupply string     `json:"circulating_supply"`
	Rankings          []*Address `json:"rankings"`
}
