package transactions_test

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/process/transactions"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/smartContractResult"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestNewTransactionProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (core.PubkeyConverter, hashing.Hasher, marshal.Marshalizer)
		expectedErr error
	}{
		{
			args: func() (core.PubkeyConverter, hashing.Hasher, marshal.Marshalizer) {
				return nil, &mock.HasherStub{}, &mock.MarshallerStub{}
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
				return &mock.PubKeyConverterStub{}, &mock.HasherStub{}, nil
			},
			expectedErr: covalent.ErrNilMarshaller,
		},
		{
			args: func() (core.PubkeyConverter, hashing.Hasher, marshal.Marshalizer) {
				return &mock.PubKeyConverterStub{}, &mock.HasherStub{}, &mock.MarshallerStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := transactions.NewTransactionProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestTransactionProcessor_ProcessTransactions(t *testing.T) {
	txHash1 := []byte("x")
	txHash2 := []byte("y")
	txHash3 := []byte("z")
	txHash4 := []byte("scr block")
	txHash5 := []byte("tx block, scr tx")

	headerHash := []byte("header hash")
	header := &block.Header{Round: 111, TimeStamp: 222}
	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{txHash1, txHash2},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            block.TxBlock},
		{
			TxHashes:        [][]byte{txHash3, []byte("tx not found in pool")},
			ReceiverShardID: 4,
			SenderShardID:   5,
			Type:            block.TxBlock},
		{
			TxHashes:        [][]byte{},
			ReceiverShardID: 6,
			SenderShardID:   7,
			Type:            block.TxBlock},
		{
			TxHashes:        [][]byte{txHash4},
			ReceiverShardID: 8,
			SenderShardID:   9,
			Type:            block.SmartContractResultBlock},
		{
			TxHashes:        [][]byte{txHash5},
			ReceiverShardID: 6,
			SenderShardID:   7,
			Type:            block.TxBlock},
	}}

	tx1 := &transaction.Transaction{
		Nonce:       1,
		Value:       big.NewInt(2),
		RcvAddr:     []byte("rcv1"),
		SndAddr:     []byte("snd1"),
		GasLimit:    3,
		GasPrice:    4,
		Signature:   []byte("sig1"),
		SndUserName: []byte("sndName1"),
		RcvUserName: []byte("rcvName1"),
	}
	tx2 := &transaction.Transaction{
		Nonce:       5,
		Value:       big.NewInt(6),
		RcvAddr:     []byte("rcv2"),
		SndAddr:     []byte("snd2"),
		GasLimit:    7,
		GasPrice:    8,
		Signature:   []byte("sig2"),
		SndUserName: []byte("sndName2"),
		RcvUserName: nil,
	}
	tx3 := &transaction.Transaction{
		Nonce:       9,
		Value:       big.NewInt(10),
		RcvAddr:     []byte("rcv3"),
		SndAddr:     []byte("snd3"),
		GasLimit:    11,
		GasPrice:    12,
		Signature:   []byte("sig3"),
		SndUserName: nil,
		RcvUserName: nil,
	}
	tx4 := &smartContractResult.SmartContractResult{
		Nonce:    13,
		Value:    big.NewInt(14),
		RcvAddr:  []byte("rcv4"),
		SndAddr:  []byte("snd4"),
		GasLimit: 15,
		GasPrice: 16,
	}
	tx5 := &smartContractResult.SmartContractResult{
		Nonce:    17,
		Value:    big.NewInt(18),
		RcvAddr:  []byte("rcv5"),
		SndAddr:  []byte("snd5"),
		GasLimit: 19,
		GasPrice: 20,
	}

	txPool := map[string]data.TransactionHandler{
		string(txHash1): tx1,
		string(txHash2): tx2,
		string(txHash3): tx3,
		string(txHash4): tx4,
		string(txHash5): tx5,
	}

	txp, _ := transactions.NewTransactionProcessor(&mock.PubKeyConverterStub{}, &mock.HasherStub{}, &mock.MarshallerStub{})
	ret, _ := txp.ProcessTransactions(header, headerHash, body, txPool)

	require.Len(t, ret, 3)

	require.Equal(t, ret[0].Hash, txHash1)
	require.Equal(t, ret[0].MiniBlockHash, []byte("ok"))
	require.Equal(t, ret[0].BlockHash, headerHash)
	require.Equal(t, ret[0].Nonce, int64(1))
	require.Equal(t, ret[0].Round, int64(111))
	require.Equal(t, ret[0].Value, big.NewInt(2).Bytes())
	require.Equal(t, ret[0].Receiver, []byte("erd1rcv1"))
	require.Equal(t, ret[0].Sender, []byte("erd1snd1"))
	require.Equal(t, ret[0].ReceiverShard, int32(1))
	require.Equal(t, ret[0].SenderShard, int32(2))
	require.Equal(t, ret[0].GasPrice, int64(4))
	require.Equal(t, ret[0].GasLimit, int64(3))
	require.Equal(t, ret[0].Signature, []byte("sig1"))
	require.Equal(t, ret[0].Timestamp, int64(222))
	require.Equal(t, ret[0].SenderUserName, []byte("sndName1"))
	require.Equal(t, ret[0].ReceiverUserName, []byte("rcvName1"))

	require.Equal(t, ret[1].Hash, txHash2)
	require.Equal(t, ret[1].MiniBlockHash, []byte("ok"))
	require.Equal(t, ret[1].BlockHash, headerHash)
	require.Equal(t, ret[1].Nonce, int64(5))
	require.Equal(t, ret[1].Round, int64(111))
	require.Equal(t, ret[1].Value, big.NewInt(6).Bytes())
	require.Equal(t, ret[1].Receiver, []byte("erd1rcv2"))
	require.Equal(t, ret[1].Sender, []byte("erd1snd2"))
	require.Equal(t, ret[1].ReceiverShard, int32(1))
	require.Equal(t, ret[1].SenderShard, int32(2))
	require.Equal(t, ret[1].GasPrice, int64(8))
	require.Equal(t, ret[1].GasLimit, int64(7))
	require.Equal(t, ret[1].Signature, []byte("sig2"))
	require.Equal(t, ret[1].Timestamp, int64(222))
	require.Equal(t, ret[1].SenderUserName, []byte("sndName2"))
	require.Equal(t, ret[1].ReceiverUserName, []byte(nil))

	require.Equal(t, ret[2].Hash, txHash3)
	require.Equal(t, ret[2].MiniBlockHash, []byte("ok"))
	require.Equal(t, ret[2].BlockHash, headerHash)
	require.Equal(t, ret[2].Nonce, int64(9))
	require.Equal(t, ret[2].Round, int64(111))
	require.Equal(t, ret[2].Value, big.NewInt(10).Bytes())
	require.Equal(t, ret[2].Receiver, []byte("erd1rcv3"))
	require.Equal(t, ret[2].Sender, []byte("erd1snd3"))
	require.Equal(t, ret[2].ReceiverShard, int32(4))
	require.Equal(t, ret[2].SenderShard, int32(5))
	require.Equal(t, ret[2].GasPrice, int64(12))
	require.Equal(t, ret[2].GasLimit, int64(11))
	require.Equal(t, ret[2].Signature, []byte("sig3"))
	require.Equal(t, ret[2].Timestamp, int64(222))
	require.Equal(t, ret[2].SenderUserName, []byte(nil))
	require.Equal(t, ret[2].ReceiverUserName, []byte(nil))
}
