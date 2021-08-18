package mock

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
)

// MiniBlockHandlerStub that will be used for testing
type MiniBlockHandlerStub struct {
	ProcessMiniBlockCalled func(header data.HeaderHandler, body *block.Body) ([]*schema.MiniBlock, error)
}

// ProcessMiniBlocks calls a custom mini blocks process function if defined, otherwise returns nil, nil
func (mbhs *MiniBlockHandlerStub) ProcessMiniBlocks(header data.HeaderHandler, body *block.Body) ([]*schema.MiniBlock, error) {
	if mbhs.ProcessMiniBlockCalled != nil {
		return mbhs.ProcessMiniBlockCalled(header, body)
	}

	return nil, nil
}
