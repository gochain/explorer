package models

import (
	"time"
)

type Block struct {
	Number    int       `json:"num" datastore:"num"`
	GasLimit  int       `json:"gas_limit" datastore:"gas_limit"`
	Hash      string    `json:"hash" datastore:"hash"`
	CreatedAt time.Time `json:"created_at" datastore:"created_at"`
}

type BlockList struct {
	Blocks []*Block `json:"blocks" datastore:"blocks"`
}
