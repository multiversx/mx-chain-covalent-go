package transactions

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/stretchr/testify/require"
)

func generateRandHexString() string {
	randBytes := testscommon.GenerateRandomFixedBytes(32)
	return hex.EncodeToString(randBytes)
}

func generateStringSlice(n int) []string {
	ret := make([]string, n)

	for i := 0; i < n; i++ {
		randStr := fmt.Sprintf("str%d", i)
		ret[i] = randStr
	}

	return ret
}

func generateBigUIntStringSlice(n int) []string {
	ret := make([]string, n)

	for i := 0; i < n; i++ {
		randBI := testscommon.GenerateRandomBigInt()
		ret[i] = randBI.String()
	}

	return ret
}

func generateUInt32Slice(n int) []uint32 {
	ret := make([]uint32, n)

	for i := 0; i < n; i++ {
		ret[i] = rand.Uint32()
	}

	return ret
}

func TestNewTransactionProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil log processor, should return error", func(t *testing.T) {
		t.Parallel()

		txp, err := NewTransactionProcessor(nil, &mock.ReceiptHandlerStub{})
		require.Nil(t, txp)
		require.Equal(t, errNilLogProcessor, err)
	})

	t.Run("nil receipt processor, should return error", func(t *testing.T) {
		t.Parallel()

		txp, err := NewTransactionProcessor(&mock.LogHandlerStub{}, nil)
		require.Nil(t, txp)
		require.Equal(t, errNilReceiptProcessor, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		txp, err := NewTransactionProcessor(&mock.LogHandlerStub{}, &mock.ReceiptHandlerStub{})
		require.Nil(t, err)
		require.NotNil(t, txp)
	})
}

func TestTransactionProcessor_ProcessTransactions(t *testing.T) {
	t.Parallel()

	tx := &transaction.ApiTransactionResult{
		Type:                              "normal",
		ProcessingTypeOnSource:            "source",
		ProcessingTypeOnDestination:       "dest",
		Hash:                              generateRandHexString(),
		HashBytes:                         nil,
		Nonce:                             rand.Uint64(),
		Round:                             rand.Uint64(),
		Epoch:                             rand.Uint32(),
		Value:                             testscommon.GenerateRandomBigInt().String(),
		Receiver:                          "erd1aa",
		Sender:                            "erd1bb",
		SenderUsername:                    testscommon.GenerateRandomBytes(),
		ReceiverUsername:                  testscommon.GenerateRandomBytes(),
		GasPrice:                          rand.Uint64(),
		GasLimit:                          rand.Uint64(),
		Data:                              testscommon.GenerateRandomBytes(),
		CodeMetadata:                      testscommon.GenerateRandomBytes(),
		Code:                              "code",
		PreviousTransactionHash:           generateRandHexString(),
		OriginalTransactionHash:           generateRandHexString(),
		ReturnMessage:                     "success",
		OriginalSender:                    "erd1cc",
		Signature:                         generateRandHexString(),
		SourceShard:                       rand.Uint32(),
		DestinationShard:                  rand.Uint32(),
		BlockNonce:                        rand.Uint64(),
		BlockHash:                         generateRandHexString(),
		NotarizedAtSourceInMetaNonce:      rand.Uint64(),
		NotarizedAtSourceInMetaHash:       generateRandHexString(),
		NotarizedAtDestinationInMetaNonce: rand.Uint64(),
		NotarizedAtDestinationInMetaHash:  generateRandHexString(),
		MiniBlockType:                     "SmartContractResultBlock",
		MiniBlockHash:                     generateRandHexString(),
		HyperblockNonce:                   rand.Uint64(),
		HyperblockHash:                    generateRandHexString(),
		Timestamp:                         rand.Int63(),
		Receipt:                           nil,
		Logs:                              nil,
		Status:                            "status",
		Tokens:                            generateStringSlice(rand.Int()%10 + 1),
		ESDTValues:                        generateBigUIntStringSlice(rand.Int()%10 + 1),
		Receivers:                         generateStringSlice(rand.Int()%10 + 1),
		ReceiversShardIDs:                 generateUInt32Slice(rand.Int() % 4),
		Operation:                         "transfer",
		Function:                          "function",
		InitiallyPaidFee:                  testscommon.GenerateRandomBigInt().String(),
		IsRelayed:                         true,
		IsRefund:                          true,
	}

	txp, _ := NewTransactionProcessor(&mock.LogHandlerStub{}, &mock.ReceiptHandlerStub{})
	ret, err := txp.ProcessTransactions([]*transaction.ApiTransactionResult{tx})
	require.Nil(t, err)
	require.Len(t, ret, 1)
}
