package covalent

import (
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/elodina/go-avro"
)

type DataHandler interface {
	ProcessData(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error)
}

// HyperBlockProcessor shall handle hyper block processing into avro schema blocks
type HyperBlockProcessor interface {
	Process(hyperBlock *hyperBlock.HyperBlock) (*schema.BlockResult, error)
}

type AccountsAdapter interface {
	LoadAccount(address []byte) (vmcommon.AccountHandler, error)
	IsInterfaceNil() bool
}

// AvroMarshaller defines what an avro marshaller should do
type AvroMarshaller interface {
	Encode(record avro.AvroRecord) ([]byte, error)
	Decode(record avro.AvroRecord, buffer []byte) error
}
