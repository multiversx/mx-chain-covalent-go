package shardBlocks

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/data/api"
)

func TestProcessShardBlocks(t *testing.T) {
	t.Parallel()

	sp := NewShardProcessor()
	apiBlocks := []*api.NotarizedBlock{
		{
			Hash:  "ff",
			Nonce: 4,
			Round: 5,
			Shard: 2,
		},
	}

	t.Run("should work", func(t *testing.T) {

	})
}

func requireShardBlockPro