package utility_test

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/stretchr/testify/require"
)

var testAvroMarshaller = &utility.AvroMarshaller{}

func TestHexSliceToByteSlice_DifferentValues(t *testing.T) {
	in := []string{"0a", "0b", "0c"}
	out, err := utility.HexSliceToByteSlice(in)
	require.Nil(t, err)

	require.Len(t, out, 3)
	require.Equal(t, []byte{0xa}, out[0])
	require.Equal(t, []byte{0xb}, out[1])
	require.Equal(t, []byte{0xc}, out[2])
}

func TestHexSliceToByteSlice_EmptyInput(t *testing.T) {
	out, err := utility.HexSliceToByteSlice([]string{})

	require.Nil(t, err)
	require.Len(t, out, 0)
}

func TestHexSliceToByteSlice_InvalidString_ExpectError(t *testing.T) {
	in := []string{"0a", "xz", "0c"}
	out, err := utility.HexSliceToByteSlice(in)

	require.NotNil(t, err)
	require.Nil(t, out)
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
	require.Equal(t, utility.GetBytes(nil), []byte{})

	x := big.NewInt(10)
	require.Equal(t, []byte{0xa}, utility.GetBytes(x))
}

func TestEncodeDecode(t *testing.T) {
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

func TestEncode_Block(t *testing.T) {
	block := schema.Block{
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
	require.NotNil(t, err)
}

func TestEncode_MiniBlock(t *testing.T) {
	mb := schema.MiniBlock{
		Hash: testscommon.GenerateRandomFixedBytes(32),
	}
	_, err := testAvroMarshaller.Encode(&mb)
	require.Nil(t, err)

	mbNilHash := mb
	mbNilHash.Hash = nil
	_, err = testAvroMarshaller.Encode(&mbNilHash)
	require.NotNil(t, err)
}

func TestEncode_EpochStartInfo(t *testing.T) {
	info := schema.EpochStartInfo{}
	_, err := testAvroMarshaller.Encode(&info)
	require.Nil(t, err)
}

func TestEncode_Transaction(t *testing.T) {
	tx := schema.Transaction{
		Hash:          testscommon.GenerateRandomFixedBytes(32),
		MiniBlockHash: testscommon.GenerateRandomFixedBytes(32),
		BlockHash:     testscommon.GenerateRandomFixedBytes(32),
		Receiver:      testscommon.GenerateRandomFixedBytes(62),
		Sender:        testscommon.GenerateRandomFixedBytes(62),
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
	require.NotNil(t, err)

	txNilReceiver := tx
	txNilReceiver.Receiver = nil
	_, err = testAvroMarshaller.Encode(&txNilReceiver)
	require.NotNil(t, err)

	txNilSender := tx
	txNilSender.Sender = nil
	_, err = testAvroMarshaller.Encode(&txNilSender)
	require.NotNil(t, err)
}

func TestEncode_SCR(t *testing.T) {
	scRes := schema.SCResult{
		Hash:           testscommon.GenerateRandomFixedBytes(32),
		Sender:         testscommon.GenerateRandomFixedBytes(62),
		Receiver:       testscommon.GenerateRandomFixedBytes(62),
		PrevTxHash:     testscommon.GenerateRandomFixedBytes(32),
		OriginalTxHash: testscommon.GenerateRandomFixedBytes(32),
	}
	_, err := testAvroMarshaller.Encode(&scRes)
	require.Nil(t, err)

	scResNilHash := scRes
	scResNilHash.Hash = nil
	_, err = testAvroMarshaller.Encode(&scResNilHash)
	require.NotNil(t, err)

	scResNilSender := scRes
	scResNilSender.Sender = nil
	_, err = testAvroMarshaller.Encode(&scResNilSender)
	require.NotNil(t, err)

	scResNilReceiver := scRes
	scResNilReceiver.Receiver = nil
	_, err = testAvroMarshaller.Encode(&scResNilReceiver)
	require.NotNil(t, err)

	scResNilPrevTxHash := scRes
	scResNilPrevTxHash.PrevTxHash = nil
	_, err = testAvroMarshaller.Encode(&scResNilPrevTxHash)
	require.NotNil(t, err)

	scResNilOriginalTxHash := scRes
	scResNilOriginalTxHash.OriginalTxHash = nil
	_, err = testAvroMarshaller.Encode(&scResNilOriginalTxHash)
	require.NotNil(t, err)
}

func TestEncode_Receipt(t *testing.T) {
	receipt := schema.Receipt{
		Hash:   testscommon.GenerateRandomFixedBytes(32),
		Sender: testscommon.GenerateRandomFixedBytes(62),
		TxHash: testscommon.GenerateRandomFixedBytes(32),
	}

	_, err := testAvroMarshaller.Encode(&receipt)
	require.Nil(t, err)

	receiptNilHash := receipt
	receiptNilHash.Hash = nil
	_, err = testAvroMarshaller.Encode(&receiptNilHash)
	require.NotNil(t, err)

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
	log := schema.Event{}
	_, err := testAvroMarshaller.Encode(&log)
	require.Nil(t, err)

	event := schema.Event{}
	_, err = testAvroMarshaller.Encode(&event)
	require.Nil(t, err)
}

func TestEncode_AccountBalanceUpdate(t *testing.T) {
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

func TestEncode_BlockResult(t *testing.T) {
	block := schema.Block{
		Hash:          testscommon.GenerateRandomFixedBytes(32),
		StateRootHash: testscommon.GenerateRandomFixedBytes(32),
	}

	blockRes := schema.BlockResult{
		Block: &block,
	}
	_, err := testAvroMarshaller.Encode(&blockRes)
	require.Nil(t, err)

	blockResNilBlock := blockRes
	blockResNilBlock.Block = nil
	_, err = testAvroMarshaller.Encode(&blockResNilBlock)
	require.NotNil(t, err)
}
