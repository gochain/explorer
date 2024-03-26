package utils

import (
	"bytes"
	"testing"

	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/crypto"
)

func TestInterfaceIdentifiers(t *testing.T) {
	for name, data := range EVMFunctions {
		t.Run(name.String(), func(t *testing.T) {
			got := crypto.Keccak256Hash([]byte(data.Signature)).Bytes()[:4]
			if !bytes.Equal(common.Hex2Bytes(data.ID), got) {
				t.Errorf("constant %x but bytes4(keccak256(sig)): %x", data.ID, got)
			}
		})
	}
}
