package mock

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
)

type MiniBlockHandlerStub struct {
}

func (mbhs *MiniBlockHandlerStub) ProcessMiniBlocks(headerHash []byte, header data.HeaderHandler, body *block.Body) ([]*schema.MiniBlock, error) {
	return nil, nil
}
