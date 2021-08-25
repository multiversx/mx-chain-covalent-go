package utility

import (
	"bytes"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/elodina/go-avro"
	"math/big"
)

func StrSliceToBytesSlice(in []string) [][]byte {
	out := make([][]byte, len(in))

	for i := range in {
		out[i] = make([]byte, len(in[i]))
		out[i] = []byte(in[i])
	}

	return out
}

func UIntSliceToIntSlice(in []uint64) []int64 {
	out := make([]int64, len(in))

	for i := range in {
		out[i] = int64(in[i])
	}

	return out
}

func GetBytes(val *big.Int) []byte {
	if val != nil {
		return val.Bytes()
	}

	return nil
}

func EncodePubKey(pubKeyConverter core.PubkeyConverter, pubKey []byte) []byte {
	return []byte(pubKeyConverter.Encode(pubKey))
}

func Write(record avro.AvroRecord) ([]byte, error) {
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

func Decode(record avro.AvroRecord, buffer []byte) error {
	reader := avro.NewSpecificDatumReader()
	reader.SetSchema(record.Schema())

	decoder := avro.NewBinaryDecoder(buffer)
	return reader.Read(record, decoder)
}
