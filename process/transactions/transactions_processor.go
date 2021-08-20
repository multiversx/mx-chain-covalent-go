package transactions

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type transactionProcessor struct{}

// NewTransactionProcessor creates a new instance of transactions processor
func NewTransactionProcessor() (*transactionProcessor, error) {
	return &transactionProcessor{}, nil
}

// ProcessTransactions converts transactions data to a specific structure defined by avro schema
func (txp *transactionProcessor) ProcessTransactions(transactions map[string]data.TransactionHandler) ([]*schema.Transaction, error) {
	return nil, nil
}
