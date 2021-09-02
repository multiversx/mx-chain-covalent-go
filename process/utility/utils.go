package utility

import (
	"bytes"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/elodina/go-avro"
)

// StrSliceToBytesSlice outputs the bytes slice representation of a string slice input
func StrSliceToBytesSlice(in []string) [][]byte {
	out := make([][]byte, len(in))

	for i := range in {
		out[i] = make([]byte, len(in[i]))
		out[i] = []byte(in[i])
	}

	return out
}

// UIntSliceToIntSlice outputs the int64 slice representation of a uint64 slice input
func UIntSliceToIntSlice(in []uint64) []int64 {
	out := make([]int64, len(in))

	for i := range in {
		out[i] = int64(in[i])
	}

	return out
}

// GetBytes returns the bytes representation of a big int input if not nil, otherwise returns []byte{}
func GetBytes(val *big.Int) []byte {
	if val != nil {
		return val.Bytes()
	}

	return big.NewInt(0).Bytes()
}

// EncodePubKey returns a byte slice of the encoded pubKey input, using a pub key converter
func EncodePubKey(pubKeyConverter core.PubkeyConverter, pubKey []byte) []byte {
	return []byte(pubKeyConverter.Encode(pubKey))
}

// Encode returns a byte slice representing the binary encoding of the input avro record
func Encode(record avro.AvroRecord) ([]byte, error) {
	writer := avro.NewSpecificDatumWriter()
	writer.SetSchema(record.Schema())

	buffer := new(bytes.Buffer)
	encoder := avro.NewBinaryEncoder(buffer)

	err := writer.Write(record, encoder)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Decode tries to decode a data buffer, read it and store it on the input record.
// If successfully, the record is filled with data from the buffer, otherwise an error might be returned
func Decode(record avro.AvroRecord, buffer []byte) error {
	reader := avro.NewSpecificDatumReader()
	reader.SetSchema(record.Schema())

	decoder := avro.NewBinaryDecoder(buffer)
	return reader.Read(record, decoder)
}
