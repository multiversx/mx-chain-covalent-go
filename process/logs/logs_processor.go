package logs

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type logsProcessor struct{}

// NewLogsProcessor creates a new instance of logs processor
func NewLogsProcessor() (*logsProcessor, error) {
	return &logsProcessor{}, nil
}

// ProcessLogs converts logs data to a specific structure defined by avro schema
func (lp *logsProcessor) ProcessLogs(logs *map[string]data.LogHandler) ([]*schema.Log, error) {
	return nil, nil
}
