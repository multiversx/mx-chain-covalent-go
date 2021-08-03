package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

type BlockHandler interface {
	ProcessBlock(block data.BodyHandler) (*schema.Block, error)
}
