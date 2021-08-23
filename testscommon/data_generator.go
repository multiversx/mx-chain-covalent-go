package testscommon

import (
	"math/big"
	"math/rand"
	"strconv"
)

func GenerateRandomBytes() []byte {
	return []byte(strconv.Itoa(rand.Int()))
}

func GenerateRandomBigInt() *big.Int {
	return big.NewInt(rand.Int63())
}
