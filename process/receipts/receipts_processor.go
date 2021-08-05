package receipts

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type receiptsProcessor struct{}

// NewReceiptsProcessor creates a new instance of receipts processor
func NewReceiptsProcessor() (*receiptsProcessor, error) {
	return &receiptsProcessor{}, nil
}

// ProcessReceipts converts receipts data to a specific structure defined by avro schema
func (rp *receiptsProcessor) ProcessReceipts(transactions map[string]data.TransactionHandler) ([]*schema.Receipt, error) {
	return nil, nil
}
