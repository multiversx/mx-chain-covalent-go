package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type DataProcessor struct {
	blockHandler       BlockHandler
	transactionHandler TransactionHandler
}

func NewDataProcessor(blockHandler BlockHandler,
	transactionHandler TransactionHandler) (*DataProcessor, error) {
	return &DataProcessor{
		blockHandler:       blockHandler,
		transactionHandler: transactionHandler,
	}, nil
}

func (d *DataProcessor) ProcessData(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
	block, err := d.blockHandler.ProcessBlock(&args.Body)
	if err != nil {
		return nil, err
	}

	transactions, _ := d.transactionHandler.ProcessTransactions(&args.TransactionsPool.Txs)

	return &schema.BlockResult{
		Block:        block,
		Transactions: transactions,
	}, nil
}
