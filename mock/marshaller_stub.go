package mock

import (
	"encoding/json"
	"errors"
)

// ErrMockMarshalizer -
var ErrMockMarshalizer = errors.New("MarshallerStub generic error")

// ErrNilObjectToMarshal -
var ErrNilObjectToMarshal = errors.New("nil object to serialize from")

// MarshallerStub that will be used for testing
type MarshallerStub struct {
	MarshalCalled   func(obj interface{}) ([]byte, error)
	UnmarshalCalled func(obj interface{}, buff []byte) error
}

// Marshal converts the input object in a slice of bytes
func (mm *MarshallerStub) Marshal(obj interface{}) ([]byte, error) {
	if mm.MarshalCalled != nil {
		return mm.MarshalCalled(obj)
	}
	return json.Marshal(obj)
}

// Unmarshal applies the serialized values over an instantiated object
func (mm *MarshallerStub) Unmarshal(obj interface{}, buff []byte) error {
	if mm.UnmarshalCalled != nil {
		return mm.UnmarshalCalled(obj, buff)
	}
	return json.Unmarshal(buff, obj)
}

// IsInterfaceNil returns true if there is no value under the interface
func (mm *MarshallerStub) IsInterfaceNil() bool {
	return mm == nil
}
