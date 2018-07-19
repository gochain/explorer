package models

import "time"

type Transaction struct {
	TxHash      string    `json:"tx_hash" firestore:"tx_hash"`
	To          string    `json:"to" firestore:"to"`
	From        string    `json:"from" firestore:"from"`
	Amount      string    `json:"amount" firestore:"amount"`
	Price       string    `json:"price" firestore:"price"`
	GasLimit    string    `json:"gas_limit" firestore:"gas_limit"`
	BlockNumber string    `json:"block_number" firestore:"block_number"`
	Nonce       string    `json:"nonce" firestore:"nonce"`
	BlockHash   string    `json:"block_hash" firestore:"hash"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
}

type TransactionList struct {
	Transactions []*Transaction `json:"transactions" firestore:"transactions"`
}
