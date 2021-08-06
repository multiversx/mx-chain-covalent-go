package block

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	erdBlock "github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"math/big"
)

var log = logger.GetOrCreate("indexer/workItems")

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

	blockSizeInBytes, err := bp.computeBlockSize(args.Header, body)
	if err != nil {
		return nil, err
	}

	txsSizeInBytes := bp.computeSizeOfTxs(args.TransactionsPool)

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
	accumulatedFees := getBlockField(args.Header.GetAccumulatedFees())
	developerFees := getBlockField(args.Header.GetDeveloperFees())
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
		Size:                  blockSizeInBytes,
		SizeTxs:               txsSizeInBytes,
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

func getBlockField(val *big.Int) []byte {
	if val != nil {
		return val.Bytes()
	}

	return nil
}

func getProposerIndex(signersIndexes []uint64) int64 {
	if len(signersIndexes) > 0 {
		return int64(signersIndexes[0])
	}

	return 0
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
		TotalSupply:                      getBlockField(economics.TotalSupply),
		TotalToDistribute:                getBlockField(economics.TotalToDistribute),
		TotalNewlyMinted:                 getBlockField(economics.TotalNewlyMinted),
		RewardsPerBlock:                  getBlockField(economics.RewardsPerBlock),
		RewardsForProtocolSustainability: getBlockField(economics.RewardsForProtocolSustainability),
		NodePrice:                        getBlockField(economics.NodePrice),
		PrevEpochStartRound:              int64(economics.PrevEpochStartRound),
		PrevEpochStartHash:               economics.PrevEpochStartHash,
	}
}

func (bp *blockProcessor) computeBlockSize(header data.HeaderHandler, body *erdBlock.Body) (int64, error) {
	headerBytes, err := bp.marshalizer.Marshal(header)
	if err != nil {
		return 0, err
	}
	bodyBytes, err := bp.marshalizer.Marshal(body)
	if err != nil {
		return 0, err
	}

	blockSize := len(headerBytes) + len(bodyBytes)

	return int64(blockSize), nil
}

func (bp *blockProcessor) computeSizeOfTxs(pool *indexer.Pool) int64 {
	if pool == nil {
		return 0
	}

	sizeTxs := 0
	sizeTxs += computeSizeOfMap(bp.marshalizer, pool.Txs)
	sizeTxs += computeSizeOfMap(bp.marshalizer, pool.Receipts)
	sizeTxs += computeSizeOfMap(bp.marshalizer, pool.Invalid)
	sizeTxs += computeSizeOfMap(bp.marshalizer, pool.Rewards)
	sizeTxs += computeSizeOfMap(bp.marshalizer, pool.Scrs)

	return int64(sizeTxs)
}

func computeSizeOfMap(marshalizer marshal.Marshalizer, mapTxs map[string]data.TransactionHandler) int {
	txsSize := 0
	for _, tx := range mapTxs {
		txBytes, err := marshalizer.Marshal(tx)
		if err != nil {
			log.Debug("block processor, computeSizeOfMap", "error", err)
			continue
		}

		txsSize += len(txBytes)
	}

	return txsSize
}
