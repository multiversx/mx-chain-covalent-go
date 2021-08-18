package mock

// HasherStub that will be used for testing
type HasherStub struct {
}

// Compute outputs a constant dummy hash
func (hs *HasherStub) Compute(s string) []byte {
	return []byte("ok")
}

// Size returns a dummy size
func (hs *HasherStub) Size() int {
	return 123
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (hs *HasherStub) IsInterfaceNil() bool {
	return hs == nil
}
