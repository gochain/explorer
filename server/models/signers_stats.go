package models

import (
	"github.com/gochain-io/gochain/v3/common"
)

type Signer struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Region string `json:"region"`
}

type SignerStats struct {
	Signer        Signer         `json:"signer"`
	SignerAddress common.Address `json:"signer_address"`
	BlocksCount   int            `json:"blocks_count"`
}
type BlockRange struct {
	StartBlock, EndBlock int64
}

type SignersStats struct {
	SignerStats []SignerStats `json:"signer_stats"`
	BlockRange  BlockRange    `json:"block_range"`
	Range       string        `json:"range"`
}
