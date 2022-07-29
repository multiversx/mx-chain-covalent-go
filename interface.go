package covalent

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
)

// DataHandler defines the behaviour of a component that is able to process the data of a block
type DataHandler interface {
	ProcessData(args *outport.ArgsSaveBlockData) (*schema.BlockResult, error)
}

// Driver defines the behaviour of an outport driver that should handle the saving or processing of data resulted from node
type Driver interface {
	SaveBlock(args *outport.ArgsSaveBlockData) error
	RevertIndexedBlock(header data.HeaderHandler, body data.BodyHandler) error
	SaveRoundsInfo(roundsInfos []*outport.RoundInfo) error
	SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) error
	SaveValidatorsRating(indexID string, infoRating []*outport.ValidatorRatingInfo) error
	SaveAccounts(blockTimestamp uint64, acc map[string]*outport.AlteredAccount) error
	FinalizedBlock(headerHash []byte) error
	Close() error
	IsInterfaceNil() bool
}
