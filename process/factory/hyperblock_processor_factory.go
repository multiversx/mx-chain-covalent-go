package factory

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/accounts"
	"github.com/ElrondNetwork/covalent-indexer-go/process/epochStart"
	"github.com/ElrondNetwork/covalent-indexer-go/process/logs"
	"github.com/ElrondNetwork/covalent-indexer-go/process/receipts"
	"github.com/ElrondNetwork/covalent-indexer-go/process/shardBlocks"
	"github.com/ElrondNetwork/covalent-indexer-go/process/transactions"
)

// CreateHyperBlockProcessor creates a new hyper block processor handler
func CreateHyperBlockProcessor() (covalent.HyperBlockProcessor, error) {
	receiptsHandler := receipts.NewReceiptsProcessor()
	logsHandler := logs.NewLogsProcessor()
	transactionsHandler, err := transactions.NewTransactionProcessor(logsHandler, receiptsHandler)
	if err != nil {
		return nil, err
	}

	shardBlocksHandler := shardBlocks.NewShardBlocksProcessor()
	epochStartInfoHandler := epochStart.NewEpochStartInfoProcessor()
	alteredAccountsHandler := accounts.NewAlteredAccountsProcessor()

	args := &process.HyperBlockProcessorArgs{
		TransactionHandler:     transactionsHandler,
		ShardBlockHandler:      shardBlocksHandler,
		EpochStartInfoHandler:  epochStartInfoHandler,
		AlteredAccountsHandler: alteredAccountsHandler,
	}
	return process.NewHyperBlockProcessor(args)
}
