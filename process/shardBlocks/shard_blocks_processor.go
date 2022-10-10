package shardBlocks

import (
	"encoding/hex"

	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
)

type shardBlocksProcessor struct {
}

func NewShardBlocksProcessor() *shardBlocksProcessor {
	return &shardBlocksProcessor{}
}

func (sbp *shardBlocksProcessor) ProcessShardBlocks(apiBlocks []*api.NotarizedBlock) ([]*schemaV2.ShardBlocks, error) {
	shardBlocks := make([]*schemaV2.ShardBlocks, 0, len(apiBlocks))

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

func processShardBlock(apiBlock *api.NotarizedBlock) (*schemaV2.ShardBlocks, error) {
	hash, err := hex.DecodeString(apiBlock.Hash)
	if err != nil {
		return nil, err
	}

	return &schemaV2.ShardBlocks{
		Hash:  hash,
		Nonce: int64(apiBlock.Nonce),
		Round: int64(apiBlock.Round),
		Shard: int32(apiBlock.Shard),
	}, nil
}
