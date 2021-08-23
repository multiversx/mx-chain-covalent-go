package receipts_test

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/process/receipts"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/receipt"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/stretchr/testify/require"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
)

func TestNewReceiptsProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() core.PubkeyConverter
		expectedErr error
	}{
		{
			args: func() core.PubkeyConverter {
				return nil
			},
			expectedErr: covalent.ErrNilPubKeyConverter,
		},
		{
			args: func() core.PubkeyConverter {
				return &mock.PubKeyConverterStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := receipts.NewReceiptsProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func generateRandomReceipt() *receipt.Receipt {
	return &receipt.Receipt{
		Value:   big.NewInt(rand.Int63()),
		SndAddr: []byte(strconv.Itoa(rand.Int())),
		Data:    []byte(strconv.Itoa(rand.Int())),
		TxHash:  []byte(strconv.Itoa(rand.Int())),
	}
}

func TestReceiptsProcessor_ProcessReceipts_TwoReceipts_OneNormalTx_ExpectTwoProcessedReceipts(t *testing.T) {
	rp, _ := receipts.NewReceiptsProcessor(&mock.PubKeyConverterStub{})

	receipt1 := generateRandomReceipt()
	receipt2 := generateRandomReceipt()

	txPool := map[string]data.TransactionHandler{
		"hash1": receipt1,
		"hash2": receipt2,
		"hash3": &transaction.Transaction{},
	}

	ret := rp.ProcessReceipts(txPool, 123)

	require.Len(t, ret, 2)

	requireProcessedReceiptEqual(t, ret[0], receipt1, "hash1", 123, &mock.PubKeyConverterStub{})
	requireProcessedReceiptEqual(t, ret[1], receipt2, "hash2", 123, &mock.PubKeyConverterStub{})
}

func requireProcessedReceiptEqual(
	t *testing.T,
	processedReceipt *schema.Receipt,
	rec *receipt.Receipt,
	receiptHash string,
	timestamp uint64,
	pubKeyConverter core.PubkeyConverter) {

	require.Equal(t, processedReceipt.Hash, []byte(receiptHash))
	require.Equal(t, processedReceipt.Value, rec.GetValue().Bytes())
	require.Equal(t, processedReceipt.Sender, utility.EncodePubKey(pubKeyConverter, rec.GetSndAddr()))
	require.Equal(t, processedReceipt.Data, rec.GetData())
	require.Equal(t, processedReceipt.TxHash, rec.GetTxHash())
	require.Equal(t, processedReceipt.Timestamp, int64(timestamp))
}
