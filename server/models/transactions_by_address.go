package models

import "time"

type TransactionsByAddress struct {
	TxHash    string    `json:"tx_hash" bson:"tx_hash"`
	Address   string    `json:"address" bson:"address"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
