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
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
)

// ArgsDataProcessor holds all input dependencies required by data processor factory
// in order to create a new data handler instance of type data processor
type ArgsDataProcessor struct {
	PubKeyConvertor core.PubkeyConverter
	Accounts        covalent.AccountsAdapter
	hasher          hashing.Hasher
	marshalizer     marshal.Marshalizer
}

// CreateDataProcessor creates a new data handler instance of type data processor
func CreateDataProcessor(args *ArgsDataProcessor) (covalent.DataHandler, error) {
	blockHandler, err := blockCovalent.NewBlockProcessor(args.hasher, args.marshalizer)
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

	accountsHandler, err := accounts.NewAccountsProcessor(args.Accounts, args.PubKeyConvertor)
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
