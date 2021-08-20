package miniblocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("process/block/miniBlocks/miniBlocksProcessor")

type miniBlocksProcessor struct {
	hasher     hashing.Hasher
	marshaller marshal.Marshalizer
}

// NewMiniBlocksProcessor will create a new instance of miniBlocksProcessor
func NewMiniBlocksProcessor(hasher hashing.Hasher, marshalizer marshal.Marshalizer) (*miniBlocksProcessor, error) {
	if check.IfNil(marshalizer) {
		return nil, covalent.ErrNilMarshaller
	}
	if check.IfNil(hasher) {
		return nil, covalent.ErrNilHasher
	}

	return &miniBlocksProcessor{
		hasher:     hasher,
		marshaller: marshalizer,
	}, nil
}

// ProcessMiniBlocks converts mini blocks core data to a specific mini blocks structure array defined by avro schema
func (mbp *miniBlocksProcessor) ProcessMiniBlocks(header data.HeaderHandler, body data.BodyHandler) ([]*schema.MiniBlock, error) {
	erdBody, castOk := body.(*block.Body)
	if !castOk {
		return nil, covalent.ErrBlockBodyAssertion
	}

	miniBlocks := make([]*schema.MiniBlock, 0)
	erdMiniBlocks := erdBody.GetMiniBlocks()

	for _, mb := range erdMiniBlocks {

		miniBlock, err := mbp.processMiniBlock(mb, header)
		if err != nil {
			log.Warn("miniBlocksProcessor.ProcessMiniBlocks cannot process miniBlock", "error", err)
			continue
		}

		miniBlocks = append(miniBlocks, miniBlock)
	}

	return miniBlocks, nil
}

func (mbp *miniBlocksProcessor) processMiniBlock(miniBlock *block.MiniBlock, header data.HeaderHandler) (*schema.MiniBlock, error) {
	miniBlockHash, err := core.CalculateHash(mbp.marshaller, mbp.hasher, miniBlock)
	if err != nil {
		return nil, err
	}

	return &schema.MiniBlock{
		Hash:            miniBlockHash,
		SenderShardID:   int32(miniBlock.SenderShardID),
		ReceiverShardID: int32(miniBlock.ReceiverShardID),
		Type:            int32(miniBlock.Type),
		Timestamp:       int64(header.GetTimeStamp()),
	}, nil
}
