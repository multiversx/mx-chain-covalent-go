package factory

import (
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	blockCovalent "github.com/ElrondNetwork/covalent-indexer-go/process/block"
	"github.com/ElrondNetwork/covalent-indexer-go/process/logs"
	"github.com/ElrondNetwork/covalent-indexer-go/process/receipts"
	"github.com/ElrondNetwork/covalent-indexer-go/process/transactions"
	"github.com/ElrondNetwork/elrond-go-core/core"
)

type ArgsDataProcessor struct {
	PubKeyConvertor core.PubkeyConverter
}

func CreateDataProcessor(args *ArgsDataProcessor) (*process.DataProcessor, error) {
	blockHandler, _ := blockCovalent.NewBlockProcessor()
	transactionsHandler, _ := transactions.NewTransactionProcessor()
	receiptsHandler, _ := receipts.NewReceiptsProcessor()
	scHandler, _ := transactions.NewSCProcessor()
	logHandler, _ := logs.NewLogsProcessor()

	return process.NewDataProcessor(
		blockHandler,
		transactionsHandler,
		scHandler,
		receiptsHandler,
		logHandler)
}
