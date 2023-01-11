package receipts

import (
	"encoding/hex"
	"fmt"

	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

type receiptsProcessor struct {
	pubKeyConverter core.PubkeyConverter
}

// NewReceiptsProcessor creates a new instance of receipts processor
func NewReceiptsProcessor() *receiptsProcessor {
	return &receiptsProcessor{}
}

// ProcessReceipt converts receipts api data to a specific structure defined by avro schema
func (rp *receiptsProcessor) ProcessReceipt(apiReceipt *transaction.ApiReceipt) (*schema.Receipt, error) {
	if apiReceipt == nil {
		return schema.NewReceipt(), nil
	}

	hash, err := hex.DecodeString(apiReceipt.TxHash)
	if err != nil {
		return nil, fmt.Errorf("receiptsProcessor.ProcessReceipt: could not decode tx hash: %s from receipt, err: %w", apiReceipt.TxHash, err)
	}

	return &schema.Receipt{
		Value:  utility.GetBytes(apiReceipt.Value),
		Sender: []byte(apiReceipt.SndAddr),
		Data:   []byte(apiReceipt.Data),
		TxHash: hash,
	}, nil
}
