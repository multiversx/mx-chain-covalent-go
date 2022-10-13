package process

import (
	"encoding/hex"

	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
)

// HyperBlockProcessorArgs holds all input dependencies required
// by hyper block processor in order to create a new hyper block processor
type HyperBlockProcessorArgs struct {
	TransactionHandler     TransactionHandler
	ShardBlockHandler      ShardBlocksHandler
	EpochStartInfoHandler  EpochStartInfoHandler
	AlteredAccountsHandler AlteredAccountsHandler
}

type hyperBlockProcessor struct {
	transactionProcessor     TransactionHandler
	shardBlocksProcessor     ShardBlocksHandler
	epochStartInfoProcessor  EpochStartInfoHandler
	alteredAccountsProcessor AlteredAccountsHandler
}

// NewHyperBlockProcessor will create a new instance of an hyper block processor
func NewHyperBlockProcessor(args *HyperBlockProcessorArgs) (*hyperBlockProcessor, error) {
	if args.TransactionHandler == nil {
		return nil, errNilTransactionHandler
	}
	if args.ShardBlockHandler == nil {
		return nil, errNilShardBlocksHandler
	}
	if args.EpochStartInfoHandler == nil {
		return nil, errNilEpochStartInfoHandler
	}
	if args.AlteredAccountsHandler == nil {
		return nil, errNilAlteredAccountsHandler
	}

	return &hyperBlockProcessor{
		transactionProcessor:     args.TransactionHandler,
		shardBlocksProcessor:     args.ShardBlockHandler,
		epochStartInfoProcessor:  args.EpochStartInfoHandler,
		alteredAccountsProcessor: args.AlteredAccountsHandler,
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
		EpochStartInfo:         epochStartInfoOrNil(epochStartInfo),
		ShardBlocks:            shardBlocksOrNil(shardBlocks),
		Transactions:           txsOrNil(txs),
		StateChanges:           nil,
		Status:                 hyperBlock.Status,
	}, nil
}

func txsOrNil(txs []*schemaV2.Transaction) []*schemaV2.Transaction {
	if len(txs) == 0 {
		return nil
	}
	return txs
}

func shardBlocksOrNil(shardBlocks []*schemaV2.ShardBlocks) []*schemaV2.ShardBlocks {
	if len(shardBlocks) == 0 {
		return nil
	}
	return shardBlocks
}

func epochStartInfoOrNil(epochStartInfo *schemaV2.EpochStartInfo) *schemaV2.EpochStartInfo {
	if emptyEpochStartInfo(epochStartInfo) {
		return nil
	}

	return epochStartInfo
}

func emptyEpochStartInfo(epochStartInfo *schemaV2.EpochStartInfo) bool {
	if epochStartInfo == nil {
		return true
	}

	return len(epochStartInfo.TotalSupply) == 0 &&
		len(epochStartInfo.TotalToDistribute) == 0 &&
		len(epochStartInfo.TotalNewlyMinted) == 0 &&
		len(epochStartInfo.RewardsPerBlock) == 0 &&
		len(epochStartInfo.RewardsForProtocolSustainability) == 0 &&
		len(epochStartInfo.NodePrice) == 0 &&
		len(epochStartInfo.PrevEpochStartHash) == 0 &&
		epochStartInfo.PrevEpochStartRound == 0
}
