package models

type Richlist struct {
	TotalSupply       string     `json:"total_supply"`
	CirculatingSupply string     `json:"circulating_supply"`
	Rankings          []*Address `json:"rankings"`
}
