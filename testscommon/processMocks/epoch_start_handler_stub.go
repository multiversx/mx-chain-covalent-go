package processMocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
)

// EpochStartInfoHandlerStub -
type EpochStartInfoHandlerStub struct {
	ProcessEpochStartInfoCalled func(apiEpochInfo *api.EpochStartInfo) (*schemaV2.EpochStartInfo, error)
}

// ProcessEpochStartInfo -
func (esi *EpochStartInfoHandlerStub) ProcessEpochStartInfo(apiEpochInfo *api.EpochStartInfo) (*schemaV2.EpochStartInfo, error) {
	if esi.ProcessEpochStartInfoCalled != nil {
		return esi.ProcessEpochStartInfoCalled(apiEpochInfo)
	}

	return nil, nil
}
