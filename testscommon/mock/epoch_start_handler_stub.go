package mock

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
)

type EpochStartInfoHandlerStub struct {
	ProcessEpochStartInfoCalled func(apiEpochInfo *api.EpochStartInfo) (*schemaV2.EpochStartInfo, error)
}

func (esi *EpochStartInfoHandlerStub) ProcessEpochStartInfo(apiEpochInfo *api.EpochStartInfo) (*schemaV2.EpochStartInfo, error) {
	if esi.ProcessEpochStartInfoCalled != nil {
		return esi.ProcessEpochStartInfoCalled(apiEpochInfo)
	}

	return nil, nil
}
