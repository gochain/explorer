package models

import (
	"encoding/json"

	"github.com/gochain-io/gochain/v3/common"
)

type Node struct {
	Address common.Address
	Name    string
	URL     string
	Region  string
}

func ParseConfig(data []byte) (map[common.Address]Node, error) {
	var nodes []Node
	m := make(map[common.Address]Node)
	err := json.Unmarshal(data, &nodes)
	if err != nil {
		return nil, err
	}
	for _, l := range nodes {
		m[l.Address] = l
	}
	return m, nil
}
