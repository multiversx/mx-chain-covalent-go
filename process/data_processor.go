package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type DataProcessor struct {
	blockHandler       BlockHandler
	transactionHandler TransactionHandler
	receiptHandler     ReceiptHandler
}

func NewDataProcessor(blockHandler BlockHandler,
	transactionHandler TransactionHandler,
	receiptHandler ReceiptHandler) (*DataProcessor, error) {

	return &DataProcessor{
		blockHandler:       blockHandler,
		transactionHandler: transactionHandler,
		receiptHandler:     receiptHandler,
	}, nil
}

func (d *DataProcessor) ProcessData(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
	block, err := d.blockHandler.ProcessBlock(&args.Body)
	if err != nil {
		return nil, err
	}

	transactions, _ := d.transactionHandler.ProcessTransactions(&args.TransactionsPool.Txs)
	receipts, _ := d.receiptHandler.ProcessReceipts(&args.TransactionsPool.Txs)

	return &schema.BlockResult{
		Block:        block,
		Transactions: transactions,
		Receipts:     receipts,
	}, nil
}
