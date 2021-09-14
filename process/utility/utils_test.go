package utility_test

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/stretchr/testify/require"
)

func TestStrSliceToBytesSlice_DifferentValues(t *testing.T) {
	in := []string{"a", "b", "c"}
	out := utility.StrSliceToBytesSlice(in)

	require.Len(t, out, 3)
	require.Equal(t, []byte("a"), out[0])
	require.Equal(t, []byte("b"), out[1])
	require.Equal(t, []byte("c"), out[2])
}

func TestStrSliceToBytesSlice_EmptyInput(t *testing.T) {
	out := utility.StrSliceToBytesSlice([]string{})

	require.Len(t, out, 0)
}

func TestUIntSliceToIntSlice_DifferentValues(t *testing.T) {
	in := []uint64{1, 2, 3}
	out := utility.UIntSliceToIntSlice(in)

	require.Len(t, out, 3)
	require.Equal(t, int64(1), out[0])
	require.Equal(t, int64(2), out[1])
	require.Equal(t, int64(3), out[2])
}

func TestUIntSliceToIntSlice_EmptyInput(t *testing.T) {
	out := utility.UIntSliceToIntSlice([]uint64{})

	require.Len(t, out, 0)
}

func TestGetBytes(t *testing.T) {
	require.Nil(t, utility.GetBytes(nil))

	x := big.NewInt(10)
	require.Equal(t, []byte{0xa}, utility.GetBytes(x))
}

func TestEncodeDecode(t *testing.T) {
	account := &schema.AccountBalanceUpdate{
		Address: testscommon.GenerateRandomFixedBytes(62),
		Balance: big.NewInt(1000).Bytes(),
		Nonce:   444,
	}

	buffer, err := utility.Encode(account)
	require.Nil(t, err)

	decodedAccount := &schema.AccountBalanceUpdate{}
	err = utility.Decode(decodedAccount, buffer)
	require.Nil(t, err)

	require.Equal(t, account, decodedAccount)
}
