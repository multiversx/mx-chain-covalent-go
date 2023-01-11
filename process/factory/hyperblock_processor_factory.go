package factory

import (
	"github.com/multiversx/mx-chain-covalent-go"
	"github.com/multiversx/mx-chain-covalent-go/process"
	"github.com/multiversx/mx-chain-covalent-go/process/accounts"
	"github.com/multiversx/mx-chain-covalent-go/process/epochStart"
	"github.com/multiversx/mx-chain-covalent-go/process/logs"
	"github.com/multiversx/mx-chain-covalent-go/process/receipts"
	"github.com/multiversx/mx-chain-covalent-go/process/shardBlocks"
	"github.com/multiversx/mx-chain-covalent-go/process/transactions"
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
