package models

import "time"

type Address struct {
	Address       string    `json:"address" bson:"address"`
	Owner         string    `json:"owner" bson:"owner"`
	Balance       string    `json:"balance" bson:"balance"`
	BalanceInt    int64     `json:"balance_int" bson:"balance_int"`
	LastUpdatedAt time.Time `json:"last_updated_at" bson:"last_updated_at"`
	TokenName     string    `json:"token_name" bson:"token_name"`
	TokenSymbol   string    `json:"token_symbol" bson:"token_symbol"`
	Contract      bool      `json:"contract" bson:"contract"`
	GO20          bool      `json:"go20" bson:"go20"`
}

type AddressesList struct {
	Adresses []*Address `json:"adresses"`
}
