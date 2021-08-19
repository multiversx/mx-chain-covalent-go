package transactions_test

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/process/transactions"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/smartContractResult"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestNewSCProcessor(t *testing.T) {
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
		_, err := transactions.NewSCProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestScProcessor_ProcessSCs(t *testing.T) {
	scp, _ := transactions.NewSCProcessor(&mock.PubKeyConverterStub{})

	tx1 := &smartContractResult.SmartContractResult{
		Nonce:          1,
		Value:          big.NewInt(2),
		RcvAddr:        []byte("rcv1"),
		SndAddr:        []byte("snd1"),
		RelayerAddr:    []byte("rly1"),
		RelayedValue:   big.NewInt(3),
		Code:           []byte("code1"),
		Data:           []byte("data1"),
		PrevTxHash:     []byte("prevHash1"),
		OriginalTxHash: []byte("origHash1"),
		GasLimit:       4,
		GasPrice:       5,
		CallType:       6,
		CodeMetadata:   []byte("codeMeta1"),
		ReturnMessage:  []byte("retMsg1"),
		OriginalSender: []byte("origSnd1"),
	}
	tx2 := &smartContractResult.SmartContractResult{
		Nonce:          7,
		Value:          big.NewInt(8),
		RcvAddr:        []byte("rcv2"),
		SndAddr:        []byte("snd2"),
		RelayerAddr:    []byte("rly2"),
		RelayedValue:   big.NewInt(9),
		Code:           []byte("code2"),
		Data:           []byte("data2"),
		PrevTxHash:     []byte("prevHash2"),
		OriginalTxHash: []byte("origHash2"),
		GasLimit:       10,
		GasPrice:       11,
		CallType:       12,
		CodeMetadata:   []byte("codeMeta2"),
		ReturnMessage:  []byte("retMsg2"),
		OriginalSender: []byte("origSnd2"),
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

	txPool := map[string]data.TransactionHandler{
		"hash1": tx1,
		"hash2": tx2,
		"hash3": tx3,
	}

	ret, _ := scp.ProcessSCs(txPool, 123)

	require.Len(t, ret, 2)

	require.Equal(t, ret[0].Hash, []byte("hash1"))
	require.Equal(t, ret[0].Nonce, int64(1))
	require.Equal(t, ret[0].GasLimit, int64(4))
	require.Equal(t, ret[0].GasPrice, int64(5))
	require.Equal(t, ret[0].Value, big.NewInt(2).Bytes())
	require.Equal(t, ret[0].Sender, []byte("erd1snd1"))
	require.Equal(t, ret[0].Receiver, []byte("erd1rcv1"))
	require.Equal(t, ret[0].RelayerAddr, []byte("erd1rly1"))
	require.Equal(t, ret[0].RelayedValue, big.NewInt(3).Bytes())
	require.Equal(t, ret[0].Code, []byte("code1"))
	require.Equal(t, ret[0].Data, []byte("data1"))
	require.Equal(t, ret[0].PrevTxHash, []byte("prevHash1"))
	require.Equal(t, ret[0].OriginalTxHash, []byte("origHash1"))
	require.Equal(t, ret[0].CallType, int32(6))
	require.Equal(t, ret[0].CodeMetadata, []byte("codeMeta1"))
	require.Equal(t, ret[0].ReturnMessage, []byte("retMsg1"))
	require.Equal(t, ret[0].Timestamp, int64(123))

	require.Equal(t, ret[1].Hash, []byte("hash2"))
	require.Equal(t, ret[1].Nonce, int64(7))
	require.Equal(t, ret[1].GasLimit, int64(10))
	require.Equal(t, ret[1].GasPrice, int64(11))
	require.Equal(t, ret[1].Value, big.NewInt(8).Bytes())
	require.Equal(t, ret[1].Sender, []byte("erd1snd2"))
	require.Equal(t, ret[1].Receiver, []byte("erd1rcv2"))
	require.Equal(t, ret[1].RelayerAddr, []byte("erd1rly2"))
	require.Equal(t, ret[1].RelayedValue, big.NewInt(9).Bytes())
	require.Equal(t, ret[1].Code, []byte("code2"))
	require.Equal(t, ret[1].Data, []byte("data2"))
	require.Equal(t, ret[1].PrevTxHash, []byte("prevHash2"))
	require.Equal(t, ret[1].OriginalTxHash, []byte("origHash2"))
	require.Equal(t, ret[1].CallType, int32(12))
	require.Equal(t, ret[1].CodeMetadata, []byte("codeMeta2"))
	require.Equal(t, ret[1].ReturnMessage, []byte("retMsg2"))
	require.Equal(t, ret[1].Timestamp, int64(123))
}
