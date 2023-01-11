package processMocks

import (
	"github.com/multiversx/mx-chain-covalent-go/schema"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
)

// TransactionHandlerStub -
type TransactionHandlerStub struct {
	ProcessTransactionsCalled func(apiTransactions []*transaction.ApiTransactionResult) ([]*schema.Transaction, error)
}

// ProcessTransactions -
func (ths *TransactionHandlerStub) ProcessTransactions(apiTransactions []*transaction.ApiTransactionResult) ([]*schema.Transaction, error) {
	if ths.ProcessTransactionsCalled != nil {
		return ths.ProcessTransactionsCalled(apiTransactions)
	}

	return nil, nil
}
