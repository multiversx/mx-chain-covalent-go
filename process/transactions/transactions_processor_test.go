package transactions

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/stretchr/testify/require"
)

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

func generateApiTx() *transaction.ApiTransactionResult {
	return &transaction.ApiTransactionResult{
		Type:                              "normal",
		ProcessingTypeOnSource:            "source",
		ProcessingTypeOnDestination:       "dest",
		Hash:                              testscommon.GenerateRandHexString(),
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
		PreviousTransactionHash:           testscommon.GenerateRandHexString(),
		OriginalTransactionHash:           testscommon.GenerateRandHexString(),
		ReturnMessage:                     "success",
		OriginalSender:                    "erd1cc",
		Signature:                         testscommon.GenerateRandHexString(),
		SourceShard:                       rand.Uint32(),
		DestinationShard:                  rand.Uint32(),
		BlockNonce:                        rand.Uint64(),
		BlockHash:                         testscommon.GenerateRandHexString(),
		NotarizedAtSourceInMetaNonce:      rand.Uint64(),
		NotarizedAtSourceInMetaHash:       testscommon.GenerateRandHexString(),
		NotarizedAtDestinationInMetaNonce: rand.Uint64(),
		NotarizedAtDestinationInMetaHash:  testscommon.GenerateRandHexString(),
		MiniBlockType:                     "SmartContractResultBlock",
		MiniBlockHash:                     testscommon.GenerateRandHexString(),
		HyperblockNonce:                   rand.Uint64(),
		HyperblockHash:                    testscommon.GenerateRandHexString(),
		Timestamp:                         rand.Int63(),
		Receipt:                           &transaction.ApiReceipt{TxHash: testscommon.GenerateRandHexString()},
		Logs:                              &transaction.ApiLogs{Address: "erd1dd"},
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
}

func generateApiTxs(n int) []*transaction.ApiTransactionResult {
	out := make([]*transaction.ApiTransactionResult, n)

	for i := 0; i < n; i++ {
		out[i] = generateApiTx()
	}

	return out
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

	processLogCalledCt := 0
	logHandler := &mock.LogHandlerStub{
		ProcessLogCalled: func(log *transaction.ApiLogs) *schemaV2.Log {
			processLogCalledCt++
			return &schemaV2.Log{Address: []byte(log.Address)}
		},
	}
	processReceiptCalledCt := 0
	receiptHandler := &mock.ReceiptHandlerStub{
		ProcessReceiptCalled: func(apiReceipt *transaction.ApiReceipt) (*schemaV2.Receipt, error) {
			processReceiptCalledCt++
			hash, err := hex.DecodeString(apiReceipt.TxHash)
			if err != nil {
				return nil, err
			}
			return &schemaV2.Receipt{TxHash: hash}, nil
		},
	}

	txp, _ := NewTransactionProcessor(logHandler, receiptHandler)

	t.Run("should work", func(t *testing.T) {
		apiTxs := generateApiTxs(10)
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Equal(t, processLogCalledCt, 10)
		require.Equal(t, processReceiptCalledCt, 10)
		require.Nil(t, err)
		requireTransactionsProcessedSuccessfully(t, apiTxs, ret, logHandler, receiptHandler)
	})

	t.Run("nil api tx, should skip it", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[0] = nil
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, err)
		requireTransactionsProcessedSuccessfully(t, apiTxs[1:], ret, logHandler, receiptHandler)
	})

	t.Run("invalid hash, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].Hash = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid value, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].Value = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid previousTransactionHash, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].PreviousTransactionHash = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid originalTransactionHash, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].OriginalTransactionHash = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid signature, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].Signature = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid blockHash, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].BlockHash = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid notarizedAtSourceInMetaHash, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].NotarizedAtSourceInMetaHash = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid notarizedAtDestinationInMetaHash, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].NotarizedAtDestinationInMetaHash = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid miniBlockHash, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].MiniBlockHash = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid hyperBlockHash, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].HyperblockHash = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid receipt, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].Receipt.TxHash = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid esdtValues, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].ESDTValues[0] = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})

	t.Run("invalid initiallyPaidFee, should err", func(t *testing.T) {
		apiTxs := generateApiTxs(3)
		apiTxs[1].InitiallyPaidFee = "invalid"
		ret, err := txp.ProcessTransactions(apiTxs)
		require.Nil(t, ret)
		require.NotNil(t, err)
	})
}

func requireTransactionsProcessedSuccessfully(
	t *testing.T,
	apiTxs []*transaction.ApiTransactionResult,
	processedTxs []*schemaV2.Transaction,
	logHandler process.LogHandler,
	receiptHandler process.ReceiptHandler,
) {
	require.Equal(t, len(apiTxs), len(processedTxs))

	for idx := range apiTxs {
		requireTransactionProcessedSuccessfully(t, apiTxs[idx], processedTxs[idx], logHandler, receiptHandler)
	}
}

func requireTransactionProcessedSuccessfully(
	t *testing.T,
	apiTx *transaction.ApiTransactionResult,
	processedTx *schemaV2.Transaction,
	logHandler process.LogHandler,
	receiptHandler process.ReceiptHandler,
) {
	txHash, err := hex.DecodeString(apiTx.Hash)
	require.Nil(t, err)
	value, err := utility.GetBigIntBytesFromStr(apiTx.Value)
	require.Nil(t, err)
	prevTxHash, err := hex.DecodeString(apiTx.PreviousTransactionHash)
	require.Nil(t, err)
	originalTxHash, err := hex.DecodeString(apiTx.OriginalTransactionHash)
	require.Nil(t, err)
	signature, err := hex.DecodeString(apiTx.Signature)
	require.Nil(t, err)
	blockHash, err := hex.DecodeString(apiTx.BlockHash)
	require.Nil(t, err)
	notarizedAtSourceInMetaHash, err := hex.DecodeString(apiTx.NotarizedAtSourceInMetaHash)
	require.Nil(t, err)
	notarizedAtDestinationInMetaHash, err := hex.DecodeString(apiTx.NotarizedAtDestinationInMetaHash)
	require.Nil(t, err)
	miniBlockHash, err := hex.DecodeString(apiTx.MiniBlockHash)
	require.Nil(t, err)
	hyperBlockHash, err := hex.DecodeString(apiTx.HyperblockHash)
	require.Nil(t, err)
	receipt, err := receiptHandler.ProcessReceipt(apiTx.Receipt)
	require.Nil(t, err)
	esdtValues, err := utility.GetBigIntBytesSliceFromStringSlice(apiTx.ESDTValues)
	require.Nil(t, err)
	initiallyPaidFee, err := utility.GetBigIntBytesFromStr(apiTx.InitiallyPaidFee)
	require.Nil(t, err)

	expectedTx := &schemaV2.Transaction{
		Type:                              apiTx.Type,
		ProcessingTypeOnSource:            apiTx.ProcessingTypeOnSource,
		ProcessingTypeOnDestination:       apiTx.ProcessingTypeOnDestination,
		Hash:                              txHash,
		Nonce:                             int64(apiTx.Nonce),
		Round:                             int64(apiTx.Round),
		Epoch:                             int32(apiTx.Epoch),
		Value:                             value,
		Receiver:                          []byte(apiTx.Receiver),
		Sender:                            []byte(apiTx.Sender),
		SenderUserName:                    apiTx.SenderUsername,
		ReceiverUserName:                  apiTx.ReceiverUsername,
		GasPrice:                          int64(apiTx.GasPrice),
		GasLimit:                          int64(apiTx.GasLimit),
		Data:                              apiTx.Data,
		CodeMetadata:                      apiTx.CodeMetadata,
		Code:                              []byte(apiTx.Code),
		PreviousTransactionHash:           prevTxHash,
		OriginalTransactionHash:           originalTxHash,
		ReturnMessage:                     apiTx.ReturnMessage,
		OriginalSender:                    []byte(apiTx.OriginalSender),
		Signature:                         signature,
		SourceShard:                       int32(apiTx.SourceShard),
		DestinationShard:                  int32(apiTx.DestinationShard),
		BlockNonce:                        int64(apiTx.BlockNonce),
		BlockHash:                         blockHash,
		NotarizedAtSourceInMetaNonce:      int64(apiTx.NotarizedAtSourceInMetaNonce),
		NotarizedAtSourceInMetaHash:       notarizedAtSourceInMetaHash,
		NotarizedAtDestinationInMetaNonce: int64(apiTx.NotarizedAtDestinationInMetaNonce),
		NotarizedAtDestinationInMetaHash:  notarizedAtDestinationInMetaHash,
		MiniBlockType:                     apiTx.MiniBlockType,
		MiniBlockHash:                     miniBlockHash,
		HyperBlockNonce:                   int64(apiTx.HyperblockNonce),
		HyperBlockHash:                    hyperBlockHash,
		Timestamp:                         apiTx.Timestamp,
		Receipt:                           receipt,
		Log:                               logHandler.ProcessLog(apiTx.Logs),
		Status:                            apiTx.Status.String(),
		Tokens:                            apiTx.Tokens,
		ESDTValues:                        esdtValues,
		Receivers:                         utility.StringSliceToByteSlice(apiTx.Receivers),
		ReceiversShardIDs:                 utility.UInt32SliceToInt32Slice(apiTx.ReceiversShardIDs),
		Operation:                         apiTx.Operation,
		Function:                          apiTx.Function,
		InitiallyPaidFee:                  initiallyPaidFee,
		IsRelayed:                         apiTx.IsRelayed,
		IsRefund:                          apiTx.IsRefund,
	}

	require.Equal(t, expectedTx, processedTx)
}
