package models

import (
	"time"
)

type Block struct {
	Number     int       `json:"num" datastore:"num"`
	GasLimit   int       `json:"gas_limit" datastore:"gas_limit"`
	BlockHash  string    `json:"hash" datastore:"hash"`
	CreatedAt  time.Time `json:"created_at" datastore:"created_at"`
	ParentHash string    `json:"parent_hash" datastore:"parent_hash"`
	TxHash     string    `json:"tx_hash" datastore:"tx_hash"`
	GasUsed    string    `json:"gas_used" datastore:"gas_used"`
	Nonce      string    `json:"nonce" datastore:"nonce"`
	Miner      string    `json:"miner" datastore:"miner"`
	TxCount    int       `json:"tx_count" datastore:"tx_count"`
}

type BlockList struct {
	Blocks []*Block `json:"blocks" datastore:"blocks"`
}
