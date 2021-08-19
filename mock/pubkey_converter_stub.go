package mock

// PubKeyConverterStub -
type PubKeyConverterStub struct {
	DecodeCalled func(humanReadable string) ([]byte, error)
	EncodeCalled func(pkBytes []byte) string
}

// Len -
func (pcs *PubKeyConverterStub) Len() int {
	return 0
}

// Decode -
func (pcs *PubKeyConverterStub) Decode(humanReadable string) ([]byte, error) {
	if pcs.DecodeCalled != nil {
		return pcs.DecodeCalled(humanReadable)
	}

	return make([]byte, 0), nil
}

// Encode -
func (pcs *PubKeyConverterStub) Encode(pkBytes []byte) string {
	if pcs.EncodeCalled != nil {
		return pcs.EncodeCalled(pkBytes)
	}

	return "erd1" + string(pkBytes)
}

// IsInterfaceNil -
func (pcs *PubKeyConverterStub) IsInterfaceNil() bool {
	return pcs == nil
}
