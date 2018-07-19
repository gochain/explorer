package models

import (
	"time"
)

type Block struct {
	Number     int       `json:"num" firestore:"num"`
	GasLimit   int       `json:"gas_limit" firestore:"gas_limit"`
	BlockHash  string    `json:"hash" firestore:"hash"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
	ParentHash string    `json:"parent_hash" firestore:"parent_hash"`
	TxHash     string    `json:"tx_hash" firestore:"tx_hash"`
	GasUsed    string    `json:"gas_used" firestore:"gas_used"`
	Nonce      string    `json:"nonce" firestore:"nonce"`
	Miner      string    `json:"miner" firestore:"miner"`
	TxCount    int       `json:"tx_count" firestore:"tx_count"`
}

type BlockList struct {
	Blocks []*Block `json:"blocks" firestore:"blocks"`
}
