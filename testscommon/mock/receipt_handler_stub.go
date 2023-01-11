package mock

import (
	"github.com/multiversx/mx-chain-covalent-go/schema"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
)

// ReceiptHandlerStub -
type ReceiptHandlerStub struct {
	ProcessReceiptCalled func(apiReceipt *transaction.ApiReceipt) (*schema.Receipt, error)
}

// ProcessReceipt -
func (rhs *ReceiptHandlerStub) ProcessReceipt(apiReceipt *transaction.ApiReceipt) (*schema.Receipt, error) {
	if rhs.ProcessReceiptCalled != nil {
		return rhs.ProcessReceiptCalled(apiReceipt)
	}

	return nil, nil
}
