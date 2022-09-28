package process

import (
	"encoding/hex"

	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
)

type hyperBlockProcessor struct {
}

// NewHyperBlockProcessor will create a new instance of an hyper block processor
func NewHyperBlockProcessor() *hyperBlockProcessor {
	return &hyperBlockProcessor{}
}

// Process will process current hyper block and convert it to an avro schema block result
func (hbp *hyperBlockProcessor) Process(hyperBlock *hyperBlock.HyperBlock) (*schema.BlockResult, error) {
	block := schema.NewBlockResult()

	hash, err := hex.DecodeString(hyperBlock.Hash)
	if err != nil {
		return nil, err
	}

	block.Block.Hash = hash
	return block, nil
}
