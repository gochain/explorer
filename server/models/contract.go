package models

import (
	"github.com/gochain-io/explorer/server/utils"
	"time"
)

type Contract struct {
	Address         string          `json:"address" bson:"address"`
	Bytecode        string          `json:"byte_code" bson:"byte_code,omitempty"`
	Valid           bool            `json:"valid" bson:"valid,omitempty"`
	ContractName    string          `json:"contract_name" bson:"contract_name,omitempty"`
	CompilerVersion string          `json:"compiler_version" bson:"compiler_version,omitempty"`
	Optimization    bool            `json:"optimization" bson:"optimization,omitempty"`
	SourceCode      string          `json:"source_code" bson:"source_code,omitempty"`
	CreatedAt       time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" bson:"updated_at"`
	Abi             []utils.AbiItem `json:"abi" bson:"abi"`
	/*RecaptchaToken  string    `json:"recaptcha_token"`*/
}
