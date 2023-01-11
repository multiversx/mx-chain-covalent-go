package processMocks

import (
	"github.com/multiversx/mx-chain-covalent-go/schema"
	"github.com/multiversx/mx-chain-core-go/data/api"
)

// EpochStartInfoHandlerStub -
type EpochStartInfoHandlerStub struct {
	ProcessEpochStartInfoCalled func(apiEpochInfo *api.EpochStartInfo) (*schema.EpochStartInfo, error)
}

// ProcessEpochStartInfo -
func (esi *EpochStartInfoHandlerStub) ProcessEpochStartInfo(apiEpochInfo *api.EpochStartInfo) (*schema.EpochStartInfo, error) {
	if esi.ProcessEpochStartInfoCalled != nil {
		return esi.ProcessEpochStartInfoCalled(apiEpochInfo)
	}

	return nil, nil
}
