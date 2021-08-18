package block

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	erdBlock "github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
)

const DefaultProposerIndex = int64(0)

type blockProcessor struct {
	marshaller        marshal.Marshalizer
	miniBlocksHandler process.MiniBlockHandler
}

// NewBlockProcessor creates a new instance of block processor
func NewBlockProcessor(marshaller marshal.Marshalizer, mbHandler process.MiniBlockHandler) (*blockProcessor, error) {
	if check.IfNil(marshaller) {
		return nil, covalent.ErrNilMarshaller
	}
	if mbHandler == nil {
		return nil, covalent.ErrNilMiniBlockHandler
	}

	return &blockProcessor{
		marshaller:        marshaller,
		miniBlocksHandler: mbHandler,
	}, nil
}

// ProcessBlock converts block data to a specific structure defined by avro schema
func (bp *blockProcessor) ProcessBlock(args *indexer.ArgsSaveBlockData) (*schema.Block, error) {
	body, ok := args.Body.(*erdBlock.Body)
	if !ok {
		return nil, covalent.ErrBlockBodyAssertion
	}

	blockSizeInBytes, err := bp.computeBlockSize(args.Header, body)
	if err != nil {
		return nil, err
	}

	miniBlocks, err := bp.miniBlocksHandler.ProcessMiniBlocks(args.Header, body)
	if err != nil {
		return nil, err
	}

	nonce := int64(args.Header.GetNonce())
	round := int64(args.Header.GetRound())
	epoch := int32(args.Header.GetEpoch())
	hash := args.HeaderHash
	notarizedBlocksHashes := utility.StrSliceToBytesSlice(args.NotarizedHeadersHashes)
	proposer := getProposerIndex(args.SignersIndexes)
	validators := utility.UIntSliceToIntSlice(args.SignersIndexes)
	pubKeysBitmap := args.Header.GetPubKeysBitmap()
	timeStamp := int64(args.Header.GetTimeStamp())
	rootHash := args.Header.GetRootHash()
	prevHash := args.Header.GetPrevHash()
	shardID := int32(args.Header.GetShardID())
	txCount := int32(args.Header.GetTxCount())
	accumulatedFees := utility.GetBytes(args.Header.GetAccumulatedFees())
	developerFees := utility.GetBytes(args.Header.GetDeveloperFees())
	isStartOfEpochBlock := args.Header.IsStartOfEpochBlock()
	epochStartInfo := getEpochStartInfo(args.Header)

	return &schema.Block{
		Nonce:                 nonce,
		Round:                 round,
		Epoch:                 epoch,
		Hash:                  hash,
		MiniBlocks:            miniBlocks,
		NotarizedBlocksHashes: notarizedBlocksHashes,
		Proposer:              proposer,
		Validators:            validators,
		PubKeysBitmap:         pubKeysBitmap,
		Size:                  blockSizeInBytes,
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

func (bp *blockProcessor) computeBlockSize(header data.HeaderHandler, body *erdBlock.Body) (int64, error) {
	headerBytes, err := bp.marshaller.Marshal(header)
	if err != nil {
		return 0, err
	}
	bodyBytes, err := bp.marshaller.Marshal(body)
	if err != nil {
		return 0, err
	}

	blockSize := len(headerBytes) + len(bodyBytes)

	return int64(blockSize), nil
}

func getProposerIndex(signersIndexes []uint64) int64 {
	if len(signersIndexes) > 0 {
		return int64(signersIndexes[0])
	}

	return DefaultProposerIndex
}

func getEpochStartInfo(header data.HeaderHandler) *schema.EpochStartInfo {
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

	economics := metaHeader.EpochStart.Economics

	return &schema.EpochStartInfo{
		TotalSupply:                      utility.GetBytes(economics.TotalSupply),
		TotalToDistribute:                utility.GetBytes(economics.TotalToDistribute),
		TotalNewlyMinted:                 utility.GetBytes(economics.TotalNewlyMinted),
		RewardsPerBlock:                  utility.GetBytes(economics.RewardsPerBlock),
		RewardsForProtocolSustainability: utility.GetBytes(economics.RewardsForProtocolSustainability),
		NodePrice:                        utility.GetBytes(economics.NodePrice),
		PrevEpochStartRound:              int32(economics.PrevEpochStartRound),
		PrevEpochStartHash:               economics.PrevEpochStartHash,
	}
}
