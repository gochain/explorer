package models

import (
	"github.com/gochain-io/gochain/v3/accounts/abi"
	"time"
)

type AbiItem struct {
	Anonymous       bool           `json:"anonymous" bson:"anonymous"`
	Constant        bool           `json:"constant" bson:"constant"`
	Inputs          []abi.Argument `json:"inputs" bson:"inputs"`
	Name            string         `json:"name" bson:"name"`
	Outputs         []abi.Argument `json:"outputs" bson:"outputs"`
	Payable         bool           `json:"payable" bson:"payable"`
	StateMutability string         `json:"stateMutability" bson:"stateMutability"`
	Type            string         `json:"type" bson:"type"`
}

type Contract struct {
	Address         string    `json:"address" bson:"address"`
	Bytecode        string    `json:"byte_code" bson:"byte_code,omitempty"`
	Valid           bool      `json:"valid" bson:"valid,omitempty"`
	ContractName    string    `json:"contract_name" bson:"contract_name,omitempty"`
	CompilerVersion string    `json:"compiler_version" bson:"compiler_version,omitempty"`
	Optimization    bool      `json:"optimization" bson:"optimization,omitempty"`
	SourceCode      string    `json:"source_code" bson:"source_code,omitempty"`
	Abi             []AbiItem `json:"abi" bson:"abi,omitempty"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at,omitempty"`
	/*RecaptchaToken  string    `json:"recaptcha_token"`*/
}
