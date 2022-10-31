package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

// TransactionHandler defines what a transaction processor shall do
type TransactionHandler interface {
	ProcessTransactions(apiTransactions []*transaction.ApiTransactionResult) ([]*schema.Transaction, error)
}

// ReceiptHandler defines what a receipt processor shall do
type ReceiptHandler interface {
	ProcessReceipt(apiReceipt *transaction.ApiReceipt) (*schema.Receipt, error)
}

// LogHandler defines what a log processor shall do
type LogHandler interface {
	ProcessLog(log *transaction.ApiLogs) *schema.Log
}

// ShardBlocksHandler defines what shard blocks processor shall do
type ShardBlocksHandler interface {
	ProcessShardBlocks(apiBlocks []*api.NotarizedBlock) ([]*schema.ShardBlocks, error)
}

// EpochStartInfoHandler defines what epoch start info processor shall do
type EpochStartInfoHandler interface {
	ProcessEpochStartInfo(apiEpochInfo *api.EpochStartInfo) (*schema.EpochStartInfo, error)
}

// AlteredAccountsHandler defines what an account processor shall do
type AlteredAccountsHandler interface {
	ProcessAccounts(notarizedBlocks []*api.NotarizedBlock) ([]*schema.AccountBalanceUpdate, error)
}
