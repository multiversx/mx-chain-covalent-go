package utility

import (
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
)

// UIntSliceToIntSlice outputs the int64 slice representation of a uint64 slice input
func UIntSliceToIntSlice(in []uint64) []int64 {
	out := make([]int64, len(in))

	for i := range in {
		out[i] = int64(in[i])
	}

	return out
}

// UInt32SliceToInt32Slice outputs the int32 slice representation of a uint32 slice input
func UInt32SliceToInt32Slice(in []uint32) []int32 {
	out := make([]int32, len(in))

	for i := range in {
		out[i] = int32(in[i])
	}

	return out
}

// GetBytes returns the bytes representation of a big int input if not nil, otherwise returns []byte{}
func GetBytes(val *big.Int) []byte {
	if val != nil {
		return val.Bytes()
	}

	return big.NewInt(0).Bytes()
}

// GetBigIntBytesFromStr returns the big int bytes representation of a string input if not empty, otherwise returns []byte{}
func GetBigIntBytesFromStr(val string) ([]byte, error) {
	if len(val) == 0 {
		return big.NewInt(0).Bytes(), nil
	}

	valBI, ok := big.NewInt(0).SetString(val, 10)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errInvalidValueInBase10, val)
	}

	return valBI.Bytes(), nil
}

// StringSliceToByteSlice converts the input string slice to a byte slice
func StringSliceToByteSlice(in []string) [][]byte {
	out := make([][]byte, len(in))

	for idx, elem := range in {
		out[idx] = []byte(elem)
	}

	return out
}

// GetBigIntBytesSliceFromStringSlice converts the input string slice in a big int byte array slice
func GetBigIntBytesSliceFromStringSlice(in []string) ([][]byte, error) {
	out := make([][]byte, len(in))

	for idx, elem := range in {
		bigIntBytes, err := GetBigIntBytesFromStr(elem)
		if err != nil {
			return nil, err
		}

		out[idx] = bigIntBytes
	}

	return out, nil
}

// GetAddressOrMetachainAddr checks if the corresponding address is metachain. This func should only be used for sender addresses.
// If so, it returns a 62 byte array address(by padding with zeros), otherwise converts the address string to byte slice.
func GetAddressOrMetachainAddr(address string) []byte {
	if address != MetachainShardName {
		return []byte(address)
	}

	return MetaChainShardAddress()
}

// MetaChainShardAddress returns core.MetachainShardId as a 62 byte array address(by padding with zeros).
// This is needed, since all addresses from avro schema are required to be 62 fixed byte array
func MetaChainShardAddress() []byte {
	ret := make([]byte, 62)
	copy(ret, fmt.Sprintf("%d", core.MetachainShardId))
	return ret
}
