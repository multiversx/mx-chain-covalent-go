package transactions

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type scProcessor struct{}

// NewSCProcessor creates a new instance of smart contracts processor
func NewSCProcessor() (*scProcessor, error) {
	return &scProcessor{}, nil
}

// ProcessSCs converts smart contracts data to a specific structure defined by avro schema
func (scp *scProcessor) ProcessSCs(transactions map[string]data.TransactionHandler) ([]*schema.SCResult, error) {
	return nil, nil
}
