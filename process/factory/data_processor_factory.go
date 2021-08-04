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

func CreateDataProcessor(args *ArgsDataProcessor) (covalent.DataHandler, error) {
	blockHandler, err := blockCovalent.NewBlockProcessor()
	if err != nil {
		return nil, err
	}

	transactionsHandler, err := transactions.NewTransactionProcessor()
	if err != nil {
		return nil, err
	}

	receiptsHandler, err := receipts.NewReceiptsProcessor()
	if err != nil {
		return nil, err
	}

	scHandler, err := transactions.NewSCProcessor()
	if err != nil {
		return nil, err
	}

	logHandler, err := logs.NewLogsProcessor()
	if err != nil {
		return nil, err
	}

	accountsHandler, err := accounts.NewAccountsProcessor(&args.Accounts, &args.PubKeyConvertor)
	if err != nil {
		return nil, err
	}

	return process.NewDataProcessor(
		blockHandler,
		transactionsHandler,
		scHandler,
		receiptsHandler,
		logHandler,
		accountsHandler)
}
