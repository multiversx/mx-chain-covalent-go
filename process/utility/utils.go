package utility

import (
	"github.com/ElrondNetwork/elrond-go-core/core"
	"math/big"
)

func StrSliceToBytesSlice(in []string) [][]byte {
	out := make([][]byte, len(in))

	for i := range in {
		out[i] = make([]byte, len(in[i]))
		out[i] = []byte(in[i])
	}

	return out
}

func UIntSliceToIntSlice(in []uint64) []int64 {
	out := make([]int64, len(in))

	for i := range in {
		out[i] = int64(in[i])
	}

	return out
}

func GetBytes(val *big.Int) []byte {
	if val != nil {
		return val.Bytes()
	}

	return nil
}

func PubKeyToBech32(pubKeyConverter core.PubkeyConverter, pubKey []byte) []byte {
	return []byte(pubKeyConverter.Encode(pubKey))
}
