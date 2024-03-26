package tokens

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"

	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/core/types"
)

func Test_unpackERC20TransferEvent(t *testing.T) {
	from := common.BigToAddress(big.NewInt(rand.Int63()))
	to := common.BigToAddress(big.NewInt(rand.Int63()))
	data := common.BigToHash(big.NewInt(1)).Bytes()
	hash := common.BigToHash(big.NewInt(rand.Int63()))
	tests := []struct {
		name    string
		log     types.Log
		want    *TransferEvent
		wantErr bool
	}{
		{"valid", types.Log{Topics: []common.Hash{{}, from.Hash(), to.Hash()}, Data: data, BlockNumber: 1, TxHash: hash},
			&TransferEvent{From: from, To: to, Value: big.NewInt(1), BlockNumber: 1, TransactionHash: hash.Hex()}, false},
		{"2topics", types.Log{Topics: make([]common.Hash, 2), Data: data}, nil, true},
		{"4topics", types.Log{Topics: make([]common.Hash, 4), Data: data}, nil, true},
		{"nilData", types.Log{Topics: make([]common.Hash, 3), Data: nil}, nil, true},
		{"emptyData", types.Log{Topics: make([]common.Hash, 3), Data: []byte{}}, nil, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := unpackERC20TransferEvent(test.log)
			if (err != nil) != test.wantErr {
				t.Errorf("error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want && ((got == nil || test.want == nil) || !reflect.DeepEqual(*got, *test.want)) {
				t.Errorf("\ngot: %#v\nwant: %#v", *got, *test.want)
			}
		})
	}
}

func Test_unpackERC721TransferEvent(t *testing.T) {
	from := common.BigToAddress(big.NewInt(rand.Int63()))
	to := common.BigToAddress(big.NewInt(rand.Int63()))
	data := common.BigToHash(big.NewInt(1)).Bytes()
	hash := common.BigToHash(big.NewInt(rand.Int63()))
	tests := []struct {
		name    string
		log     types.Log
		want    *TransferEvent
		wantErr bool
	}{
		{"valid", types.Log{Topics: []common.Hash{{}, from.Hash(), to.Hash(), common.BytesToHash(big.NewInt(1).Bytes())}, BlockNumber: 1, TxHash: hash},
			&TransferEvent{From: from, To: to, Value: big.NewInt(1), BlockNumber: 1, TransactionHash: hash.Hex()}, false},
		{"3topics", types.Log{Topics: make([]common.Hash, 3)}, nil, true},
		{"5topics", types.Log{Topics: make([]common.Hash, 5)}, nil, true},
		{"data", types.Log{Topics: make([]common.Hash, 4), Data: data}, nil, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := unpackERC721TransferEvent(test.log)
			if (err != nil) != test.wantErr {
				t.Errorf("error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want && ((got == nil || test.want == nil) || !reflect.DeepEqual(*got, *test.want)) {
				t.Errorf("\ngot: %#v\nwant: %#v", *got, *test.want)
			}
		})
	}
}
