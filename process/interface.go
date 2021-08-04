package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type BlockHandler interface {
	ProcessBlock(block *data.BodyHandler) (*schema.Block, error)
}

type TransactionHandler interface {
	ProcessTransactions(transactions *map[string]data.TransactionHandler) ([]*schema.Transaction, error)
}

type SCHandler interface {
	ProcessSCs(transactions *map[string]data.TransactionHandler) ([]*schema.SCResult, error)
}

type ReceiptHandler interface {
	ProcessReceipts(transactions *map[string]data.TransactionHandler) ([]*schema.Receipt, error)
}

type LogHandler interface {
	ProcessLogs(logs *map[string]data.LogHandler) ([]*schema.Log, error)
}

type AccountsHandler interface {
	ProcessAccounts() ([]*schema.AccountBalanceUpdate, error)
}
