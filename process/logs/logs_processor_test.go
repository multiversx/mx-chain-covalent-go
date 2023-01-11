package logs_test

import (
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/process/logs"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/stretchr/testify/require"
)

func TestLogsProcessor_ProcessLog(t *testing.T) {
	t.Parallel()

	lp := logs.NewLogsProcessor()

	t.Run("no log, expect empty log", func(t *testing.T) {
		t.Parallel()

		log := lp.ProcessLog(nil)
		require.Equal(t, schema.NewLog(), log)
	})

	t.Run("no events, expect log only filled with address", func(t *testing.T) {
		t.Parallel()

		apiLog := &transaction.ApiLogs{Address: "erd1qq", Events: nil}

		processedLog := lp.ProcessLog(apiLog)
		require.Equal(t, &schema.Log{
			Address: []byte(apiLog.Address),
			Events:  []*schema.Event{},
		}, processedLog)
	})

	t.Run("4 events, 2 nil events, expect one log with 2 events", func(t *testing.T) {
		t.Parallel()

		event1 := &transaction.Events{
			Address:    "erd1aa",
			Identifier: "id1",
			Topics:     [][]byte{testscommon.GenerateRandomBytes()},
			Data:       testscommon.GenerateRandomBytes(),
		}

		event2 := &transaction.Events{
			Address:    "erd1bb",
			Identifier: "id2",
			Topics:     [][]byte{testscommon.GenerateRandomBytes()},
			Data:       testscommon.GenerateRandomBytes(),
		}

		apiLog := &transaction.ApiLogs{
			Address: "erd1cc",
			Events:  []*transaction.Events{event1, nil, event2, nil},
		}

		processedLog := lp.ProcessLog(apiLog)
		expectedLog := &schema.Log{
			Address: []byte(apiLog.Address),
			Events: []*schema.Event{
				{
					Address:    []byte(event1.Address),
					Identifier: []byte(event1.Identifier),
					Topics:     event1.Topics,
					Data:       event1.Data,
				},
				{
					Address:    []byte(event2.Address),
					Identifier: []byte(event2.Identifier),
					Topics:     event2.Topics,
					Data:       event2.Data,
				},
			},
		}
		require.Equal(t, expectedLog, processedLog)
	})
}
