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

	alteredAccountsHandler := accounts.NewAlteredAccountsProcessor()
	shardBlocksHandler, err := shardBlocks.NewShardBlocksProcessor(alteredAccountsHandler)
	if err != nil {
		return nil, err
	}

	epochStartInfoHandler := epochStart.NewEpochStartInfoProcessor()
	args := &process.HyperBlockProcessorArgs{
		TransactionHandler:    transactionsHandler,
		ShardBlockHandler:     shardBlocksHandler,
		EpochStartInfoHandler: epochStartInfoHandler,
	}
	return process.NewHyperBlockProcessor(args)
}
