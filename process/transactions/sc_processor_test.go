package transactions_test

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/transactions"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/smartContractResult"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	"github.com/stretchr/testify/require"
	"math/rand"
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

func generateRandomSCR() *smartContractResult.SmartContractResult {
	return &smartContractResult.SmartContractResult{
		Nonce:          rand.Uint64(),
		Value:          testscommon.GenerateRandomBigInt(),
		RcvAddr:        testscommon.GenerateRandomBytes(),
		SndAddr:        testscommon.GenerateRandomBytes(),
		RelayerAddr:    testscommon.GenerateRandomBytes(),
		RelayedValue:   testscommon.GenerateRandomBigInt(),
		Code:           testscommon.GenerateRandomBytes(),
		Data:           testscommon.GenerateRandomBytes(),
		PrevTxHash:     testscommon.GenerateRandomBytes(),
		OriginalTxHash: testscommon.GenerateRandomBytes(),
		GasLimit:       rand.Uint64(),
		GasPrice:       rand.Uint64(),
		CallType:       vm.CallType(rand.Int()),
		CodeMetadata:   testscommon.GenerateRandomBytes(),
		ReturnMessage:  testscommon.GenerateRandomBytes(),
		OriginalSender: testscommon.GenerateRandomBytes(),
	}
}

func TestScProcessor_ProcessSCs_TwoSCRs_OneNormalTx_ExpectTwoProcessedSCRs(t *testing.T) {
	scp, _ := transactions.NewSCProcessor(&mock.PubKeyConverterStub{})

	tx1 := generateRandomSCR()
	tx2 := generateRandomSCR()
	tx3 := generateRandomTx()

	txPool := map[string]data.TransactionHandler{
		"hash1": tx1,
		"hash2": tx2,
		"hash3": tx3,
	}

	ret := scp.ProcessSCs(txPool, 123)

	require.Len(t, ret, 2)
	requireProcessedSCREqual(t, ret[0], tx1, "hash1", 123, &mock.PubKeyConverterStub{})
	requireProcessedSCREqual(t, ret[1], tx2, "hash2", 123, &mock.PubKeyConverterStub{})
}

func requireProcessedSCREqual(
	t *testing.T,
	processedSCR *schema.SCResult,
	scr *smartContractResult.SmartContractResult,
	hash string,
	timeStamp uint64,
	pubKeyConverter core.PubkeyConverter) {

	require.Equal(t, processedSCR.Hash, []byte(hash))
	require.Equal(t, processedSCR.Nonce, int64(scr.GetNonce()))
	require.Equal(t, processedSCR.GasLimit, int64(scr.GetGasLimit()))
	require.Equal(t, processedSCR.GasPrice, int64(scr.GetGasPrice()))
	require.Equal(t, processedSCR.Value, scr.GetValue().Bytes())
	require.Equal(t, processedSCR.Sender, utility.EncodePubKey(pubKeyConverter, scr.GetSndAddr()))
	require.Equal(t, processedSCR.Receiver, utility.EncodePubKey(pubKeyConverter, scr.GetRcvAddr()))
	require.Equal(t, processedSCR.RelayerAddr, utility.EncodePubKey(pubKeyConverter, scr.GetRelayerAddr()))
	require.Equal(t, processedSCR.RelayedValue, scr.GetRelayedValue().Bytes())
	require.Equal(t, processedSCR.Code, scr.GetCode())
	require.Equal(t, processedSCR.Data, scr.GetData())
	require.Equal(t, processedSCR.PrevTxHash, scr.GetPrevTxHash())
	require.Equal(t, processedSCR.OriginalTxHash, scr.GetOriginalTxHash())
	require.Equal(t, processedSCR.CallType, int32(scr.GetCallType()))
	require.Equal(t, processedSCR.CodeMetadata, scr.GetCodeMetadata())
	require.Equal(t, processedSCR.ReturnMessage, scr.GetReturnMessage())
	require.Equal(t, processedSCR.Timestamp, int64(timeStamp))
}
