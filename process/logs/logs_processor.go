package logs

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type logsProcessor struct {
	pubKeyConverter core.PubkeyConverter
}

// NewLogsProcessor creates a new instance of logs processor
func NewLogsProcessor(pubKeyConverter core.PubkeyConverter) (*logsProcessor, error) {
	if check.IfNil(pubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}

	return &logsProcessor{
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// ProcessLogs converts logs data to a specific structure defined by avro schema
func (lp *logsProcessor) ProcessLogs(logs map[string]data.LogHandler) []*schema.Log {
	allLogs := make([]*schema.Log, 0, len(logs))

	for currHash, currLog := range logs {
		processedLog := lp.processLog(currHash, currLog)
		if processedLog != nil {
			allLogs = append(allLogs, processedLog)
		}
	}

	return allLogs
}

func (lp *logsProcessor) processLog(hash string, log data.LogHandler) *schema.Log {
	if check.IfNil(log) {
		return nil
	}

	return &schema.Log{
		ID:      []byte(hash),
		Address: utility.EncodePubKey(lp.pubKeyConverter, log.GetAddress()),
		Events:  lp.processEvents(log.GetLogEvents()),
	}
}

func (lp *logsProcessor) processEvents(events []data.EventHandler) []*schema.Event {
	allEvents := make([]*schema.Event, 0, len(events))

	for _, currEvent := range events {
		processedEvent := lp.processEvent(currEvent)

		if processedEvent != nil {
			allEvents = append(allEvents, processedEvent)
		}
	}

	return allEvents
}

func (lp *logsProcessor) processEvent(event data.EventHandler) *schema.Event {
	if check.IfNil(event) {
		return nil
	}

	return &schema.Event{
		Address:    utility.EncodePubKey(lp.pubKeyConverter, event.GetAddress()),
		Identifier: event.GetIdentifier(),
		Topics:     event.GetTopics(),
		Data:       event.GetData(),
	}
}
