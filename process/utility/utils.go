package utility

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
)

// HexSliceToByteSlice outputs a decoded byte slice representation of a hex string encoded slice input
func HexSliceToByteSlice(in []string) ([][]byte, error) {
	if in == nil {
		return nil, nil
	}
	out := make([][]byte, len(in))

	for i := range in {
		tmp, err := hex.DecodeString(in[i])
		if err != nil {
			return nil, err
		}
		out[i] = tmp
	}

	return out, nil
}

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

// EncodePubKey returns a byte slice of the encoded pubKey input, using a pub key converter
func EncodePubKey(pubKeyConverter core.PubkeyConverter, pubKey []byte) []byte {
	return []byte(pubKeyConverter.Encode(pubKey))
}

// MetaChainShardAddress returns core.MetachainShardId as a 62 byte array address(by padding with zeros).
// This is needed, since all addresses from avro schema are required to be 62 fixed byte array
func MetaChainShardAddress() []byte {
	ret := make([]byte, 62)
	copy(ret, fmt.Sprintf("%d", core.MetachainShardId))
	return ret
}
