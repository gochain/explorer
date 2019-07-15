package models

import (
	"github.com/gochain-io/gochain/v3/common"
)

type SignerStats struct {
	Name        string
	URL         string
	Region      string
	Signer      common.Address `json:"signer"`
	BlocksCount int            `json:"blocks_count"`
}
type BlockRange struct {
	StartBlock, EndBlock int64
}

type SignersStats struct {
	SignerStats []SignerStats `json:"signer_stats"`
	BlockRange  BlockRange    `json:"block_range"`
	Range       string        `json:"range"`
}
