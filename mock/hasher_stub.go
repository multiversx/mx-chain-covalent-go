package mock

// HasherStub that will be used for testing
type HasherStub struct {
}

// Compute will output a dummy hash
func (hs HasherStub) Compute(s string) []byte {
	return []byte("ok")
}

// Size returns a dummy size
func (HasherStub) Size() int {
	return 123
}

// IsInterfaceNil returns false
func (hs HasherStub) IsInterfaceNil() bool {
	return false
}
