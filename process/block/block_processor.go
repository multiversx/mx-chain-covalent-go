package block

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type blockProcessor struct{}

func NewBlockProcessor() (*blockProcessor, error) {
	return &blockProcessor{}, nil
}

func (b *blockProcessor) ProcessBlock(block *data.BodyHandler) (*schema.Block, error) {
	return nil, nil
}
