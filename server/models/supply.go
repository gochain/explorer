package models

import (
	"encoding/json"
	"math/big"
)

type SupplyStats struct {
	Total       *big.Int
	Circulating *big.Int
	Locked      *big.Int
	FeesBurned  *big.Int
}

func (s *SupplyStats) MarshalJSON() ([]byte, error) {
	var encoded = struct {
		Total       string `json:"total,omitempty"`
		Circulating string `json:"circulating,omitempty"`
		Locked      string `json:"locked,omitempty"`
		FeesBurned  string `json:"fees_burned,omitempty"`
	}{}
	if s.Total != nil {
		encoded.Total = s.Total.String()
	}
	if s.Circulating != nil {
		encoded.Circulating = s.Circulating.String()
	}
	if s.Locked != nil {
		encoded.Locked = s.Locked.String()
	}
	if s.FeesBurned != nil {
		encoded.FeesBurned = s.FeesBurned.String()
	}
	return json.Marshal(encoded)
}
