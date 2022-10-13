package mock

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

// LogHandlerStub -
type LogHandlerStub struct {
	ProcessLogCalled func(log *transaction.ApiLogs) *schema.Log
}

// ProcessLog -
func (lhs *LogHandlerStub) ProcessLog(log *transaction.ApiLogs) *schema.Log {
	if lhs.ProcessLogCalled != nil {
		return lhs.ProcessLogCalled(log)
	}

	return nil
}
