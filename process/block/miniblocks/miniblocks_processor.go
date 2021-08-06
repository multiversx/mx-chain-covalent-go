package miniblocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
)

type miniBlocksProcessor struct {
	hasher      hashing.Hasher
	marshalizer marshal.Marshalizer
}

// NewMiniBlocksProcessor will create a new instance of miniBlocksProcessor
func NewMiniBlocksProcessor(
	hasher hashing.Hasher,
	marshalizer marshal.Marshalizer,
) (*miniBlocksProcessor, error) {
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

func (mbp *miniBlocksProcessor) ProcessMiniBlocks() ([]*schema.MiniBlock, error) {

}
