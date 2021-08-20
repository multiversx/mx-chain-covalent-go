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

// MiniBlockHandler defines what a mini blocks processor shall do
type MiniBlockHandler interface {
	ProcessMiniBlocks(header data.HeaderHandler, body data.BodyHandler) ([]*schema.MiniBlock, error)
}

// TransactionHandler defines what a transaction processor shall do
type TransactionHandler interface {
	ProcessTransactions(
		header data.HeaderHandler,
		headerHash []byte,
		bodyHandler data.BodyHandler,
		transactions map[string]data.TransactionHandler) ([]*schema.Transaction, error)
}

// SCHandler defines what a smart contract processor shall do
type SCHandler interface {
	ProcessSCs(transactions map[string]data.TransactionHandler, timeStamp uint64) []*schema.SCResult
}

// ReceiptHandler defines what a receipt processor shall do
type ReceiptHandler interface {
	ProcessReceipts(receipts map[string]data.TransactionHandler, timeStamp uint64) []*schema.Receipt
}

// LogHandler defines what a log processor shall do
type LogHandler interface {
	ProcessLogs(logs map[string]data.LogHandler) ([]*schema.Log, error)
}

// AccountsHandler defines what an account processor shall do
type AccountsHandler interface {
	ProcessAccounts() ([]*schema.AccountBalanceUpdate, error)
}
