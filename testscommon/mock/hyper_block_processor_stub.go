package mock

import (
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
)

// HyperBlockProcessorStub -
type HyperBlockProcessorStub struct {
	ProcessCalled func(hyperBlock *hyperBlock.HyperBlock) (*schema.HyperBlock, error)
}

// Process -
func (hbp *HyperBlockProcessorStub) Process(hyperBlock *hyperBlock.HyperBlock) (*schema.HyperBlock, error) {
	if hbp.ProcessCalled != nil {
		return hbp.ProcessCalled(hyperBlock)
	}

	return nil, nil
}
