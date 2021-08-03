package transactions

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type scProcessor struct{}

func NewSCProcessor() (*scProcessor, error) {
	return &scProcessor{}, nil
}

func (b *scProcessor) ProcessSCs(transactions *map[string]data.TransactionHandler) ([]*schema.SCResult, error) {
	return nil, nil
}
