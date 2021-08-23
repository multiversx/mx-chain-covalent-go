package testscommon

import (
	"math/rand"
	"strconv"
)

func GenerateRandomBytes() []byte {
	return []byte(strconv.Itoa(rand.Int()))
}
