package covalent

import (
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type covalentIndexer struct {
	processor *process.DataProcessor
}

func NewCovalentDataIndexer(proc *process.DataProcessor) (*covalentIndexer, error) {
	return &covalentIndexer{
		processor: proc,
	}, nil
}

func (c *covalentIndexer) SaveBlock(args *indexer.ArgsSaveBlockData) {
	// TODO this function
	// 1. Process data from args, format it according to avro schema
	blockResult, err := c.processor.ProcessData(args)
	if err != nil {

	}
	_ = blockResult

	// 2. Prepare blockResult data to be sent in binary format
	// 3. Send blockResult binary data to covalent
}

func (c covalentIndexer) RevertIndexedBlock(header data.HeaderHandler, body data.BodyHandler) {
	panic("implement me")
}

func (c covalentIndexer) SaveRoundsInfo(roundsInfos []*indexer.RoundInfo) {
	panic("implement me")
}

func (c covalentIndexer) SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) {
	panic("implement me")
}

func (c covalentIndexer) SaveValidatorsRating(indexID string, infoRating []*indexer.ValidatorRatingInfo) {
	panic("implement me")
}

func (c covalentIndexer) SaveAccounts(blockTimestamp uint64, acc []data.UserAccountHandler) {
	panic("implement me")
}

func (c covalentIndexer) Close() error {
	panic("implement me")
}

func (c covalentIndexer) IsInterfaceNil() bool {
	panic("implement me")
}
