package processMocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
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
