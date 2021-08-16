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

var log = logger.GetOrCreate("process/block/miniblocks/miniBlocksProcessor")

type miniBlocksProcessor struct {
	hasher      hashing.Hasher
	marshalizer marshal.Marshalizer
}

// NewMiniBlocksProcessor will create a new instance of miniBlocksProcessor
func NewMiniBlocksProcessor(hasher hashing.Hasher, marshalizer marshal.Marshalizer) (*miniBlocksProcessor, error) {
	if check.IfNil(marshalizer) {
		return nil, covalent.ErrNilMarshalizer
	}
	if check.IfNil(hasher) {
		return nil, covalent.ErrNilHasher
	}

	return &miniBlocksProcessor{
		hasher:      hasher,
		marshalizer: marshalizer,
	}, nil
}

func (mbp *miniBlocksProcessor) ProcessMiniBlocks(headerHash []byte, header data.HeaderHandler, body *block.Body) ([]*schema.MiniBlock, error) {

	miniBlocks := make([]*schema.MiniBlock, 0)

	for _, mb := range body.MiniBlocks {

		miniBlock, err := mbp.processMiniBlock(mb, header, headerHash)
		if err != nil {
			log.Warn("miniBlocksProcessor.ProcessMiniBlocks cannot process miniBlock", "error", err)
			continue
		}

		miniBlocks = append(miniBlocks, miniBlock)
	}

	return miniBlocks, nil
}

func (mbp *miniBlocksProcessor) processMiniBlock(miniBlock *block.MiniBlock, header data.HeaderHandler, headerHash []byte) (*schema.MiniBlock, error) {

	mbHash, err := core.CalculateHash(mbp.marshalizer, mbp.hasher, miniBlock)
	if err != nil {
		return nil, err
	}

	covalentMiniBlock := &schema.MiniBlock{
		Hash:            mbHash,
		SenderShardID:   int32(miniBlock.SenderShardID),
		ReceiverShardID: int32(miniBlock.ReceiverShardID),
		Type:            int32(miniBlock.Type),
		Timestamp:       int64(header.GetTimeStamp()),
	}

	if covalentMiniBlock.SenderShardID == int32(header.GetShardID()) {
		covalentMiniBlock.SenderBlockHash = headerHash
	} else {
		covalentMiniBlock.ReceiverBlockHash = headerHash
	}

	if covalentMiniBlock.SenderShardID == covalentMiniBlock.ReceiverShardID {
		covalentMiniBlock.ReceiverBlockHash = headerHash
	}

	return covalentMiniBlock, nil
}
