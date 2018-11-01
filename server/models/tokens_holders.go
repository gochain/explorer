package models

import "time"

type TokenHolder struct {
	ContractAddress    string    `json:"contract_address" bson:"contract_address"`
	TokenHolderAddress string    `json:"token_holder_address" bson:"token_holder_address"`
	Balance            string    `json:"balance" bson:"balance"`
	BalanceInt         int64     `json:"balance_int" bson:"balance_int"`
	UpdatedAt          time.Time `json:"updated_at" bson:"updated_at"`
}

type TokenHolderList struct {
	Holders []*TokenHolder `json:"token_holders"`
}
