package models

import "time"

type Address struct {
	Address       string    `json:"address" datastore:"address"`
	Owner         string    `json:"owner" datastore:"owner"`
	From          string    `json:"from" datastore:"from"`
	Balance       string    `json:"balance" datastore:"balance"`
	LastUpdatedAt time.Time `json:"last_updated_at" datastore:"last_updated_at"`
}

type AddressesList struct {
	Adresses []*Address `json:"adresses" datastore:"adresses"`
}
