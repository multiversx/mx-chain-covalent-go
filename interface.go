package covalent

import (
	"context"

	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type DataHandler interface {
	ProcessData(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error)
}

type Driver interface {
	SaveBlock(ctx context.Context, args *indexer.ArgsSaveBlockData)
	RevertIndexedBlock(ctx context.Context, header data.HeaderHandler, body data.BodyHandler)
	SaveRoundsInfo(ctx context.Context, roundsInfos []*indexer.RoundInfo)
	SaveValidatorsPubKeys(ctx context.Context, validatorsPubKeys map[uint32][][]byte, epoch uint32)
	SaveValidatorsRating(ctx context.Context, indexID string, infoRating []*indexer.ValidatorRatingInfo)
	SaveAccounts(ctx context.Context, blockTimestamp uint64, acc []data.UserAccountHandler)
	FinalizedBlock(ctx context.Context, headerHash []byte)
	Close() error
	IsInterfaceNil() bool
}

type AccountsAdapter interface {
	LoadAccount(address []byte) (vmcommon.AccountHandler, error)
	IsInterfaceNil() bool
}
