package processMocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
)

// ShardBlocksHandlerStub -
type ShardBlocksHandlerStub struct {
	ProcessShardBlocksCalled func(apiBlocks []*api.NotarizedBlock) ([]*schemaV2.ShardBlocks, error)
}

// ProcessShardBlocks -
func (sbh *ShardBlocksHandlerStub) ProcessShardBlocks(apiBlocks []*api.NotarizedBlock) ([]*schemaV2.ShardBlocks, error) {
	if sbh.ProcessShardBlocksCalled != nil {
		return sbh.ProcessShardBlocksCalled(apiBlocks)
	}

	return nil, nil
}