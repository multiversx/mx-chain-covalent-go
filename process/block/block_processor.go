package block

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	erdBlock "github.com/ElrondNetwork/elrond-go-core/data/block"
)

type blockProcessor struct{}

// NewBlockProcessor creates a new instance of block processor
func NewBlockProcessor() (*blockProcessor, error) {
	return &blockProcessor{}, nil
}

// ProcessBlock converts block data to a specific structure defined by avro schema
func (bp *blockProcessor) ProcessBlock(signersIndexes []uint64, hash []byte, header data.HeaderHandler, block data.BodyHandler) (*schema.Block, error) {
	body, ok := block.(*erdBlock.Body)
	if !ok {
		return nil, covalent.ErrBlockBodyAssertion
	}

	nonce := int64(header.GetNonce())
	round := int64(header.GetRound())
	epoch := int32(header.GetEpoch())
	proposer := getProposerIndex(signersIndexes)
	validators := uIntSliceToIntSlice(signersIndexes)
	pubKeysBitmap := header.GetPubKeysBitmap()
	size := int64(321321321) // TODO

	timeStamp := int64(header.GetTimeStamp())
	rootHash := header.GetRootHash()
	prevHash := header.GetPrevHash()
	shardID := int32(header.GetShardID())
	txCount := int32(header.GetTxCount())

	accumulatedFees := header.GetAccumulatedFees().Bytes()
	developerFees := header.GetDeveloperFees().Bytes()
	isStartOfEpochBlock := header.IsStartOfEpochBlock()
	epochStartInfo := getEpochStartInfoForMeta(header)

	return &schema.Block{
		Nonce:                 nonce,
		Round:                 round,
		Epoch:                 epoch,
		Hash:                  hash,
		MiniBlocks:            nil,
		NotarizedBlocksHashes: nil, /*TODO*/
		Proposer:              proposer,
		Validators:            validators,
		PubKeysBitmap:         pubKeysBitmap,
		Size:                  size,
		SizeTxs:               0,
		Timestamp:             timeStamp,
		StateRootHash:         rootHash,
		PrevHash:              prevHash,
		ShardID:               shardID,
		TxCount:               txCount,
		AccumulatedFees:       accumulatedFees,
		DeveloperFees:         developerFees,
		EpochStartBlock:       isStartOfEpochBlock,
		EpochStartInfo:        epochStartInfo,
	}, nil

	body.GetMiniBlocks()

	_ = hash
	//miniblocks from body
	getProposerIndex(signersIndexes)
	_ = signersIndexes
	header.GetPubKeysBitmap()
	body.Size()
	//size txs?
	header.GetTimeStamp()
	header.GetRootHash() // ?? stateroothash?
	header.GetPrevHash()
	header.GetShardID()
	header.GetTxCount()
	header.GetAccumulatedFees()
	header.GetDeveloperFees()
	header.IsStartOfEpochBlock()
	getEpochStartInfoForMeta(header)

	return nil, nil
}

func uIntSliceToIntSlice(in []uint64) []int64 {
	out := make([]int64, len(in))

	for i := range in {
		out[i] = int64(in[i])
	}

	return out
}

func getProposerIndex(signersIndexes []uint64) int64 {
	if len(signersIndexes) > 0 {
		return int64(signersIndexes[0])
	}

	return 0
}

func getEpochStartInfoForMeta(header data.HeaderHandler) *schema.EpochStartInfo {
	if header.GetShardID() != core.MetachainShardId {
		return nil
	}

	metaHeader, ok := header.(*erdBlock.MetaBlock)
	if !ok {
		return nil
	}

	if !metaHeader.IsStartOfEpochBlock() {
		return nil
	}

	metaHeaderEconomics := metaHeader.EpochStart.Economics

	return &schema.EpochStartInfo{
		TotalSupply:                      metaHeaderEconomics.TotalSupply.Bytes(),
		TotalToDistribute:                metaHeaderEconomics.TotalToDistribute.Bytes(),
		TotalNewlyMinted:                 metaHeaderEconomics.TotalNewlyMinted.Bytes(),
		RewardsPerBlock:                  metaHeaderEconomics.RewardsPerBlock.Bytes(),
		RewardsForProtocolSustainability: metaHeaderEconomics.RewardsForProtocolSustainability.Bytes(),
		NodePrice:                        metaHeaderEconomics.NodePrice.Bytes(),
		PrevEpochStartRound:              int64(metaHeaderEconomics.PrevEpochStartRound),
		PrevEpochStartHash:               metaHeaderEconomics.PrevEpochStartHash,
	}
}
