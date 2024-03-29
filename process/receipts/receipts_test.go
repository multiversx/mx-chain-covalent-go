package receipts_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-covalent-go/process/receipts"
	"github.com/multiversx/mx-chain-covalent-go/schema"
	"github.com/multiversx/mx-chain-covalent-go/testscommon"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/stretchr/testify/require"
)

func TestReceiptsProcessor_ProcessReceipt(t *testing.T) {
	t.Parallel()

	rp := receipts.NewReceiptsProcessor()
	receipt := &transaction.ApiReceipt{
		Value:   testscommon.GenerateRandomBigInt(),
		SndAddr: "erd1qqq",
		Data:    "ESDTTransfer@555344432d633736663166@1061ed82",
		TxHash:  "975ca52570",
	}

	t.Run("nil receipt, should return empty receipt", func(t *testing.T) {
		t.Parallel()

		processedReceipt, err := rp.ProcessReceipt(nil)
		require.Nil(t, err)
		require.Equal(t, schema.NewReceipt(), processedReceipt)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		processedReceipt, err := rp.ProcessReceipt(receipt)
		require.Nil(t, err)

		hashBytes, err := hex.DecodeString(receipt.TxHash)
		require.Nil(t, err)

		require.Equal(t, &schema.Receipt{
			TxHash: hashBytes,
			Value:  receipt.Value.Bytes(),
			Sender: []byte(receipt.SndAddr),
			Data:   []byte(receipt.Data),
		}, processedReceipt)
	})

	t.Run("invalid hash, should return error", func(t *testing.T) {
		t.Parallel()

		receiptCopy := *receipt
		receiptCopy.TxHash = "rr"
		processedReceipt, err := rp.ProcessReceipt(&receiptCopy)

		require.Nil(t, processedReceipt)
		require.NotNil(t, err)
		require.True(t, strings.Contains(err.Error(), receiptCopy.TxHash))
	})
}
