package receipts

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/receipt"
)

type receiptsProcessor struct {
	pubKeyConverter core.PubkeyConverter
}

// NewReceiptsProcessor creates a new instance of receipts processor
func NewReceiptsProcessor(pubKeyConverter core.PubkeyConverter) (*receiptsProcessor, error) {
	if check.IfNil(pubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}

	return &receiptsProcessor{pubKeyConverter: pubKeyConverter}, nil
}

// ProcessReceipts converts receipts data to a specific structure defined by avro schema
func (rp *receiptsProcessor) ProcessReceipts(receipts map[string]data.TransactionHandler, timeStamp uint64) ([]*schema.Receipt, error) {
	allReceipts := make([]*schema.Receipt, 0)

	for currHash, currReceipt := range receipts {

		rec, _ := currReceipt.(*receipt.Receipt)
		receipt := &schema.Receipt{
			Hash:      []byte(currHash),
			Value:     utility.GetBytes(rec.GetValue()),
			Sender:    utility.EncodePubKey(rp.pubKeyConverter, rec.GetSndAddr()),
			Data:      rec.GetData(),
			TxHash:    rec.GetTxHash(),
			Timestamp: int64(timeStamp),
		}

		allReceipts = append(allReceipts, receipt)
	}

	return allReceipts, nil
}
