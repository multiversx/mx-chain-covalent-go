package block

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type blockProcessor struct{}

// NewBlockProcessor creates a new instance of block processor
func NewBlockProcessor() (*blockProcessor, error) {
	return &blockProcessor{}, nil
}

// ProcessBlock converts block data to a specific structure defined by avro schema
func (bp *blockProcessor) ProcessBlock(block *data.BodyHandler) (*schema.Block, error) {
	return nil, nil
}
