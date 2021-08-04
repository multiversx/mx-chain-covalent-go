package transactions

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type transactionProcessor struct{}

func NewTransactionProcessor() (*transactionProcessor, error) {
	return &transactionProcessor{}, nil
}

func (txp *transactionProcessor) ProcessTransactions(transactions *map[string]data.TransactionHandler) ([]*schema.Transaction, error) {
	return nil, nil
}
