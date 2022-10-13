package shardBlocks

import (
	"encoding/hex"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/stretchr/testify/require"
)

func generateShardBlocks(n int) []*api.NotarizedBlock {
	ret := make([]*api.NotarizedBlock, n)

	for i := 0; i < n; i++ {
		ret[i] = generateShardBlock()
	}

	return ret
}

func generateShardBlock() *api.NotarizedBlock {
	return &api.NotarizedBlock{
		Hash:  testscommon.GenerateRandHexString(),
		Nonce: 4,
		Round: 5,
		Shard: 2,
	}
}

func TestProcessShardBlocks(t *testing.T) {
	t.Parallel()

	sp := NewShardBlocksProcessor()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		apiBlocks := generateShardBlocks(10)
		shardBlocks, err := sp.ProcessShardBlocks(apiBlocks)
		require.Nil(t, err)
		requireShardBlocksProcessedSuccessfully(t, apiBlocks, shardBlocks)
	})

	t.Run("invalid block hash", func(t *testing.T) {
		t.Parallel()

		apiBlocks := generateShardBlocks(10)
		apiBlocks[4].Hash = "invalid"
		shardBlocks, err := sp.ProcessShardBlocks(apiBlocks)
		require.Nil(t, shardBlocks)
		require.NotNil(t, err)
	})

	t.Run("nil api block, should skip it", func(t *testing.T) {
		t.Parallel()

		apiBlocks := generateShardBlocks(10)
		apiBlocks[0] = nil
		shardBlocks, err := sp.ProcessShardBlocks(apiBlocks)
		require.Nil(t, err)
		requireShardBlocksProcessedSuccessfully(t, apiBlocks[1:], shardBlocks)
	})

}

func requireShardBlocksProcessedSuccessfully(
	t *testing.T,
	apiBlocks []*api.NotarizedBlock,
	processedBlocks []*schema.ShardBlocks,
) {
	require.Equal(t, len(apiBlocks), len(processedBlocks))

	for idx, apiBlock := range apiBlocks {
		processedBlock := processedBlocks[idx]
		requireShardBlockProcessedSuccessfully(t, apiBlock, processedBlock)
	}
}

func requireShardBlockProcessedSuccessfully(
	t *testing.T,
	apiBlock *api.NotarizedBlock,
	processedBlock *schema.ShardBlocks,
) {
	hash, err := hex.DecodeString(apiBlock.Hash)
	require.Nil(t, err)

	require.Equal(t, &schema.ShardBlocks{
		Hash:  hash,
		Nonce: int64(apiBlock.Nonce),
		Round: int64(apiBlock.Round),
		Shard: int32(apiBlock.Shard),
	}, processedBlock)
}
