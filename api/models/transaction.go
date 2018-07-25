package models

import "time"

type Transaction struct {
	TxHash      string    `json:"tx_hash" bson:"tx_hash"`
	To          string    `json:"to" bson:"to"`
	From        string    `json:"from" bson:"from"`
	Amount      int64     `json:"amount" bson:"amount"`
	Price       string    `json:"price" bson:"price"`
	GasLimit    string    `json:"gas_limit" bson:"gas_limit"`
	BlockNumber int64     `json:"block_number" bson:"block_number"`
	Nonce       string    `json:"nonce" bson:"nonce"`
	BlockHash   string    `json:"block_hash" bson:"hash"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

type TransactionList struct {
	Transactions []*Transaction `json:"transactions" bson:"transactions"`
}
