package factory

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/logs"
	"github.com/ElrondNetwork/covalent-indexer-go/process/receipts"
	"github.com/ElrondNetwork/covalent-indexer-go/process/transactions"
)

// ArgsHyperBlockProcessor holds all input dependencies required by hyper block processor factory
// in order to create a new hyper block processor
type ArgsHyperBlockProcessor struct {
}

// CreateHyperBlockProcessor creates a new hyper block processor handler
func CreateHyperBlockProcessor(args *ArgsHyperBlockProcessor) (covalent.HyperBlockProcessor, error) {
	receiptsHandler := receipts.NewReceiptsProcessor()
	logsHandler := logs.NewLogsProcessor()
	transactionsHandler, err := transactions.NewTransactionProcessor(logsHandler, receiptsHandler)
	if err != nil {
		return nil, err
	}

	return process.NewHyperBlockProcessor(transactionsHandler)
}
