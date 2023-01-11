package processMocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
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
