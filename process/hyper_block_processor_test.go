package process

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/stretchr/testify/require"
)

func TestHyperBlockProcessor_Process(t *testing.T) {
	t.Parallel()

	apiTxs := []*transaction.ApiTransactionResult{{Hash: "hash1"}}
	shardBlocks := []*api.NotarizedBlock{{Hash: "hash2"}}
	epochStartInfo := &api.EpochStartInfo{NodePrice: "100"}

	apiHyperBLock := &hyperBlock.HyperBlock{
		Hash:                   "0a",
		PrevBlockHash:          "0b",
		StateRootHash:          "0c",
		Nonce:                  4,
		Round:                  5,
		Epoch:                  6,
		NumTxs:                 7,
		AccumulatedFees:        "8",
		DeveloperFees:          "9",
		AccumulatedFeesInEpoch: "10",
		DeveloperFeesInEpoch:   "11",
		Timestamp:              12,
		EpochStartInfo:         epochStartInfo,
		ShardBlocks:            shardBlocks,
		Transactions:           apiTxs,
		Status:                 "status",
	}

	t.Run("should work", func(t *testing.T) {
		processedTxs := []*schemaV2.Transaction{{Hash: []byte(apiTxs[0].Hash)}}
		processedShardBlocks := []*schemaV2.ShardBlocks{{Hash: []byte(shardBlocks[0].Hash)}}
		processedEpochStartInfo := &schemaV2.EpochStartInfo{NodePrice: big.NewInt(100).Bytes()}

		txProcessor := &mock.TransactionHandlerStub{
			ProcessTransactionsCalled: func(apiTransactions []*transaction.ApiTransactionResult) ([]*schemaV2.Transaction, error) {
				require.Equal(t, apiTxs, apiTransactions)
				return processedTxs, nil
			},
		}
		shardBlocksProcessor := &mock.ShardBlocksHandlerStub{
			ProcessShardBlocksCalled: func(apiBlocks []*api.NotarizedBlock) ([]*schemaV2.ShardBlocks, error) {
				require.Equal(t, shardBlocks, apiBlocks)
				return processedShardBlocks, nil
			},
		}
		epochStartInfoProcessor := &mock.EpochStartInfoHandlerStub{
			ProcessEpochStartInfoCalled: func(apiEpochInfo *api.EpochStartInfo) (*schemaV2.EpochStartInfo, error) {
				require.Equal(t, epochStartInfo, apiEpochInfo)
				return processedEpochStartInfo, nil
			},
		}
		hbp, _ := NewHyperBlockProcessor(txProcessor, shardBlocksProcessor, epochStartInfoProcessor)

		processedHyperBlock, err := hbp.Process(apiHyperBLock)
		require.Nil(t, err)
		require.Equal(t, &schemaV2.HyperBlock{
			Hash:                   []byte{10},
			PrevBlockHash:          []byte{11},
			StateRootHash:          []byte{12},
			Nonce:                  4,
			Round:                  5,
			Epoch:                  6,
			NumTxs:                 7,
			AccumulatedFees:        big.NewInt(8).Bytes(),
			DeveloperFees:          big.NewInt(9).Bytes(),
			AccumulatedFeesInEpoch: big.NewInt(10).Bytes(),
			DeveloperFeesInEpoch:   big.NewInt(11).Bytes(),
			Timestamp:              12,
			EpochStartInfo:         processedEpochStartInfo,
			ShardBlocks:            processedShardBlocks,
			Transactions:           processedTxs,
			StateChanges:           nil,
			Status:                 "status",
		}, processedHyperBlock)
	})
}
