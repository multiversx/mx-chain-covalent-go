package utility

import (
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestStrSliceToBytesSlice_DifferentValues(t *testing.T) {
	in := []string{"a", "b", "c"}
	out := StrSliceToBytesSlice(in)

	require.Equal(t, len(out), 3)
	require.Equal(t, out[0], []byte("a"))
	require.Equal(t, out[1], []byte("b"))
	require.Equal(t, out[2], []byte("c"))
}

func TestStrSliceToBytesSlice_EmptyInput(t *testing.T) {
	out := StrSliceToBytesSlice([]string{})

	require.Equal(t, len(out), 0)
}

func TestUIntSliceToIntSlice_DifferentValues(t *testing.T) {
	in := []uint64{1, 2, 3}
	out := UIntSliceToIntSlice(in)

	require.Equal(t, len(out), 3)
	require.Equal(t, out[0], int64(1))
	require.Equal(t, out[1], int64(2))
	require.Equal(t, out[2], int64(3))
}

func TestUIntSliceToIntSlice_EmptyInput(t *testing.T) {
	out := UIntSliceToIntSlice([]uint64{})

	require.Equal(t, len(out), 0)
}

func TestGetBytes(t *testing.T) {
	require.Nil(t, GetBytes(nil))

	x := big.NewInt(10)
	require.Equal(t, GetBytes(x), []byte{0xa})
}
