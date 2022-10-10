package mock

import (
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
)

// HyperBlockProcessorStub -
type HyperBlockProcessorStub struct {
	ProcessCalled func(hyperBlock *hyperBlock.HyperBlock) (*schemaV2.HyperBlock, error)
}

// Process -
func (hbp *HyperBlockProcessorStub) Process(hyperBlock *hyperBlock.HyperBlock) (*schemaV2.HyperBlock, error) {
	if hbp.ProcessCalled != nil {
		return hbp.ProcessCalled(hyperBlock)
	}

	return nil, nil
}
