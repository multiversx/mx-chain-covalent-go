package process

import (
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/processMocks"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/stretchr/testify/require"
)

func createHyperBlockProcessorArgs() *HyperBlockProcessorArgs {
	return &HyperBlockProcessorArgs{
		TransactionHandler:     &processMocks.TransactionHandlerStub{},
		ShardBlockHandler:      &processMocks.ShardBlocksHandlerStub{},
		EpochStartInfoHandler:  &processMocks.EpochStartInfoHandlerStub{},
		AlteredAccountsHandler: &processMocks.AlteredAccountsHandlerStub{},
	}
}

func TestNewHyperBlockProcessor(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createHyperBlockProcessorArgs()
		hbp, err := NewHyperBlockProcessor(args)
		require.Nil(t, err)
		require.NotNil(t, hbp)
	})

	t.Run("nil transaction processor, should return error", func(t *testing.T) {
		t.Parallel()

		args := createHyperBlockProcessorArgs()
		args.TransactionHandler = nil

		hbp, err := NewHyperBlockProcessor(args)
		require.Nil(t, hbp)
		require.Equal(t, errNilTransactionHandler, err)
	})

	t.Run("nil shard block processor, should return error", func(t *testing.T) {
		t.Parallel()

		args := createHyperBlockProcessorArgs()
		args.ShardBlockHandler = nil

		hbp, err := NewHyperBlockProcessor(args)
		require.Nil(t, hbp)
		require.Equal(t, errNilShardBlocksHandler, err)
	})

	t.Run("nil epoch start info processor, should return error", func(t *testing.T) {
		t.Parallel()

		args := createHyperBlockProcessorArgs()
		args.EpochStartInfoHandler = nil

		hbp, err := NewHyperBlockProcessor(args)
		require.Nil(t, hbp)
		require.Equal(t, errNilEpochStartInfoHandler, err)
	})

	t.Run("nil altered accounts processor, should return error", func(t *testing.T) {
		t.Parallel()

		args := createHyperBlockProcessorArgs()
		args.AlteredAccountsHandler = nil

		hbp, err := NewHyperBlockProcessor(args)
		require.Nil(t, hbp)
		require.Equal(t, errNilAlteredAccountsHandler, err)
	})
}

func TestHyperBlockProcessor_Process(t *testing.T) {
	t.Parallel()

	apiTxs := []*transaction.ApiTransactionResult{{Hash: "hash1"}}
	alteredAcc := &outport.AlteredAccount{Balance: "100"}
	shardBlocks := []*api.NotarizedBlock{{Hash: "hash2", AlteredAccounts: []*outport.AlteredAccount{alteredAcc}}}
	epochStartInfo := &api.EpochStartInfo{NodePrice: "100"}

	processedTxs := []*schema.Transaction{{Hash: []byte(apiTxs[0].Hash)}}
	processedAlteredAcc := &schema.AccountBalanceUpdate{Balance: big.NewInt(100).Bytes()}
	processedShardBlocks := []*schema.ShardBlocks{{Hash: []byte(shardBlocks[0].Hash)}}
	processedEpochStartInfo := &schema.EpochStartInfo{NodePrice: big.NewInt(100).Bytes()}

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

	expectedProcessedHyperBlock := &schema.HyperBlock{
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
		StateChanges:           []*schema.AccountBalanceUpdate{processedAlteredAcc},
		Status:                 "status",
	}

	txProcessor := &processMocks.TransactionHandlerStub{
		ProcessTransactionsCalled: func(apiTransactions []*transaction.ApiTransactionResult) ([]*schema.Transaction, error) {
			require.Equal(t, apiTxs, apiTransactions)
			return processedTxs, nil
		},
	}
	shardBlocksProcessor := &processMocks.ShardBlocksHandlerStub{
		ProcessShardBlocksCalled: func(apiBlocks []*api.NotarizedBlock) ([]*schema.ShardBlocks, error) {
			require.Equal(t, shardBlocks, apiBlocks)
			return processedShardBlocks, nil
		},
	}
	epochStartInfoProcessor := &processMocks.EpochStartInfoHandlerStub{
		ProcessEpochStartInfoCalled: func(apiEpochInfo *api.EpochStartInfo) (*schema.EpochStartInfo, error) {
			require.Equal(t, epochStartInfo, apiEpochInfo)
			return processedEpochStartInfo, nil
		},
	}
	alteredAccProcessor := &processMocks.AlteredAccountsHandlerStub{
		ProcessAccountsCalled: func(apiNotarizedBlocks []*api.NotarizedBlock) ([]*schema.AccountBalanceUpdate, error) {
			require.Equal(t, shardBlocks, apiNotarizedBlocks)
			return []*schema.AccountBalanceUpdate{processedAlteredAcc}, nil
		},
	}

	t.Run("should work", func(t *testing.T) {
		t.Parallel()
		args := &HyperBlockProcessorArgs{
			TransactionHandler:     txProcessor,
			ShardBlockHandler:      shardBlocksProcessor,
			EpochStartInfoHandler:  epochStartInfoProcessor,
			AlteredAccountsHandler: alteredAccProcessor,
		}
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(apiHyperBLock)
		require.Nil(t, err)
		require.Equal(t, expectedProcessedHyperBlock, processedHyperBlock)
	})

	t.Run("invalid hash, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		apiHyperBLockCopy.Hash = "hash"
		args := createHyperBlockProcessorArgs()
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.NotNil(t, err)
	})

	t.Run("invalid prev block hash, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		apiHyperBLockCopy.PrevBlockHash = "prev block hash"
		args := createHyperBlockProcessorArgs()
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.NotNil(t, err)
	})

	t.Run("invalid state root hash, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		apiHyperBLockCopy.StateRootHash = "state root hash"
		args := createHyperBlockProcessorArgs()
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.NotNil(t, err)
	})

	t.Run("invalid accumulated fees, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		apiHyperBLockCopy.AccumulatedFees = "accumulated fees"
		args := createHyperBlockProcessorArgs()
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "accumulated fees"))
	})

	t.Run("invalid developer fees, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		apiHyperBLockCopy.DeveloperFees = "developer fees"
		args := createHyperBlockProcessorArgs()
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "developer fees"))
	})

	t.Run("invalid accumulated fees in epoch, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		apiHyperBLockCopy.AccumulatedFeesInEpoch = "accumulated fees in epoch"
		args := createHyperBlockProcessorArgs()
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "accumulated fees in epoch"))
	})

	t.Run("invalid developer fees in epoch, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		apiHyperBLockCopy.DeveloperFeesInEpoch = "developer fees in epoch"
		args := createHyperBlockProcessorArgs()
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "developer fees in epoch"))
	})

	t.Run("empty processed txs, should fill txs field with nil", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		args := &HyperBlockProcessorArgs{
			TransactionHandler: &processMocks.TransactionHandlerStub{
				ProcessTransactionsCalled: func(apiTransactions []*transaction.ApiTransactionResult) ([]*schema.Transaction, error) {
					return []*schema.Transaction{}, nil
				},
			},
			ShardBlockHandler:      shardBlocksProcessor,
			EpochStartInfoHandler:  epochStartInfoProcessor,
			AlteredAccountsHandler: alteredAccProcessor,
		}
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, err)

		expectedProcessedHyperBlockCopy := *expectedProcessedHyperBlock
		expectedProcessedHyperBlockCopy.Transactions = nil
		require.Equal(t, &expectedProcessedHyperBlockCopy, processedHyperBlock)
	})

	t.Run("invalid txs, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		args := createHyperBlockProcessorArgs()
		errProcessTransactions := errors.New("error processing transactions")
		args.TransactionHandler = &processMocks.TransactionHandlerStub{
			ProcessTransactionsCalled: func(apiTransactions []*transaction.ApiTransactionResult) ([]*schema.Transaction, error) {
				return nil, errProcessTransactions
			},
		}
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.Equal(t, errProcessTransactions, err)
	})

	t.Run("empty shard blocks, should fill shard blocks field with nil", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		args := &HyperBlockProcessorArgs{
			TransactionHandler: txProcessor,
			ShardBlockHandler: &processMocks.ShardBlocksHandlerStub{
				ProcessShardBlocksCalled: func(apiBlocks []*api.NotarizedBlock) ([]*schema.ShardBlocks, error) {
					return []*schema.ShardBlocks{}, nil
				},
			},
			EpochStartInfoHandler:  epochStartInfoProcessor,
			AlteredAccountsHandler: alteredAccProcessor,
		}
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, err)

		expectedProcessedHyperBlockCopy := *expectedProcessedHyperBlock
		expectedProcessedHyperBlockCopy.ShardBlocks = nil
		require.Equal(t, &expectedProcessedHyperBlockCopy, processedHyperBlock)
	})

	t.Run("invalid shard blocks, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		args := createHyperBlockProcessorArgs()
		errProcessShardBlocks := errors.New("error processing shard blocks")
		args.ShardBlockHandler = &processMocks.ShardBlocksHandlerStub{
			ProcessShardBlocksCalled: func(apiBlocks []*api.NotarizedBlock) ([]*schema.ShardBlocks, error) {
				return nil, errProcessShardBlocks
			},
		}
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.Equal(t, errProcessShardBlocks, err)
	})

	t.Run("nil epoch start info, should fill epoch start info with nil", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		args := &HyperBlockProcessorArgs{
			TransactionHandler: txProcessor,
			ShardBlockHandler:  shardBlocksProcessor,
			EpochStartInfoHandler: &processMocks.EpochStartInfoHandlerStub{
				ProcessEpochStartInfoCalled: func(apiEpochInfo *api.EpochStartInfo) (*schema.EpochStartInfo, error) {
					return nil, nil
				},
			},
			AlteredAccountsHandler: alteredAccProcessor,
		}
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, err)

		expectedProcessedHyperBlockCopy := *expectedProcessedHyperBlock
		expectedProcessedHyperBlockCopy.EpochStartInfo = nil
		require.Equal(t, &expectedProcessedHyperBlockCopy, processedHyperBlock)
	})

	t.Run("empty epoch start info, should fill epoch start info with nil", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		args := &HyperBlockProcessorArgs{
			TransactionHandler: txProcessor,
			ShardBlockHandler:  shardBlocksProcessor,
			EpochStartInfoHandler: &processMocks.EpochStartInfoHandlerStub{
				ProcessEpochStartInfoCalled: func(apiEpochInfo *api.EpochStartInfo) (*schema.EpochStartInfo, error) {
					return schema.NewEpochStartInfo(), nil
				},
			},
			AlteredAccountsHandler: alteredAccProcessor,
		}
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, err)

		expectedProcessedHyperBlockCopy := *expectedProcessedHyperBlock
		expectedProcessedHyperBlockCopy.EpochStartInfo = nil
		require.Equal(t, &expectedProcessedHyperBlockCopy, processedHyperBlock)
	})

	t.Run("invalid epoch start info, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		args := createHyperBlockProcessorArgs()
		errProcessEpochStartInfo := errors.New("error processing epoch start info")
		args.EpochStartInfoHandler = &processMocks.EpochStartInfoHandlerStub{
			ProcessEpochStartInfoCalled: func(apiEpochInfo *api.EpochStartInfo) (*schema.EpochStartInfo, error) {
				return nil, errProcessEpochStartInfo
			},
		}
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.Equal(t, errProcessEpochStartInfo, err)
	})

	t.Run("empty altered accounts, should not fill it", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		args := &HyperBlockProcessorArgs{
			TransactionHandler:    txProcessor,
			ShardBlockHandler:     shardBlocksProcessor,
			EpochStartInfoHandler: epochStartInfoProcessor,
			AlteredAccountsHandler: &processMocks.AlteredAccountsHandlerStub{
				ProcessAccountsCalled: func(apiNotarizedBlocks []*api.NotarizedBlock) ([]*schema.AccountBalanceUpdate, error) {
					return []*schema.AccountBalanceUpdate{}, nil
				},
			},
		}
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, err)

		expectedProcessedHyperBlockCopy := *expectedProcessedHyperBlock
		expectedProcessedHyperBlockCopy.StateChanges = nil
		require.Equal(t, &expectedProcessedHyperBlockCopy, processedHyperBlock)
	})

	t.Run("invalid altered accounts, should return error", func(t *testing.T) {
		t.Parallel()

		apiHyperBLockCopy := *apiHyperBLock
		args := createHyperBlockProcessorArgs()
		errProcessAlteredAccounts := errors.New("error processing altered accounts")
		args.AlteredAccountsHandler = &processMocks.AlteredAccountsHandlerStub{
			ProcessAccountsCalled: func(apiNotarizedBlocks []*api.NotarizedBlock) ([]*schema.AccountBalanceUpdate, error) {
				return nil, errProcessAlteredAccounts
			},
		}
		hbp, _ := NewHyperBlockProcessor(args)

		processedHyperBlock, err := hbp.Process(&apiHyperBLockCopy)
		require.Nil(t, processedHyperBlock)
		require.Equal(t, errProcessAlteredAccounts, err)
	})
}
