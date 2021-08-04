package factory

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/accounts"
	blockCovalent "github.com/ElrondNetwork/covalent-indexer-go/process/block"
	"github.com/ElrondNetwork/covalent-indexer-go/process/logs"
	"github.com/ElrondNetwork/covalent-indexer-go/process/receipts"
	"github.com/ElrondNetwork/covalent-indexer-go/process/transactions"
	"github.com/ElrondNetwork/elrond-go-core/core"
)

type ArgsDataProcessor struct {
	PubKeyConvertor core.PubkeyConverter
	Accounts        covalent.AccountsAdapter
}

func CreateDataProcessor(args *ArgsDataProcessor) (*process.DataProcessor, error) {
	blockHandler, _ := blockCovalent.NewBlockProcessor()
	transactionsHandler, _ := transactions.NewTransactionProcessor()
	receiptsHandler, _ := receipts.NewReceiptsProcessor()
	scHandler, _ := transactions.NewSCProcessor()
	logHandler, _ := logs.NewLogsProcessor()
	accountsHandler, _ := accounts.NewAccountsProcessor(&args.Accounts, &args.PubKeyConvertor)

	return process.NewDataProcessor(
		blockHandler,
		transactionsHandler,
		scHandler,
		receiptsHandler,
		logHandler,
		accountsHandler)
}
