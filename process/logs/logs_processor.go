package logs

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

type logsProcessor struct {
}

// NewLogsProcessor creates a new instance of logs processor
func NewLogsProcessor() *logsProcessor {
	return &logsProcessor{}
}

// ProcessLog converts logs data to a specific structure defined by avro schema
func (lp *logsProcessor) ProcessLog(log *transaction.ApiLogs) *schema.Log {
	if log == nil {
		return schema.NewLog()
	}

	return &schema.Log{
		Address: []byte(log.Address),
		Events:  lp.processEvents(log.Events),
	}
}

func (lp *logsProcessor) processEvents(events []*transaction.Events) []*schema.Event {
	allEvents := make([]*schema.Event, 0, len(events))

	for _, currEvent := range events {
		processedEvent := lp.processEvent(currEvent)

		if processedEvent != nil {
			allEvents = append(allEvents, processedEvent)
		}
	}

	return allEvents
}

func (lp *logsProcessor) processEvent(event *transaction.Events) *schema.Event {
	if event == nil {
		return nil
	}

	return &schema.Event{
		Address:    []byte(event.Address),
		Identifier: []byte(event.Identifier),
		Topics:     event.Topics,
		Data:       event.Data,
	}
}
