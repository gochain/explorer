package models

import (
	"github.com/gochain-io/gochain/v3/common"
)

type Signer struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Region string `json:"region"`
}

type BlockRange struct {
	StartBlock int64 `json:"start_block"`
	EndBlock   int64 `json:"end_block"`
}

type SignersStats struct {
	SignerStats map[common.Address]int64 `json:"signer_stats"`
	BlockRange  BlockRange               `json:"block_range"`
	Range       string                   `json:"range"`
}
