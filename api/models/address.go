package models

import "time"

type Address struct {
	Address       string    `json:"address" firestore:"address"`
	Owner         string    `json:"owner" firestore:"owner"`
	Balance       string    `json:"balance" firestore:"balance"`
	LastUpdatedAt time.Time `json:"last_updated_at" firestore:"last_updated_at"`
}

type AddressesList struct {
	Adresses []*Address `json:"adresses" firestore:"adresses"`
}
