package storage

import (
	"crypto/sha256"
	"fmt"

	"github.com/AccumulateNetwork/accumulated/smt/common"
	"github.com/AccumulateNetwork/accumulated/types"
)

const (
	KeyLength = 32 // Total bytes used for keys
)

// var ElementIndexKey = []byte("ElementIndex")
// var BptKey = []byte("BPT")
// var BptRootKey = []byte("BPT.Root")
// var StatesKey = []byte("States")
// var NextElementKey = []byte("NextElement")
// var StatesHeadKey = []byte("States.Head")

type Key = [KeyLength]byte

func ComputeKey(keys ...interface{}) Key {
	n := -1
	bkeys := make([][]byte, len(keys))
	for i, key := range keys {
		switch key := key.(type) {
		case nil:
			bkeys[i] = []byte{}
		case []byte:
			bkeys[i] = key
		case [32]uint8:
			bkeys[i] = key[:]
		case types.Bytes:
			bkeys[i] = key
		case types.Bytes32:
			bkeys[i] = key[:]
		case string:
			bkeys[i] = []byte(key)
		case types.String:
			bkeys[i] = []byte(key)
		case uint:
			bkeys[i] = common.Uint64Bytes(uint64(key))
		case uint8:
			bkeys[i] = common.Uint64Bytes(uint64(key))
		case uint16:
			bkeys[i] = common.Uint64Bytes(uint64(key))
		case uint32:
			bkeys[i] = common.Uint64Bytes(uint64(key))
		case uint64:
			bkeys[i] = common.Uint64Bytes(key)
		case int:
			bkeys[i] = common.Int64Bytes(int64(key))
		case int8:
			bkeys[i] = common.Int64Bytes(int64(key))
		case int16:
			bkeys[i] = common.Int64Bytes(int64(key))
		case int32:
			bkeys[i] = common.Int64Bytes(int64(key))
		case int64:
			bkeys[i] = common.Int64Bytes(key)
		default:
			panic(fmt.Errorf("cannot use %T as a key", key))
		}

		n += len(bkeys[i]) + 1
	}

	composite := make([]byte, 0, n)
	for i, k := range bkeys {
		if i > 0 {
			composite = append(composite, '.')
		}
		composite = append(composite, k...)
	}

	return sha256.Sum256(composite)
}