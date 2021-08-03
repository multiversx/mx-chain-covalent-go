package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type DataProcessor struct {
	blockHandler       BlockHandler
	transactionHandler TransactionHandler
	receiptHandler     ReceiptHandler
	scHandler          SCHandler
}

func NewDataProcessor(
	blockHandler BlockHandler,
	transactionHandler TransactionHandler,
	scHandler SCHandler,
	receiptHandler ReceiptHandler) (*DataProcessor, error) {
	return &DataProcessor{
		blockHandler:       blockHandler,
		transactionHandler: transactionHandler,
		scHandler:          scHandler,
		receiptHandler:     receiptHandler,
	}, nil
}

func (d *DataProcessor) ProcessData(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {

	block, _ := d.blockHandler.ProcessBlock(&args.Body)
	transactions, _ := d.transactionHandler.ProcessTransactions(&args.TransactionsPool.Txs)
	smartContracts, _ := d.scHandler.ProcessSCs(&args.TransactionsPool.Scrs)
	receipts, _ := d.receiptHandler.ProcessReceipts(&args.TransactionsPool.Receipts)

	return &schema.BlockResult{
		Block:        block,
		Transactions: transactions,
		Receipts:     receipts,
		SCResults:    smartContracts,
	}, nil
}
