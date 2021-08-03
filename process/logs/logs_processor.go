package logs

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type logsProcessor struct{}

func NewLogsProcessor() (*logsProcessor, error) {
	return &logsProcessor{}, nil
}

func (b *logsProcessor) ProcessLogs(logs *map[string]data.LogHandler) ([]*schema.Log, error) {
	return nil, nil
}
