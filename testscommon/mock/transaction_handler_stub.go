package mock

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

type TransactionHandlerStub struct {
	ProcessTransactionsCalled func(apiTransactions []*transaction.ApiTransactionResult) ([]*schemaV2.Transaction, error)
}

func (ths *TransactionHandlerStub) ProcessTransactions(apiTransactions []*transaction.ApiTransactionResult) ([]*schemaV2.Transaction, error) {
	if ths.ProcessTransactionsCalled != nil {
		return ths.ProcessTransactionsCalled(apiTransactions)
	}

	return nil, nil
}
