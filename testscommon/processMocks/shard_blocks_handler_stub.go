package processMocks

import (
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-covalent-go/schema"
)

// ShardBlocksHandlerStub -
type ShardBlocksHandlerStub struct {
	ProcessShardBlocksCalled func(apiBlocks []*api.NotarizedBlock) ([]*schema.ShardBlocks, error)
}

// ProcessShardBlocks -
func (sbh *ShardBlocksHandlerStub) ProcessShardBlocks(apiBlocks []*api.NotarizedBlock) ([]*schema.ShardBlocks, error) {
	if sbh.ProcessShardBlocksCalled != nil {
		return sbh.ProcessShardBlocksCalled(apiBlocks)
	}

	return nil, nil
}
