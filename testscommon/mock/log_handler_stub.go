package mock

import (
	"github.com/multiversx/mx-chain-covalent-go/schema"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
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
