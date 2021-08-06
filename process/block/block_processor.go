package block

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	erdBlock "github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
)

type blockProcessor struct {
	hasher      hashing.Hasher
	marshalizer marshal.Marshalizer
}

// NewBlockProcessor creates a new instance of block processor
func NewBlockProcessor(hasher hashing.Hasher, marshalizer marshal.Marshalizer) (*blockProcessor, error) {
	if check.IfNil(hasher) {
		return nil, covalent.ErrNilHasher
	}
	if check.IfNil(marshalizer) {
		return nil, covalent.ErrNilMarshalizer
	}

	return &blockProcessor{
		hasher:      hasher,
		marshalizer: marshalizer,
	}, nil
}

// ProcessBlock converts block data to a specific structure defined by avro schema
func (bp *blockProcessor) ProcessBlock(args *indexer.ArgsSaveBlockData) (*schema.Block, error) {
	body, ok := args.Body.(*erdBlock.Body)
	if !ok {
		return nil, covalent.ErrBlockBodyAssertion
	}
	_ = body
	nonce := int64(args.Header.GetNonce())
	round := int64(args.Header.GetRound())
	epoch := int32(args.Header.GetEpoch())
	hash := args.HeaderHash
	notarizedBlocksHashes := strSliceToBytesSlice(args.NotarizedHeadersHashes)
	proposer := getProposerIndex(args.SignersIndexes)
	validators := uIntSliceToIntSlice(args.SignersIndexes)
	pubKeysBitmap := args.Header.GetPubKeysBitmap()
	size := int64(321321321) // TODO
	timeStamp := int64(args.Header.GetTimeStamp())
	rootHash := args.Header.GetRootHash()
	prevHash := args.Header.GetPrevHash()
	shardID := int32(args.Header.GetShardID())
	txCount := int32(args.Header.GetTxCount())
	accumulatedFees := args.Header.GetAccumulatedFees().Bytes()
	developerFees := args.Header.GetDeveloperFees().Bytes()
	isStartOfEpochBlock := args.Header.IsStartOfEpochBlock()
	epochStartInfo := getEpochStartInfo(args.Header)

	return &schema.Block{
		Nonce:                 nonce,
		Round:                 round,
		Epoch:                 epoch,
		Hash:                  hash,
		MiniBlocks:            nil, /*TODO*/
		NotarizedBlocksHashes: notarizedBlocksHashes,
		Proposer:              proposer,
		Validators:            validators,
		PubKeysBitmap:         pubKeysBitmap,
		Size:                  size, /*TODO*/
		SizeTxs:               0,    /*TODO*/
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
}

func strSliceToBytesSlice(in []string) [][]byte {
	out := make([][]byte, len(in))

	for i := range in {
		out[i] = make([]byte, len(in[i]))
		tmp := []byte(in[i])
		out = append(out, tmp)
	}

	return out
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

func getEpochStartInfo(header data.HeaderHandler) *schema.EpochStartInfo {

	economics, exists := getEconomicsIfExists(header)
	if !exists {
		return nil
	}

	return &schema.EpochStartInfo{
		TotalSupply:                      economics.TotalSupply.Bytes(),
		TotalToDistribute:                economics.TotalToDistribute.Bytes(),
		TotalNewlyMinted:                 economics.TotalNewlyMinted.Bytes(),
		RewardsPerBlock:                  economics.RewardsPerBlock.Bytes(),
		RewardsForProtocolSustainability: economics.RewardsForProtocolSustainability.Bytes(),
		NodePrice:                        economics.NodePrice.Bytes(),
		PrevEpochStartRound:              int64(economics.PrevEpochStartRound),
		PrevEpochStartHash:               economics.PrevEpochStartHash,
	}
}

func getEconomicsIfExists(header data.HeaderHandler) (*erdBlock.Economics, bool) {
	if header.GetShardID() != core.MetachainShardId {
		return nil, false
	}

	metaHeader, ok := header.(*erdBlock.MetaBlock)
	if !ok {
		return nil, false
	}

	if !metaHeader.IsStartOfEpochBlock() {
		return nil, false
	}

	return &metaHeader.EpochStart.Economics, true
}

func (bp *blockProcessor) computeBlockSizeAndHeaderHash(header data.HeaderHandler, body *erdBlock.Body) (int, []byte, error) {
	headerBytes, err := bp.marshalizer.Marshal(header)
	if err != nil {
		return 0, nil, err
	}
	bodyBytes, err := bp.marshalizer.Marshal(body)
	if err != nil {
		return 0, nil, err
	}

	blockSize := len(headerBytes) + len(bodyBytes)

	headerHash := bp.hasher.Compute(string(headerBytes))

	return blockSize, headerHash, nil
}
