package shardBlocks

import (
	"encoding/hex"

	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
)

type shardBlocksProcessor struct {
}

// NewShardBlocksProcessor creates a new instance of a shard block processor
func NewShardBlocksProcessor() *shardBlocksProcessor {
	return &shardBlocksProcessor{}
}

// ProcessShardBlocks converts api notarized shard blocks to avro schema shard blocks
func (sbp *shardBlocksProcessor) ProcessShardBlocks(apiBlocks []*api.NotarizedBlock) ([]*schema.ShardBlocks, error) {
	shardBlocks := make([]*schema.ShardBlocks, 0, len(apiBlocks))

	for _, apiBlock := range apiBlocks {
		if apiBlock == nil {
			continue
		}

		shardBlock, err := processShardBlock(apiBlock)
		if err != nil {
			return nil, err
		}

		shardBlocks = append(shardBlocks, shardBlock)
	}

	return shardBlocks, nil
}

func processShardBlock(apiBlock *api.NotarizedBlock) (*schema.ShardBlocks, error) {
	hash, err := hex.DecodeString(apiBlock.Hash)
	if err != nil {
		return nil, err
	}

	return &schema.ShardBlocks{
		Hash:  hash,
		Nonce: int64(apiBlock.Nonce),
		Round: int64(apiBlock.Round),
		Shard: int32(apiBlock.Shard),
	}, nil
}
