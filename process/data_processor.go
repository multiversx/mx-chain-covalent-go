package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type dataProcessor struct {
	blockHandler       BlockHandler
	transactionHandler TransactionHandler
	receiptHandler     ReceiptHandler
	scHandler          SCHandler
	logHandler         LogHandler
	accountsHandler    AccountsHandler
}

// NewDataProcessor creates a new instance of data processor, which handles all sub-processes
func NewDataProcessor(
	blockHandler BlockHandler,
	transactionHandler TransactionHandler,
	scHandler SCHandler,
	receiptHandler ReceiptHandler,
	logHandler LogHandler,
	accountsHandler AccountsHandler,
) (*dataProcessor, error) {

	return &dataProcessor{
		blockHandler:       blockHandler,
		transactionHandler: transactionHandler,
		scHandler:          scHandler,
		receiptHandler:     receiptHandler,
		logHandler:         logHandler,
		accountsHandler:    accountsHandler,
	}, nil
}

// ProcessData converts all covalent necessary data to a specific structure defined by avro schema
func (dp *dataProcessor) ProcessData(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {

	block, err := dp.blockHandler.ProcessBlock(args)
	if err != nil {
		return nil, err
	}

	transactions, err := dp.transactionHandler.ProcessTransactions(args.Header, args.HeaderHash, args.Body, args.TransactionsPool.Txs)
	if err != nil {
		return nil, err
	}

	smartContracts := dp.scHandler.ProcessSCs(args.TransactionsPool.Scrs, args.Header.GetTimeStamp())
	receipts := dp.receiptHandler.ProcessReceipts(args.TransactionsPool.Receipts, args.Header.GetTimeStamp())
	logs := dp.logHandler.ProcessLogs(args.TransactionsPool.Logs)

	accountUpdates, err := dp.accountsHandler.ProcessAccounts()
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
