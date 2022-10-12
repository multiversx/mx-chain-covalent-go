package process

import (
	covalentApi "github.com/ElrondNetwork/covalent-indexer-go/api"
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

// AlteredAccountsHandler defines what an account processor shall do
type AlteredAccountsHandler interface {
	ProcessAccounts(apiAlteredAccounts []*covalentApi.AlteredAccount) ([]*schemaV2.AccountBalanceUpdate, error)
}
