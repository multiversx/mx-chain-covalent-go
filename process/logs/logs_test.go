package logs_test

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/process/logs"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
)

func TestNewLogsProcessor(t *testing.T) {
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
		_, err := logs.NewLogsProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestLogsProcessor_ProcessLogs_OneNilLog_ExpectZeroProcessedLogs(t *testing.T) {
	lp, _ := logs.NewLogsProcessor(&mock.PubKeyConverterStub{})

	logsAndEvents := map[string]data.LogHandler{
		"hash1": nil,
	}

	ret, _ := lp.ProcessLogs(logsAndEvents)
	require.Len(t, ret, 0)

}

func TestLogsProcessor_ProcessLogs_OneLog_NoEvent_ExpectOneProcessedLogAndZeroEvents(t *testing.T) {
	lp, _ := logs.NewLogsProcessor(&mock.PubKeyConverterStub{})

	event := generateRandomEvent()
	log := &transaction.Log{
		Address: []byte(strconv.Itoa(rand.Int())),
		Events:  []*transaction.Event{event},
	}

	logsAndEvents := map[string]data.LogHandler{
		"hash1": log,
	}

	ret, _ := lp.ProcessLogs(logsAndEvents)
	require.Len(t, ret, 1)
	require.Len(t, ret[0].Events, 1)

}

func generateRandomEvent() *transaction.Event {
	return &transaction.Event{
		Address:    []byte(strconv.Itoa(rand.Int())),
		Identifier: []byte(strconv.Itoa(rand.Int())),
		Topics:     [][]byte{[]byte(strconv.Itoa(rand.Int()))},
		Data:       []byte(strconv.Itoa(rand.Int())),
	}
}
