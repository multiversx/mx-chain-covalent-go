package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

// TransactionHandler defines what a transaction processor shall do
type TransactionHandler interface {
	ProcessTransactions(apiTransactions []*transaction.ApiTransactionResult) ([]*schemaV2.Transaction, error)
}

// SCResultsHandler defines what a smart contract processor shall do
type SCResultsHandler interface {
	ProcessSCRs(transactions map[string]data.TransactionHandler, timeStamp uint64) []*schema.SCResult
}

// ReceiptHandler defines what a receipt processor shall do
type ReceiptHandler interface {
	ProcessReceipt(apiReceipt *transaction.ApiReceipt) (*schemaV2.Receipt, error)
}

// LogHandler defines what a log processor shall do
type LogHandler interface {
	ProcessLog(log *transaction.ApiLogs) *schemaV2.Log
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
