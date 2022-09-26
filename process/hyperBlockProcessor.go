package process

import (
	"encoding/hex"

	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
)

type hyperBlockProcessor struct {
}

func NewHyperBlockProcessor() *hyperBlockProcessor {
	return &hyperBlockProcessor{}
}

func (hbp *hyperBlockProcessor) Process(hyperBlock *hyperBlock.HyperBlock) (*schema.BlockResult, error) {
	block := schema.NewBlockResult()

	hash, err := hex.DecodeString(hyperBlock.Hash)
	if err != nil {
		return nil, err
	}

	block.Block.Hash = hash
	return block, nil
}
