package models

import "time"

type Transaction struct {
	TxHash          string    `json:"tx_hash" bson:"tx_hash"`
	To              string    `json:"to" bson:"to"`
	From            string    `json:"from" bson:"from"`
	Status          bool      `json:"status" bson:"status"`
	ContractAddress string    `json:"contract_address" bson:"contract_address"`
	Value           string    `json:"value" bson:"value"`
	GasPrice        string    `json:"gas_price" bson:"gas_price"`
	GasFee          string    `json:"gas_fee" bson:"gas_fee"`
	GasLimit        uint64    `json:"gas_limit" bson:"gas_limit"`
	BlockNumber     int64     `json:"block_number" bson:"block_number"`
	Nonce           uint64    `json:"nonce,string" bson:"nonce"`
	BlockHash       string    `json:"block_hash" bson:"hash"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	InputData       string    `json:"input_data" bson:"input_data"`
	Logs            string    `json:"logs" bson:"logs"`
	ReceiptReceived bool      `json:"-" bson:"receipt_received"`
}

type TransactionList struct {
	Transactions []*Transaction `json:"transactions"`
}
