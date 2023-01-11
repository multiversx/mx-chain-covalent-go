package utility_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/stretchr/testify/require"
)

var testAvroMarshaller = &utility.AvroMarshaller{}

func TestUIntSliceToIntSlice_DifferentValues(t *testing.T) {
	t.Parallel()

	in := []uint64{1, 2, 3}
	out := utility.UIntSliceToIntSlice(in)

	require.Len(t, out, 3)
	require.Equal(t, int64(1), out[0])
	require.Equal(t, int64(2), out[1])
	require.Equal(t, int64(3), out[2])
}

func TestUIntSliceToIntSlice_EmptyInput(t *testing.T) {
	t.Parallel()

	out := utility.UIntSliceToIntSlice([]uint64{})
	require.Len(t, out, 0)
}

func TestGetBytes(t *testing.T) {
	t.Parallel()

	require.Equal(t, utility.GetBytes(nil), []byte{})

	x := big.NewInt(10)
	require.Equal(t, []byte{0xa}, utility.GetBytes(x))
}

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	account := &schema.AccountBalanceUpdate{
		Address: testscommon.GenerateRandomFixedBytes(62),
		Balance: big.NewInt(1000).Bytes(),
		Nonce:   444,
	}

	buffer, err := testAvroMarshaller.Encode(account)
	require.Nil(t, err)

	decodedAccount := &schema.AccountBalanceUpdate{}
	err = testAvroMarshaller.Decode(decodedAccount, buffer)
	require.Nil(t, err)

	require.Equal(t, account, decodedAccount)
}

func TestEncode_HyperBlock(t *testing.T) {
	t.Parallel()

	block := schema.HyperBlock{
		Hash:          testscommon.GenerateRandomFixedBytes(32),
		StateRootHash: testscommon.GenerateRandomFixedBytes(32),
	}
	_, err := testAvroMarshaller.Encode(&block)
	require.Nil(t, err)

	blockNilHash := block
	blockNilHash.Hash = nil
	_, err = testAvroMarshaller.Encode(&blockNilHash)
	require.NotNil(t, err)

	blockNilStateRootHash := block
	blockNilStateRootHash.StateRootHash = nil
	_, err = testAvroMarshaller.Encode(&blockNilStateRootHash)
	require.Nil(t, err)

	blockNilPrevRootHash := block
	blockNilPrevRootHash.PrevBlockHash = nil
	_, err = testAvroMarshaller.Encode(&blockNilPrevRootHash)
	require.Nil(t, err)
}

func TestEncode_EpochStartInfo(t *testing.T) {
	t.Parallel()

	info := schema.EpochStartInfo{}
	_, err := testAvroMarshaller.Encode(&info)
	require.Nil(t, err)
}

func TestEncode_Transaction(t *testing.T) {
	t.Parallel()

	tx := schema.Transaction{
		Hash:           testscommon.GenerateRandomFixedBytes(32),
		MiniBlockHash:  testscommon.GenerateRandomFixedBytes(32),
		BlockHash:      testscommon.GenerateRandomFixedBytes(32),
		Receiver:       testscommon.GenerateRandomFixedBytes(62),
		Sender:         testscommon.GenerateRandomFixedBytes(62),
		HyperBlockHash: testscommon.GenerateRandomFixedBytes(32),
	}
	_, err := testAvroMarshaller.Encode(&tx)
	require.Nil(t, err)

	txNilHash := tx
	txNilHash.Hash = nil
	_, err = testAvroMarshaller.Encode(&txNilHash)
	require.NotNil(t, err)

	txNilMiniBlockHash := tx
	txNilMiniBlockHash.MiniBlockHash = nil
	_, err = testAvroMarshaller.Encode(&txNilMiniBlockHash)
	require.NotNil(t, err)

	txNilBlockHash := tx
	txNilBlockHash.BlockHash = nil
	_, err = testAvroMarshaller.Encode(&txNilBlockHash)
	require.Nil(t, err)

	txNilReceiver := tx
	txNilReceiver.Receiver = nil
	_, err = testAvroMarshaller.Encode(&txNilReceiver)
	require.NotNil(t, err)

	txNilSender := tx
	txNilSender.Sender = nil
	_, err = testAvroMarshaller.Encode(&txNilSender)
	require.NotNil(t, err)
}

func TestEncode_Receipt(t *testing.T) {
	t.Parallel()

	receipt := schema.Receipt{
		Sender: testscommon.GenerateRandomFixedBytes(62),
		TxHash: testscommon.GenerateRandomFixedBytes(32),
	}

	_, err := testAvroMarshaller.Encode(&receipt)
	require.Nil(t, err)

	receiptNilSender := receipt
	receiptNilSender.Sender = nil
	_, err = testAvroMarshaller.Encode(&receiptNilSender)
	require.NotNil(t, err)

	receiptNilTxHash := receipt
	receiptNilTxHash.TxHash = nil
	_, err = testAvroMarshaller.Encode(&receiptNilTxHash)
	require.NotNil(t, err)
}

func TestEncode_LogAndEvent(t *testing.T) {
	t.Parallel()

	log := schema.Event{}
	_, err := testAvroMarshaller.Encode(&log)
	require.Nil(t, err)

	event := schema.Event{}
	_, err = testAvroMarshaller.Encode(&event)
	require.Nil(t, err)
}

func TestEncode_AccountBalanceUpdate(t *testing.T) {
	t.Parallel()

	acc := schema.AccountBalanceUpdate{
		Address: testscommon.GenerateRandomFixedBytes(62),
	}
	_, err := testAvroMarshaller.Encode(&acc)
	require.Nil(t, err)

	accNilAddress := acc
	accNilAddress.Address = nil
	_, err = testAvroMarshaller.Encode(&accNilAddress)
	require.NotNil(t, err)
}

func TestUInt32SliceToInt32Slice(t *testing.T) {
	t.Parallel()

	out := utility.UInt32SliceToInt32Slice(nil)
	require.Equal(t, []int32{}, out)

	out = utility.UInt32SliceToInt32Slice([]uint32{44, 555})
	require.Equal(t, []int32{44, 555}, out)
}

func TestGetBigIntBytesFromStr(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ret, err := utility.GetBigIntBytesFromStr("4321")
		require.Equal(t, big.NewInt(4321).Bytes(), ret)
		require.Nil(t, err)
	})

	t.Run("empty value, should return bigInt(0)", func(t *testing.T) {
		t.Parallel()

		ret, err := utility.GetBigIntBytesFromStr("")
		require.Equal(t, big.NewInt(0).Bytes(), ret)
		require.Nil(t, err)
	})

	t.Run("invalid value in base 10", func(t *testing.T) {
		t.Parallel()

		ret, err := utility.GetBigIntBytesFromStr("ff")
		require.Nil(t, ret)
		require.NotNil(t, err)
		require.ErrorIs(t, err, utility.ErrInvalidValueInBase10)
		require.True(t, strings.Contains(err.Error(), "ff"))
	})
}

func TestStringSliceToByteSlice(t *testing.T) {
	t.Parallel()

	in := []string{"a", "bb", "ccc"}
	out := utility.StringSliceToByteSlice(in)
	require.Equal(t, [][]byte{[]byte("a"), []byte("bb"), []byte("ccc")}, out)
}

func TestBigIntBytesSliceFromStringSlice(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		in := []string{"12", "345"}
		out, err := utility.GetBigIntBytesSliceFromStringSlice(in)
		require.Nil(t, err)
		require.Equal(t, [][]byte{big.NewInt(12).Bytes(), big.NewInt(345).Bytes()}, out)
	})

	t.Run("invalid string number", func(t *testing.T) {
		t.Parallel()

		in := []string{"12", "34f"}
		out, err := utility.GetBigIntBytesSliceFromStringSlice(in)
		require.Nil(t, out)
		require.NotNil(t, err)
	})
}

func TestGetAddressOrMetachainAddr(t *testing.T) {
	t.Parallel()

	t.Run("normal address", func(t *testing.T) {
		t.Parallel()

		address := utility.GetAddressOrMetachainAddr("address")
		require.Equal(t, []byte("address"), address)
	})

	t.Run("metachain address", func(t *testing.T) {
		t.Parallel()

		address := utility.GetAddressOrMetachainAddr(utility.MetachainShardName)
		require.Equal(t, utility.MetaChainShardAddress(), address)
	})
}
