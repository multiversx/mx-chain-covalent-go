package mock

import "github.com/elodina/go-avro"

type AvroEncoderStub struct {
	EncodeCalled func(record avro.AvroRecord) ([]byte, error)
}

func (aes *AvroEncoderStub) Encode(record avro.AvroRecord) ([]byte, error) {
	if aes.EncodeCalled != nil {
		return aes.EncodeCalled(record)
	}

	return nil, nil
}
