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

//Logs         []*Log
//StateChanges []*AccountBalanceUpdate
//
//type LogsAndEventsHandler interface {
//	ProcessLogsAndEvents() ([]schema.Log, error)
//}
//
//type AccountsHandler interface {
//	ProcessAccounts() ([]schema., error)
//}
