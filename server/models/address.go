package models

import "time"

type Address struct {
	Address                      string    `json:"address" bson:"address"`
	Owner                        string    `json:"owner,omitempty" bson:"owner"`
	Balance                      float64   `json:"balance,string" bson:"balance"` //backward compability
	BalanceWei                   string    `json:"balance_wei" bson:"balance_wei"`
	LastUpdatedAt                time.Time `json:"last_updated_at" bson:"last_updated_at"`
	TokenName                    string    `json:"token_name,omitempty" bson:"token_name"`
	TokenSymbol                  string    `json:"token_symbol,omitempty" bson:"token_symbol"`
	Contract                     bool      `json:"contract" bson:"contract"`
	GO20                         bool      `json:"go20" bson:"go20"`
	NumberOfTransactions         int       `json:"number_of_transactions" bson:"number_of_transactions"`
	NumberOfTokenHolders         int       `json:"number_of_token_holders,omitempty" bson:"number_of_token_holders"`
	NumberOfInternalTransactions int       `json:"number_of_internal_transactions,omitempty" bson:"number_of_internal_transactions"`
}

type AddressesList struct {
	Adresses []*Address `json:"adresses"`
}
