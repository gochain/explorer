package models

import "time"

type Transaction struct {
	TxHash      string    `json:"tx_hash" datastore:"tx_hash"`
	To          string    `json:"to" datastore:"to"`
	From        string    `json:"from" datastore:"from"`
	Amount      string    `json:"amount" datastore:"amount"`
	Price       string    `json:"price" datastore:"price"`
	GasLimit    string    `json:"gas_limit" datastore:"gas_limit"`
	BlockNumber string    `json:"block_number" datastore:"block_number"`
	Nonce       string    `json:"nonce" datastore:"nonce"`
	BlockHash   string    `json:"block_hash" datastore:"hash"`
	CreatedAt   time.Time `json:"created_at" datastore:"created_at"`
}

type TransactionList struct {
	Transactions []*Transaction `json:"transactions" datastore:"transactions"`
}
