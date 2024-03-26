package models

import (
	"github.com/gochain/gochain/v4/common"
)

type Signer struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Region string `json:"region"`
}

type SignerStats struct {
	SignerAddress common.Address `json:"signer_address"`
	BlocksCount   int            `json:"blocks_count"`
}

type BlockRange struct {
	StartBlock int64 `json:"start_block"`
	EndBlock   int64 `json:"end_block"`
}

type SignersStats struct {
	// front needs arr here
	SignerStats []SignerStats `json:"signer_stats"`
	BlockRange  BlockRange    `json:"block_range"`
	Range       string        `json:"range"`
}
