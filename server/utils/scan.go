package utils

import (
	"strings"
	"sync"

	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/crypto"
)

var contractInfoCache = struct {
	sync.RWMutex
	data map[common.Hash]contractInfo
}{data: make(map[common.Hash]contractInfo)}

type contractInfo struct {
	types map[EVMInterface]struct{}
	funcs map[EVMFunction]struct{}
}

// implementsAll returns true if funcs contains all of the EVMFunctions.
func (ci *contractInfo) implementsAll(fns []EVMFunction) bool {
	for _, fn := range fns {
		if _, ok := ci.funcs[fn]; !ok {
			return false
		}
	}
	return true
}

// scan scans the hex byte code string for function IDs to initialize funcs and types.
func (ci *contractInfo) scan(byteCode string) {
	if len(byteCode) < 8 {
		return
	}
	ci.funcs = map[EVMFunction]struct{}{}
	for k, v := range EVMFunctions {
		if strings.Contains(byteCode, v.ID) {
			ci.funcs[k] = struct{}{}
		}
	}
	ci.types = map[EVMInterface]struct{}{}
	for i, fns := range EVMFunctionsByInterface {
		if !ci.implementsAll(fns) {
			continue
		}
		ci.types[EVMInterface(i)] = struct{}{}
	}
}

// ScanContract returns the interface and function sets for the hex encoded byte code.
// Results are cached, and the final return is the number of cached entries.
func ScanContract(byteCode string) (map[EVMInterface]struct{}, map[EVMFunction]struct{}, int) {
	hash := crypto.Keccak256Hash([]byte(byteCode))
	contractInfoCache.RLock()
	ci, ok := contractInfoCache.data[hash]
	l := len(contractInfoCache.data)
	contractInfoCache.RUnlock()
	if ok {
		return ci.types, ci.funcs, l
	}

	ci.scan(byteCode)

	contractInfoCache.Lock()
	contractInfoCache.data[hash] = ci
	l = len(contractInfoCache.data)
	contractInfoCache.Unlock()

	return ci.types, ci.funcs, l
}
