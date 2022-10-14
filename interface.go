package covalent

import (
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/elodina/go-avro"
)

// HyperBlockProcessor shall handle hyper block processing into avro schema blocks
type HyperBlockProcessor interface {
	Process(hyperBlock *hyperBlock.HyperBlock) (*schema.HyperBlock, error)
}

// AvroMarshaller defines what an avro marshaller should do
type AvroMarshaller interface {
	Encode(record avro.AvroRecord) ([]byte, error)
	Decode(record avro.AvroRecord, buffer []byte) error
}
