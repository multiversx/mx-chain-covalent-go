package covalent

import (
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type covalentIndexer struct {
	processor DataHandler
}

// NewCovalentDataIndexer creates a new instance of covalent data indexer, which implements Driver interface and
// converts protocol input data to covalent required data
func NewCovalentDataIndexer(processor DataHandler) (Driver, error) {
	return &covalentIndexer{
		processor: processor,
	}, nil
}

// SaveBlock saves the block info and converts it in order to be sent to covalent
func (c *covalentIndexer) SaveBlock(args *indexer.ArgsSaveBlockData) {
	// TODO this function in future PRs
	// 1. Process data from args, format it according to avro schema
	blockResult, err := c.processor.ProcessData(args)
	if err != nil {

	}
	_ = blockResult

	// 2. Prepare blockResult data to be sent in binary format
	// 3. Send blockResult binary data to covalent
}

// RevertIndexedBlock DUMMY
func (c *covalentIndexer) RevertIndexedBlock(header data.HeaderHandler, body data.BodyHandler) {}

// SaveRoundsInfo DUMMY
func (c *covalentIndexer) SaveRoundsInfo(roundsInfos []*indexer.RoundInfo) {}

// SaveValidatorsPubKeys DUMMY
func (c *covalentIndexer) SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) {
}

// SaveValidatorsRating DUMMY
func (c *covalentIndexer) SaveValidatorsRating(indexID string, infoRating []*indexer.ValidatorRatingInfo) {
}

// SaveAccounts DUMMY
func (c *covalentIndexer) SaveAccounts(blockTimestamp uint64, acc []data.UserAccountHandler) {}

// Close DUMMY
func (c *covalentIndexer) Close() error {
	return nil
}

// IsInterfaceNil DUMMY
func (c *covalentIndexer) IsInterfaceNil() bool {
	return false
}
