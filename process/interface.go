package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

// BlockHandler defines what a block processor shall do
type BlockHandler interface {
	ProcessBlock(args *indexer.ArgsSaveBlockData) (*schema.Block, error)
}

// TransactionHandler defines what a transaction processor shall do
type TransactionHandler interface {
	ProcessTransactions(transactions map[string]data.TransactionHandler) ([]*schema.Transaction, error)
}

// SCHandler defines what a smart contract processor shall do
type SCHandler interface {
	ProcessSCs(transactions map[string]data.TransactionHandler) ([]*schema.SCResult, error)
}

// ReceiptHandler defines what a receipt processor shall do
type ReceiptHandler interface {
	ProcessReceipts(transactions map[string]data.TransactionHandler) ([]*schema.Receipt, error)
}

// LogHandler defines what a log processor shall do
type LogHandler interface {
	ProcessLogs(logs map[string]data.LogHandler) ([]*schema.Log, error)
}

// AccountsHandler defines what an account processor shall do
type AccountsHandler interface {
	ProcessAccounts() ([]*schema.AccountBalanceUpdate, error)
}
