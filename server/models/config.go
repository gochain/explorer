package models

import "github.com/gochain-io/gochain/v3/common"

type Node struct {
	Address common.Address
	Name    string
	URL     string
	Region  string
}
