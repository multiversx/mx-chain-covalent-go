package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type DataProcessor struct {
	blockHandler BlockHandler
}

func NewDataProcessor(blockHandler BlockHandler) (*DataProcessor, error) {
	return &DataProcessor{
		blockHandler: blockHandler,
	}, nil
}

func (d *DataProcessor) ProcessData(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
	block, err := d.blockHandler.ProcessBlock(args.Body)
	if err != nil {
		return nil, err
	}

	return &schema.BlockResult{
		Block: block,
	}, nil
}
