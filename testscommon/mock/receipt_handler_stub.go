package mock

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

// ReceiptHandlerStub -
type ReceiptHandlerStub struct {
	ProcessReceiptCalled func(apiReceipt *transaction.ApiReceipt) (*schemaV2.Receipt, error)
}

// ProcessReceipt -
func (rhs *ReceiptHandlerStub) ProcessReceipt(apiReceipt *transaction.ApiReceipt) (*schemaV2.Receipt, error) {
	if rhs.ProcessReceiptCalled != nil {
		return rhs.ProcessReceiptCalled(apiReceipt)
	}

	return nil, nil
}
