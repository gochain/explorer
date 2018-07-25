package models

import (
	"time"
)

type Block struct {
	Number     int64     `json:"number" bson:"num"`
	GasLimit   int       `json:"gas_limit" bson:"gas_limit"`
	BlockHash  string    `json:"hash" bson:"hash"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	ParentHash string    `json:"parent_hash" bson:"parent_hash"`
	TxHash     string    `json:"tx_hash" bson:"tx_hash"`
	GasUsed    string    `json:"gas_used" bson:"gas_used"`
	Nonce      string    `json:"nonce" bson:"nonce"`
	Miner      string    `json:"miner" bson:"miner"`
	TxCount    int       `json:"tx_count" bson:"tx_count"`
}

type BlockList struct {
	Blocks []*Block `json:"blocks" bson:"blocks"`
}
