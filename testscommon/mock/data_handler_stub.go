package mock

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
)

type DataHandlerStub struct {
	ProcessDataCalled func(args *outport.ArgsSaveBlockData) (*schema.BlockResult, error)
}

func (dhs *DataHandlerStub) ProcessData(args *outport.ArgsSaveBlockData) (*schema.BlockResult, error) {
	if dhs.ProcessDataCalled != nil {
		return dhs.ProcessDataCalled(args)
	}
	return nil, nil
}
