package receipts

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type receiptsProcessor struct{}

func NewReceiptsProcessor() (*receiptsProcessor, error) {
	return &receiptsProcessor{}, nil
}

func (b *receiptsProcessor) ProcessReceipts(transactions *map[string]data.TransactionHandler) ([]*schema.Receipt, error) {
	return nil, nil
}
