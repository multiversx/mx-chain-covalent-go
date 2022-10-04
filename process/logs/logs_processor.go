package logs

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

type logsProcessor struct {
}

// NewLogsProcessor creates a new instance of logs processor
func NewLogsProcessor() *logsProcessor {
	return &logsProcessor{}
}

// ProcessLog converts logs data to a specific structure defined by avro schema
func (lp *logsProcessor) ProcessLog(log *transaction.ApiLogs) *schemaV2.Log {
	if log == nil {
		return schemaV2.NewLog()
	}

	return &schemaV2.Log{
		Address: []byte(log.Address),
		Events:  lp.processEvents(log.Events),
	}
}

func (lp *logsProcessor) processEvents(events []*transaction.Events) []*schemaV2.Event {
	allEvents := make([]*schemaV2.Event, 0, len(events))

	for _, currEvent := range events {
		processedEvent := lp.processEvent(currEvent)

		if processedEvent != nil {
			allEvents = append(allEvents, processedEvent)
		}
	}

	return allEvents
}

func (lp *logsProcessor) processEvent(event *transaction.Events) *schemaV2.Event {
	if event == nil {
		return nil
	}

	return &schemaV2.Event{
		Address:    []byte(event.Address),
		Identifier: []byte(event.Identifier),
		Topics:     event.Topics,
		Data:       event.Data,
	}
}
