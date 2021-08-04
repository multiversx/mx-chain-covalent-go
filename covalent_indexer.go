package covalent

import (
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type covalentIndexer struct {
	processor DataHandler
}

func NewCovalentDataIndexer(processor DataHandler) (Driver, error) {
	return &covalentIndexer{
		processor: processor,
	}, nil
}

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

func (c covalentIndexer) RevertIndexedBlock(header data.HeaderHandler, body data.BodyHandler) {}

func (c covalentIndexer) SaveRoundsInfo(roundsInfos []*indexer.RoundInfo) {}

func (c covalentIndexer) SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) {}

func (c covalentIndexer) SaveValidatorsRating(indexID string, infoRating []*indexer.ValidatorRatingInfo) {
}

func (c covalentIndexer) SaveAccounts(blockTimestamp uint64, acc []data.UserAccountHandler) {}

func (c covalentIndexer) Close() error {
	return nil
}

func (c covalentIndexer) IsInterfaceNil() bool {
	return false
}
