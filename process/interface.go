package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

// TransactionHandler defines what a transaction processor shall do
type TransactionHandler interface {
	ProcessTransactions(apiTransactions []*transaction.ApiTransactionResult) ([]*schemaV2.Transaction, error)
}

// ReceiptHandler defines what a receipt processor shall do
type ReceiptHandler interface {
	ProcessReceipt(apiReceipt *transaction.ApiReceipt) (*schemaV2.Receipt, error)
}

// LogHandler defines what a log processor shall do
type LogHandler interface {
	ProcessLog(log *transaction.ApiLogs) *schemaV2.Log
}

// ShardBlocksHandler defines what shard blocks processor shall do
type ShardBlocksHandler interface {
	ProcessShardBlocks(apiBlocks []*api.NotarizedBlock) ([]*schemaV2.ShardBlocks, error)
}

// EpochStartInfoHandler defines what epoch start info processor shall do
type EpochStartInfoHandler interface {
	ProcessEpochStartInfo(apiEpochInfo *api.EpochStartInfo) (*schemaV2.EpochStartInfo, error)
}

// AccountsHandler defines what an account processor shall do
type AccountsHandler interface {
	ProcessAccounts(
		processedTxs []*schema.Transaction,
		processedSCRs []*schema.SCResult,
		processedReceipts []*schema.Receipt) []*schema.AccountBalanceUpdate
}

// ShardCoordinator defines what a shard coordinator shall do
type ShardCoordinator interface {
	SelfId() uint32
	ComputeId(address []byte) uint32
	IsInterfaceNil() bool
}
