package mock

import (
	"github.com/multiversx/mx-chain-covalent-go/hyperBlock"
	"github.com/multiversx/mx-chain-covalent-go/schema"
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
