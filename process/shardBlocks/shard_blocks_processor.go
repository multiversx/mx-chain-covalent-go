package shardBlocks

import (
	"encoding/hex"

	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
)

type shardBlocksProcessor struct {
	accountsProcessor process.AlteredAccountsHandler
}

// NewShardBlocksProcessor creates a new instance of a shard block processor
func NewShardBlocksProcessor(accountsProcessor process.AlteredAccountsHandler) (*shardBlocksProcessor, error) {
	if accountsProcessor == nil {
		return nil, errNilAlteredAccountsHandler
	}

	return &shardBlocksProcessor{
		accountsProcessor: accountsProcessor,
	}, nil
}

// ProcessShardBlocks converts api notarized shard blocks to avro schema shard blocks
func (sbp *shardBlocksProcessor) ProcessShardBlocks(apiBlocks []*api.NotarizedBlock) ([]*schema.ShardBlocks, error) {
	shardBlocks := make([]*schema.ShardBlocks, 0, len(apiBlocks))

	for _, apiBlock := range apiBlocks {
		if apiBlock == nil {
			continue
		}

		shardBlock, err := sbp.processShardBlock(apiBlock)
		if err != nil {
			return nil, err
		}

		shardBlocks = append(shardBlocks, shardBlock)
	}

	return shardBlocks, nil
}

func (sbp *shardBlocksProcessor) processShardBlock(apiBlock *api.NotarizedBlock) (*schema.ShardBlocks, error) {
	hash, err := hex.DecodeString(apiBlock.Hash)
	if err != nil {
		return nil, err
	}
	rootHash, err := hex.DecodeString(apiBlock.RootHash)
	if err != nil {
		return nil, err
	}
	mbHashes, err := hexStringSliceToBytesSlice(apiBlock.MiniBlockHashes)
	if err != nil {
		return nil, err
	}
	alteredAccounts, err := sbp.accountsProcessor.ProcessAccounts(apiBlock.AlteredAccounts)
	if err != nil {
		return nil, err
	}

	return &schema.ShardBlocks{
		Hash:            hash,
		Nonce:           int64(apiBlock.Nonce),
		Round:           int64(apiBlock.Round),
		Shard:           int32(apiBlock.Shard),
		RootHash:        rootHash,
		MiniBlockHashes: mbHashes,
		StateChanges:    alteredAccountsOrNil(alteredAccounts),
	}, nil
}

func hexStringSliceToBytesSlice(hashes []string) ([][]byte, error) {
	ret := make([][]byte, 0, len(hashes))
	for _, hash := range hashes {
		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			return nil, err
		}

		ret = append(ret, hashBytes)
	}

	return ret, nil
}

func alteredAccountsOrNil(alteredAccounts []*schema.AccountBalanceUpdate) []*schema.AccountBalanceUpdate {
	if len(alteredAccounts) == 0 {
		return nil
	}

	return alteredAccounts
}
