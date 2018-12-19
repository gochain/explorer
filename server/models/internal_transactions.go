package models

import "time"

type InternalTransaction struct {
	Contract        string    `json:"contract_address" bson:"contract_address"`
	From            string    `json:"from_address" bson:"from_address"`
	To              string    `json:"to_address" bson:"to_address"`
	Value           string    `json:"value" bson:"value"`
	BlockNumber     int64     `json:"block_number" bson:"block_number"`
	TransactionHash string    `json:"transaction_hash" bson:"transaction_hash"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
}

type InternalTransactionsList struct {
	Transactions []*InternalTransaction `json:"internal_transactions"`
}
