package logs_test

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/process/logs"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go-logger/check"
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

	log := &transaction.Log{
		Address: []byte(strconv.Itoa(rand.Int())),
		Events:  []*transaction.Event{},
	}
	logsAndEvents := map[string]data.LogHandler{
		"hash1": log,
	}

	ret, _ := lp.ProcessLogs(logsAndEvents)

	require.Len(t, ret, 1)
	require.Len(t, ret[0].Events, 0)
}

func TestLogsProcessor_ProcessLogs_OneLog_OneEvent_ExpectOneProcessedLogAndOneEvent(t *testing.T) {
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

	requireProcessedLogEqual(t, ret[0], log, "hash1", &mock.PubKeyConverterStub{})
}

func TestLogsProcessor_ProcessLogs_ThreeLogs_FourEvents_ExpectTwoProcessedLogAndThreeEvents(t *testing.T) {
	lp, _ := logs.NewLogsProcessor(&mock.PubKeyConverterStub{})

	event1 := generateRandomEvent()
	event2 := generateRandomEvent()
	event3 := generateRandomEvent()
	log1 := &transaction.Log{
		Address: []byte(strconv.Itoa(rand.Int())),
		Events:  []*transaction.Event{event1, nil, event2},
	}

	log2 := &transaction.Log{
		Address: []byte(strconv.Itoa(rand.Int())),
		Events:  []*transaction.Event{event3},
	}

	logsAndEvents := map[string]data.LogHandler{
		"hash1": log1,
		"hash2": nil,
		"hash3": log2,
	}

	ret, _ := lp.ProcessLogs(logsAndEvents)
	require.Len(t, ret, 2)
	require.Len(t, ret[0].Events, 2)
	require.Len(t, ret[1].Events, 1)

	requireProcessedLogEqual(t, ret[0], log1, "hash1", &mock.PubKeyConverterStub{})
	requireProcessedLogEqual(t, ret[1], log2, "hash3", &mock.PubKeyConverterStub{})
}

func generateRandomEvent() *transaction.Event {
	return &transaction.Event{
		Address:    []byte(strconv.Itoa(rand.Int())),
		Identifier: []byte(strconv.Itoa(rand.Int())),
		Topics:     [][]byte{[]byte(strconv.Itoa(rand.Int())), []byte(strconv.Itoa(rand.Int()))},
		Data:       []byte(strconv.Itoa(rand.Int())),
	}
}

func requireProcessedLogEqual(
	t *testing.T,
	processedLog *schema.Log,
	log *transaction.Log,
	hash string,
	pubKeyConverter core.PubkeyConverter) {

	require.Equal(t, processedLog.ID, []byte(hash))
	require.Equal(t, processedLog.Address, utility.EncodePubKey(pubKeyConverter, log.GetAddress()))

	notNilEvents := getNotNilEvents(log.GetEvents())
	for idx, currEvent := range notNilEvents {
		requireProcessedEventEqual(t, processedLog.Events[idx], currEvent, pubKeyConverter)
	}
}

func getNotNilEvents(events []*transaction.Event) []*transaction.Event {
	ret := make([]*transaction.Event, 0)

	for _, event := range events {
		if !check.IfNil(event) {
			ret = append(ret, event)
		}
	}

	return ret
}

func requireProcessedEventEqual(
	t *testing.T,
	processedEvent *schema.Event,
	event *transaction.Event,
	pubKeyConverter core.PubkeyConverter) {

	require.Equal(t, processedEvent.Address, utility.EncodePubKey(pubKeyConverter, event.GetAddress()))
	require.Equal(t, processedEvent.Identifier, event.GetIdentifier())
	require.Equal(t, processedEvent.Topics, event.GetTopics())
	require.Equal(t, processedEvent.Data, event.GetData())
}
