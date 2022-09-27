package covalent

import (
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/elodina/go-avro"
)

type DataHandler interface {
	ProcessData(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error)
}

// HyperBlockProcessor shall handle hyper block processing into avro schema blocks
type HyperBlockProcessor interface {
	Process(hyperBlock *hyperBlock.HyperBlock) (*schema.BlockResult, error)
}

type Driver interface {
	SaveBlock(args *indexer.ArgsSaveBlockData) error
	RevertIndexedBlock(header data.HeaderHandler, body data.BodyHandler) error
	SaveRoundsInfo(roundsInfos []*indexer.RoundInfo) error
	SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) error
	SaveValidatorsRating(indexID string, infoRating []*indexer.ValidatorRatingInfo) error
	SaveAccounts(blockTimestamp uint64, acc []data.UserAccountHandler) error
	FinalizedBlock(headerHash []byte) error
	Close() error
	IsInterfaceNil() bool
}

type AccountsAdapter interface {
	LoadAccount(address []byte) (vmcommon.AccountHandler, error)
	IsInterfaceNil() bool
}

type HttpServer interface {
	ListenAndServe() error
	Close() error
}

// AvroMarshaller defines what an avro marshaller should do
type AvroMarshaller interface {
	Encode(record avro.AvroRecord) ([]byte, error)
	Decode(record avro.AvroRecord, buffer []byte) error
}
