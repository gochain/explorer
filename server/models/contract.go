package models

type Contract struct {
	Address         string `json:"address" bson:"address"`
	Bytecode        string `json:"byteCode" bson:"byteCode"`
	Valid           bool   `json:"valid" bson:"valid"`
	ContractName    string `json:"contractName" bson:"contractName"`
	CompilerVersion string `json:"compilerVersion" bson:"compilerVersion"`
	Optimization    bool   `json:"optimization" bson:"optimization"`
	SourceCode      string `json:"sourceCode" bson:"sourceCode"`
	Abi             string `json:"abi" bson:"abi"`
}
