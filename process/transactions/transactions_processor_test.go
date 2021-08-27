package transactions_test

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/transactions"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/elrond-go-core/data/smartContractResult"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/stretchr/testify/require"
)

type headerData struct {
	header     data.HeaderHandler
	headerHash []byte
}

type transactionData struct {
	tx         *transaction.Transaction
	txHash     []byte
	headerData *headerData
}

func generateRandomTx() *transaction.Transaction {
	return &transaction.Transaction{
		Nonce:       rand.Uint64(),
		Value:       testscommon.GenerateRandomBigInt(),
		RcvAddr:     testscommon.GenerateRandomBytes(),
		SndAddr:     testscommon.GenerateRandomBytes(),
		GasLimit:    rand.Uint64(),
		GasPrice:    rand.Uint64(),
		Signature:   testscommon.GenerateRandomBytes(),
		SndUserName: testscommon.GenerateRandomBytes(),
		RcvUserName: testscommon.GenerateRandomBytes(),
	}
}

func generateRandomHeaderData() *headerData {
	return &headerData{
		header:     &block.Header{Round: rand.Uint64(), TimeStamp: rand.Uint64()},
		headerHash: testscommon.GenerateRandomBytes(),
	}
}

func generateRandomTxData(headerData *headerData) *transactionData {
	return &transactionData{
		txHash:     testscommon.GenerateRandomBytes(),
		tx:         generateRandomTx(),
		headerData: headerData,
	}
}

func TestNewTransactionProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (core.PubkeyConverter, hashing.Hasher, marshal.Marshalizer)
		expectedErr error
	}{
		{
			args: func() (core.PubkeyConverter, hashing.Hasher, marshal.Marshalizer) {
				return nil, &mock.HasherMock{}, &mock.MarshallerStub{}
			},
			expectedErr: covalent.ErrNilPubKeyConverter,
		},
		{
			args: func() (core.PubkeyConverter, hashing.Hasher, marshal.Marshalizer) {
				return &mock.PubKeyConverterStub{}, nil, &mock.MarshallerStub{}
			},
			expectedErr: covalent.ErrNilHasher,
		},
		{
			args: func() (core.PubkeyConverter, hashing.Hasher, marshal.Marshalizer) {
				return &mock.PubKeyConverterStub{}, &mock.HasherMock{}, nil
			},
			expectedErr: covalent.ErrNilMarshaller,
		},
		{
			args: func() (core.PubkeyConverter, hashing.Hasher, marshal.Marshalizer) {
				return &mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := transactions.NewTransactionProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestTransactionProcessor_ProcessTransactions_InvalidBody_ExpectError(t *testing.T) {
	t.Parallel()
	//TODO: ADD PARALEL EVERYWHEREEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE

	hData := generateRandomHeaderData()
	txPool := map[string]data.TransactionHandler{}
	pool := &indexer.Pool{
		Txs: txPool,
	}
	body := data.BodyHandler(nil)

	txp, _ := transactions.NewTransactionProcessor(&mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
	_, err := txp.ProcessTransactions(hData.header, hData.headerHash, body, pool)

	require.Equal(t, covalent.ErrBlockBodyAssertion, err)
}

func TestTransactionProcessor_ProcessTransactions_InvalidMarshaller_ExpectZeroProcessedTxs(t *testing.T) {
	hData := generateRandomHeaderData()
	txData1 := generateRandomTxData(hData)

	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{txData1.txHash},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            block.TxBlock},
	},
	}

	txPool := map[string]data.TransactionHandler{
		string(txData1.txHash): txData1.tx,
	}
	pool := &indexer.Pool{
		Txs: txPool,
	}

	errMarshaller := errors.New("err marshaller")
	txp, _ := transactions.NewTransactionProcessor(
		&mock.PubKeyConverterStub{},
		&mock.HasherMock{},
		&mock.MarshallerStub{
			MarshalCalled: func(obj interface{}) ([]byte, error) {
				return nil, errMarshaller
			},
		})
	ret, err := txp.ProcessTransactions(hData.header, hData.headerHash, body, pool)

	require.Nil(t, err)
	require.Len(t, ret, 0)
}

func TestTransactionProcessor_ProcessTransactions_OneEmptyTxBlock_ExpectZeroProcessedTxs(t *testing.T) {
	hData := generateRandomHeaderData()

	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            block.TxBlock},
	},
	}

	txp, _ := transactions.NewTransactionProcessor(&mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
	ret, err := txp.ProcessTransactions(hData.header, hData.headerHash, body, &indexer.Pool{})

	require.Nil(t, err)
	require.Len(t, ret, 0)
}

func TestTransactionProcessor_ProcessTransactions_OneTxBlock_TxNotFoundInPool_ExpectZeroProcessedTxs(t *testing.T) {
	hData := generateRandomHeaderData()

	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{[]byte("tx not found")},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            block.TxBlock},
	},
	}

	txp, _ := transactions.NewTransactionProcessor(&mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
	ret, err := txp.ProcessTransactions(hData.header, hData.headerHash, body, &indexer.Pool{})

	require.Nil(t, err)
	require.Len(t, ret, 0)
}

func TestTransactionProcessor_ProcessTransactions_OneTxBlock_OneTx_ExpectOneProcessedTx(t *testing.T) {
	hData := generateRandomHeaderData()
	txData1 := generateRandomTxData(hData)

	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{txData1.txHash},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            block.TxBlock},
	},
	}

	txPool := map[string]data.TransactionHandler{
		string(txData1.txHash): txData1.tx,
	}
	pool := &indexer.Pool{
		Txs: txPool,
	}

	txp, _ := transactions.NewTransactionProcessor(&mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
	ret, _ := txp.ProcessTransactions(hData.header, hData.headerHash, body, pool)

	require.Len(t, ret, 1)
	requireProcessedTransactionEqual(t, ret[0], txData1, body.GetMiniBlocks()[0], &mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
}

func TestTransactionProcessor_ProcessTransactions_OneTxBLock_TwoNormalTxs_ExpectTwoProcessedTxs(t *testing.T) {
	hData := generateRandomHeaderData()

	txData1 := generateRandomTxData(hData)
	txData2 := generateRandomTxData(hData)

	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{txData1.txHash, txData2.txHash},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            block.TxBlock},
	},
	}

	txPool := map[string]data.TransactionHandler{
		string(txData1.txHash): txData1.tx,
		string(txData2.txHash): txData2.tx,
	}
	pool := &indexer.Pool{
		Txs: txPool,
	}

	txp, _ := transactions.NewTransactionProcessor(&mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
	ret, _ := txp.ProcessTransactions(hData.header, hData.headerHash, body, pool)

	require.Len(t, ret, 2)

	requireProcessedTransactionEqual(t, ret[0], txData1, body.GetMiniBlocks()[0], &mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
	requireProcessedTransactionEqual(t, ret[1], txData2, body.GetMiniBlocks()[0], &mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
}

func TestTransactionProcessor_ProcessTransactions_TwoTxBlocks_TwoTxs_ExpectTwoProcessedTx(t *testing.T) {
	hData := generateRandomHeaderData()

	txData1 := generateRandomTxData(hData)
	txData2 := generateRandomTxData(hData)

	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{txData1.txHash},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            block.TxBlock},
		{
			TxHashes:        [][]byte{txData2.txHash},
			ReceiverShardID: 3,
			SenderShardID:   4,
			Type:            block.TxBlock},
	},
	}

	txPool := map[string]data.TransactionHandler{
		string(txData1.txHash): txData1.tx,
		string(txData2.txHash): txData2.tx,
	}
	pool := &indexer.Pool{
		Txs: txPool,
	}

	txp, _ := transactions.NewTransactionProcessor(&mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
	ret, _ := txp.ProcessTransactions(hData.header, hData.headerHash, body, pool)

	require.Len(t, ret, 2)

	requireProcessedTransactionEqual(t, ret[0], txData1, body.GetMiniBlocks()[0], &mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
	requireProcessedTransactionEqual(t, ret[1], txData2, body.GetMiniBlocks()[1], &mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
}

func TestTransactionProcessor_ProcessTransactions_OneTxBlock_OneSCRTx_ExpectZeroProcessedTxs(t *testing.T) {
	hData := generateRandomHeaderData()
	scrHash := []byte("scr tx hash")

	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{scrHash},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            block.TxBlock},
	},
	}

	txPool := map[string]data.TransactionHandler{
		string(scrHash): &smartContractResult.SmartContractResult{},
	}
	pool := &indexer.Pool{
		Txs: txPool,
	}
	txp, _ := transactions.NewTransactionProcessor(&mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
	ret, err := txp.ProcessTransactions(hData.header, hData.headerHash, body, pool)

	require.Nil(t, err)
	require.Len(t, ret, 0)
}

func TestTransactionProcessor_ProcessTransactions_OneSCRBlock_OneSCRTx_ExpectZeroProcessedTxs(t *testing.T) {
	hData := generateRandomHeaderData()
	scrHash := []byte("scr tx hash")
	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{scrHash},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            block.SmartContractResultBlock},
	},
	}

	txPool := map[string]data.TransactionHandler{
		string(scrHash): &smartContractResult.SmartContractResult{}}
	pool := &indexer.Pool{
		Txs: txPool,
	}

	txp, _ := transactions.NewTransactionProcessor(&mock.PubKeyConverterStub{}, &mock.HasherMock{}, &mock.MarshallerStub{})
	ret, err := txp.ProcessTransactions(hData.header, hData.headerHash, body, pool)

	require.Nil(t, err)
	require.Len(t, ret, 0)
}

func requireProcessedTransactionEqual(
	t *testing.T,
	processedTx *schema.Transaction,
	td *transactionData,
	miniBlock *block.MiniBlock,
	pubKeyConverter core.PubkeyConverter,
	hasher hashing.Hasher,
	marshaller marshal.Marshalizer) {

	tx := td.tx
	hData := td.headerData
	mbHash, _ := core.CalculateHash(marshaller, hasher, miniBlock)

	require.Equal(t, td.txHash, processedTx.Hash)
	require.Equal(t, int64(tx.GetNonce()), processedTx.Nonce)
	require.Equal(t, tx.GetValue().Bytes(), processedTx.Value)
	require.Equal(t, utility.EncodePubKey(pubKeyConverter, tx.GetRcvAddr()), processedTx.Receiver)
	require.Equal(t, utility.EncodePubKey(pubKeyConverter, tx.GetSndAddr()), processedTx.Sender)
	require.Equal(t, int32(miniBlock.GetReceiverShardID()), processedTx.ReceiverShard)
	require.Equal(t, int32(miniBlock.GetSenderShardID()), processedTx.SenderShard)
	require.Equal(t, int64(tx.GetGasPrice()), processedTx.GasPrice)
	require.Equal(t, int64(tx.GetGasLimit()), processedTx.GasLimit)
	require.Equal(t, tx.GetSignature(), processedTx.Signature)
	require.Equal(t, tx.GetSndUserName(), processedTx.SenderUserName)
	require.Equal(t, tx.GetRcvUserName(), processedTx.ReceiverUserName)
	require.Equal(t, mbHash, processedTx.MiniBlockHash)
	require.Equal(t, hData.headerHash, processedTx.BlockHash)
	require.Equal(t, int64(hData.header.GetRound()), processedTx.Round)
	require.Equal(t, int64(hData.header.GetTimeStamp()), processedTx.Timestamp)
}
