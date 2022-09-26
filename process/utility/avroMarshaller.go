package utility

import (
	"bytes"

	"github.com/elodina/go-avro"
)

type AvroMarshaller struct {
}

// Encode returns a byte slice representing the binary encoding of the input avro record
func (av *AvroMarshaller) Encode(record avro.AvroRecord) ([]byte, error) {
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
func (av *AvroMarshaller) Decode(record avro.AvroRecord, buffer []byte) error {
	reader := avro.NewSpecificDatumReader()
	reader.SetSchema(record.Schema())

	decoder := avro.NewBinaryDecoder(buffer)
	return reader.Read(record, decoder)
}
