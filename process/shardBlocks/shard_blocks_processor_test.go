package shardBlocks

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/processMocks"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
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
		Hash:            testscommon.GenerateRandHexString(),
		Nonce:           4,
		Round:           5,
		Shard:           2,
		RootHash:        testscommon.GenerateRandHexString(),
		MiniBlockHashes: []string{testscommon.GenerateRandHexString(), testscommon.GenerateRandHexString()},
		AlteredAccounts: []*outport.AlteredAccount{{Address: "erd1q"}, {Address: "erd1b"}},
	}
}

func TestNewShardBlocksProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil accounts processor, should return error", func(t *testing.T) {
		t.Parallel()

		sbp, err := NewShardBlocksProcessor(nil)
		require.Nil(t, sbp)
		require.Equal(t, errNilAlteredAccountsHandler, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		sbp, err := NewShardBlocksProcessor(&processMocks.AlteredAccountsHandlerStub{})
		require.Nil(t, err)
		require.NotNil(t, sbp)
	})
}

func TestProcessShardBlocks(t *testing.T) {
	t.Parallel()

	accountsProcessor := &processMocks.AlteredAccountsHandlerStub{
		ProcessAccountsCalled: func(apiAlteredAccounts []*outport.AlteredAccount) ([]*schema.AccountBalanceUpdate, error) {
			ret := make([]*schema.AccountBalanceUpdate, 0, len(apiAlteredAccounts))

			for _, apiAlteredAcc := range apiAlteredAccounts {
				ret = append(ret, &schema.AccountBalanceUpdate{
					Address: []byte(apiAlteredAcc.Address),
				})
			}

			return ret, nil
		},
	}
	sp, _ := NewShardBlocksProcessor(accountsProcessor)

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		apiBlocks := generateShardBlocks(10)
		shardBlocks, err := sp.ProcessShardBlocks(apiBlocks)
		require.Nil(t, err)
		requireShardBlocksProcessedSuccessfully(t, apiBlocks, shardBlocks, accountsProcessor)
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
		requireShardBlocksProcessedSuccessfully(t, apiBlocks[1:], shardBlocks, accountsProcessor)
	})

	t.Run("invalid rootHash, should return error", func(t *testing.T) {
		t.Parallel()

		apiBlocks := generateShardBlocks(10)
		apiBlocks[4].RootHash = "invalidRootHash"
		shardBlocks, err := sp.ProcessShardBlocks(apiBlocks)
		require.Nil(t, shardBlocks)
		require.NotNil(t, err)
	})

	t.Run("invalid mini block hash, should return error", func(t *testing.T) {
		t.Parallel()

		apiBlocks := generateShardBlocks(10)
		apiBlocks[4].MiniBlockHashes[0] = "invalidMbHash"
		shardBlocks, err := sp.ProcessShardBlocks(apiBlocks)
		require.Nil(t, shardBlocks)
		require.NotNil(t, err)
	})

	t.Run("empty altered accounts, should not fill it", func(t *testing.T) {
		t.Parallel()

		alteredAccountsProcessor := &processMocks.AlteredAccountsHandlerStub{
			ProcessAccountsCalled: func(apiAlteredAccounts []*outport.AlteredAccount) ([]*schema.AccountBalanceUpdate, error) {
				return []*schema.AccountBalanceUpdate{}, nil
			},
		}
		shardBlockProcessor, _ := NewShardBlocksProcessor(alteredAccountsProcessor)

		apiBlocks := generateShardBlocks(1)
		processedHyperBlock, err := shardBlockProcessor.ProcessShardBlocks(apiBlocks)
		require.Nil(t, err)
		require.Nil(t, processedHyperBlock[0].StateChanges)
	})

	t.Run("invalid altered accounts, should return error", func(t *testing.T) {
		t.Parallel()

		errProcessAlteredAccounts := errors.New("error processing altered accounts")
		alteredAccountsProcessor := &processMocks.AlteredAccountsHandlerStub{
			ProcessAccountsCalled: func(apiAlteredAccounts []*outport.AlteredAccount) ([]*schema.AccountBalanceUpdate, error) {
				return nil, errProcessAlteredAccounts
			},
		}
		shardBlockProcessor, _ := NewShardBlocksProcessor(alteredAccountsProcessor)

		apiBlocks := generateShardBlocks(10)
		processedHyperBlock, err := shardBlockProcessor.ProcessShardBlocks(apiBlocks)
		require.Nil(t, processedHyperBlock)
		require.Equal(t, errProcessAlteredAccounts, err)
	})

}

func requireShardBlocksProcessedSuccessfully(
	t *testing.T,
	apiBlocks []*api.NotarizedBlock,
	processedBlocks []*schema.ShardBlocks,
	alteredAccountsHandler process.AlteredAccountsHandler,
) {
	require.Equal(t, len(apiBlocks), len(processedBlocks))

	for idx, apiBlock := range apiBlocks {
		processedBlock := processedBlocks[idx]
		requireShardBlockProcessedSuccessfully(t, apiBlock, processedBlock, alteredAccountsHandler)
	}
}

func requireShardBlockProcessedSuccessfully(
	t *testing.T,
	apiBlock *api.NotarizedBlock,
	processedBlock *schema.ShardBlocks,
	alteredAccountsHandler process.AlteredAccountsHandler,
) {
	hash, err := hex.DecodeString(apiBlock.Hash)
	require.Nil(t, err)
	rootHash, err := hex.DecodeString(apiBlock.RootHash)
	require.Nil(t, err)

	mbHashes := make([][]byte, 0)
	for _, mbHash := range apiBlock.MiniBlockHashes {
		decodedMbHash, err := hex.DecodeString(mbHash)
		require.Nil(t, err)

		mbHashes = append(mbHashes, decodedMbHash)
	}

	alteredAccounts, err := alteredAccountsHandler.ProcessAccounts(apiBlock.AlteredAccounts)
	require.Nil(t, err)

	require.Equal(t, &schema.ShardBlocks{
		Hash:            hash,
		Nonce:           int64(apiBlock.Nonce),
		Round:           int64(apiBlock.Round),
		Shard:           int32(apiBlock.Shard),
		RootHash:        rootHash,
		MiniBlockHashes: mbHashes,
		StateChanges:    alteredAccounts,
	}, processedBlock)
}
