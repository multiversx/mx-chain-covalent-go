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
	logHandler         LogHandler
	accountsHandler    AccountsHandler
}

func NewDataProcessor(
	blockHandler BlockHandler,
	transactionHandler TransactionHandler,
	scHandler SCHandler,
	receiptHandler ReceiptHandler,
	logHandler LogHandler,
	accountsHandler AccountsHandler,
) (*DataProcessor, error) {

	return &DataProcessor{
		blockHandler:       blockHandler,
		transactionHandler: transactionHandler,
		scHandler:          scHandler,
		receiptHandler:     receiptHandler,
		logHandler:         logHandler,
		accountsHandler:    accountsHandler,
	}, nil
}

func (d *DataProcessor) ProcessData(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {

	block, err := d.blockHandler.ProcessBlock(&args.Body)
	if err != nil {
		return nil, err
	}

	transactions, err := d.transactionHandler.ProcessTransactions(&args.TransactionsPool.Txs)
	if err != nil {
		return nil, err
	}

	smartContracts, err := d.scHandler.ProcessSCs(&args.TransactionsPool.Scrs)
	if err != nil {
		return nil, err
	}

	receipts, err := d.receiptHandler.ProcessReceipts(&args.TransactionsPool.Receipts)
	if err != nil {
		return nil, err
	}

	logs, err := d.logHandler.ProcessLogs(&args.TransactionsPool.Logs)
	if err != nil {
		return nil, err
	}

	accountUpdates, err := d.accountsHandler.ProcessAccounts()
	if err != nil {
		return nil, err
	}

	return &schema.BlockResult{
		Block:        block,
		Transactions: transactions,
		Receipts:     receipts,
		SCResults:    smartContracts,
		Logs:         logs,
		StateChanges: accountUpdates,
	}, nil
}
