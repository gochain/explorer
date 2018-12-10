package models

import (
	"time"
)

type Block struct {
	Number          int64     `json:"number" bson:"number"`
	GasLimit        int       `json:"gas_limit" bson:"gas_limit"`
	BlockHash       string    `json:"hash" bson:"hash"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	ParentHash      string    `json:"parent_hash" bson:"parent_hash"`
	TxHash          string    `json:"tx_hash" bson:"tx_hash"`
	GasUsed         string    `json:"gas_used" bson:"gas_used"`
	Nonce           string    `json:"nonce" bson:"nonce"`
	Miner           string    `json:"miner" bson:"miner"`
	TxCount         int       `json:"tx_count" bson:"tx_count"`
	Difficulty      int64     `json:"difficulty" bson:"difficulty"`
	TotalDifficulty int64     `json:"total_difficulty" bson:"total_difficulty"`
	Sha3Uncles      string    `json:"sha3_uncles" bson:"sha3_uncles"`
	ExtraData       string    `json:"extra_data" bson:"extra_data"`
	// Transactions    []string  `json:"transactions" bson:"transactions"`
}

type LightBlock struct {
	Number    int64     `json:"number" bson:"number"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Miner     string    `json:"miner" bson:"miner"`
	TxCount   int       `json:"tx_count" bson:"tx_count"`
	ExtraData string    `json:"extra_data" bson:"extra_data"`
}

type LightBlockList struct {
	Blocks []*LightBlock `json:"blocks" bson:"blocks"`
}

type BlockList struct {
	Blocks []*Block `json:"blocks" bson:"blocks"`
}
