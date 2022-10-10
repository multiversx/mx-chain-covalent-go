package mock

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

// LogHandlerStub -
type LogHandlerStub struct {
	ProcessLogCalled func(log *transaction.ApiLogs) *schemaV2.Log
}

// ProcessLog -
func (lhs *LogHandlerStub) ProcessLog(log *transaction.ApiLogs) *schemaV2.Log {
	if lhs.ProcessLogCalled != nil {
		return lhs.ProcessLogCalled(log)
	}

	return nil
}
