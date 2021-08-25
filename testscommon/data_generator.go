package testscommon

import (
	"math/big"
	"math/rand"
	"strconv"
)

func GenerateRandomFixedBytes(n int) []byte {
	ret := make([]byte, n)

	for i := 0; i < n; i++ {
		ret[i] = byte(rand.Int())
	}

	return ret
}

func GenerateRandomBytes() []byte {
	return []byte(strconv.Itoa(rand.Int()))
}

func GenerateRandomBigInt() *big.Int {
	return big.NewInt(rand.Int63())
}
