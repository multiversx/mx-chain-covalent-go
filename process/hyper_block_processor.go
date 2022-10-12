package process

import (
	"encoding/hex"

	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
)

type hyperBlockProcessor struct {
	transactionProcessor    TransactionHandler
	shardBlocksProcessor    ShardBlocksHandler
	epochStartInfoProcessor EpochStartInfoHandler
}

// NewHyperBlockProcessor will create a new instance of an hyper block processor
func NewHyperBlockProcessor(
	transactionHandler TransactionHandler,
	shardBlockHandler ShardBlocksHandler,
	epochStartInfoProcessor EpochStartInfoHandler,
) (*hyperBlockProcessor, error) {
	if transactionHandler == nil {
		return nil, errNilTransactionHandler
	}
	if shardBlockHandler == nil {
		return nil, errNilShardBlocksHandler
	}
	if epochStartInfoProcessor == nil {
		return nil, errNilEpochStartInfoHandler
	}

	return &hyperBlockProcessor{
		transactionProcessor:    transactionHandler,
		shardBlocksProcessor:    shardBlockHandler,
		epochStartInfoProcessor: epochStartInfoProcessor,
	}, nil
}

// Process will process current hyper block and convert it to an avro schema block result
func (hbp *hyperBlockProcessor) Process(hyperBlock *hyperBlock.HyperBlock) (*schemaV2.HyperBlock, error) {
	hash, err := hex.DecodeString(hyperBlock.Hash)
	if err != nil {
		return nil, err
	}
	prevBlockHash, err := hex.DecodeString(hyperBlock.PrevBlockHash)
	if err != nil {
		return nil, err
	}
	stateRootHash, err := hex.DecodeString(hyperBlock.StateRootHash)
	if err != nil {
		return nil, err
	}
	accumulatedFees, err := utility.GetBigIntBytesFromStr(hyperBlock.AccumulatedFees)
	if err != nil {
		return nil, err
	}
	developerFees, err := utility.GetBigIntBytesFromStr(hyperBlock.DeveloperFees)
	if err != nil {
		return nil, err
	}
	accumulatedFeesInEpoch, err := utility.GetBigIntBytesFromStr(hyperBlock.AccumulatedFeesInEpoch)
	if err != nil {
		return nil, err
	}
	developerFeesInEpoch, err := utility.GetBigIntBytesFromStr(hyperBlock.DeveloperFeesInEpoch)
	if err != nil {
		return nil, err
	}
	txs, err := hbp.transactionProcessor.ProcessTransactions(hyperBlock.Transactions)
	if err != nil {
		return nil, err
	}
	shardBlocks, err := hbp.shardBlocksProcessor.ProcessShardBlocks(hyperBlock.ShardBlocks)
	if err != nil {
		return nil, err
	}
	epochStartInfo, err := hbp.epochStartInfoProcessor.ProcessEpochStartInfo(hyperBlock.EpochStartInfo)
	if err != nil {
		return nil, err
	}

	return &schemaV2.HyperBlock{
		Hash:                   hash,
		PrevBlockHash:          prevBlockHash,
		StateRootHash:          stateRootHash,
		Nonce:                  int64(hyperBlock.Nonce),
		Round:                  int64(hyperBlock.Round),
		Epoch:                  int32(hyperBlock.Epoch),
		NumTxs:                 int32(hyperBlock.NumTxs),
		AccumulatedFees:        accumulatedFees,
		DeveloperFees:          developerFees,
		AccumulatedFeesInEpoch: accumulatedFeesInEpoch,
		DeveloperFeesInEpoch:   developerFeesInEpoch,
		Timestamp:              int64(hyperBlock.Timestamp),
		EpochStartInfo:         epochStartInfo,
		ShardBlocks:            shardBlocks,
		Transactions:           txs,
		StateChanges:           nil,
		Status:                 hyperBlock.Status,
	}, nil
}
