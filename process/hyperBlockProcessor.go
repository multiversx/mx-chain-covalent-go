package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
)

type hyperBlockProcessor struct {
	transactionProcessor TransactionHandler
}

// NewHyperBlockProcessor will create a new instance of an hyper block processor
func NewHyperBlockProcessor(transactionHandler TransactionHandler) (*hyperBlockProcessor, error) {
	if transactionHandler == nil {
		return nil, errNilTransactionHandler
	}

	return &hyperBlockProcessor{
		transactionProcessor: transactionHandler,
	}, nil
}

// Process will process current hyper block and convert it to an avro schema block result
func (hbp *hyperBlockProcessor) Process(hyperBlock *hyperBlock.HyperBlock) (*schemaV2.HyperBlock, error) {
	avroHyperBlock := schemaV2.NewHyperBlock()

	txs, err := hbp.transactionProcessor.ProcessTransactions(hyperBlock.Transactions)
	if err != nil {
		return nil, err
	}

	avroHyperBlock.Transactions = txs
	return avroHyperBlock, nil
}
