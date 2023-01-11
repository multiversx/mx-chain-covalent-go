package facade

import "github.com/elodina/go-avro"

// AvroEncoder should be able to encode any avro schema in a byte array
type AvroEncoder interface {
	Encode(record avro.AvroRecord) ([]byte, error)
}
