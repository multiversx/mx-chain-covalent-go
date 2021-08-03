package covalent

import (
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type Driver interface {
	SaveBlock(args *indexer.ArgsSaveBlockData)
	RevertIndexedBlock(header data.HeaderHandler, body data.BodyHandler)
	SaveRoundsInfo(roundsInfos []*indexer.RoundInfo)
	SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32)
	SaveValidatorsRating(indexID string, infoRating []*indexer.ValidatorRatingInfo)
	SaveAccounts(blockTimestamp uint64, acc []data.UserAccountHandler)
	Close() error
	IsInterfaceNil() bool
}

//
//type TransactionHandler interface {
//	ProcessTransactions () ([]schema.Transaction, error)
//}
//
//type SCHandler interface {
//	ProcessSCs ()
//}
//
//SCResults    []*SCResult
//Receipts     []*Receipt
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
