package mock

import (
	"math/big"
	"strconv"
)

// UserAccountStub -
type UserAccountStub struct {
	CurrentBalance     int64
	CurrentNonce       uint64
	GetBalanceCalled   func() *big.Int
	GetNonceCalled     func() uint64
	AddressBytesCalled func() []byte
}

// IncreaseNonce -
func (us *UserAccountStub) IncreaseNonce(_ uint64) {
}

// GetBalance calls a custom GetBalance function if defined, otherwise returns a dummy value
func (us *UserAccountStub) GetBalance() *big.Int {
	if us.GetBalanceCalled != nil {
		return us.GetBalanceCalled()
	}

	us.CurrentBalance++
	return big.NewInt(us.CurrentBalance)
}

// AddressBytes calls a custom AddressBytes function if defined, otherwise returns a dummy byte slice
func (us *UserAccountStub) AddressBytes() []byte {
	if us.AddressBytesCalled != nil {
		return us.AddressBytesCalled()
	}
	return []byte("addr" + strconv.Itoa(int(us.CurrentBalance)))
}

// GetNonce calls a custom GetNonce function if defined, otherwise returns a dummy value
func (us *UserAccountStub) GetNonce() uint64 {
	if us.GetNonceCalled != nil {
		return us.GetNonceCalled()
	}
	us.CurrentNonce++
	return us.CurrentNonce
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (us *UserAccountStub) IsInterfaceNil() bool {
	return us == nil
}

// RetrieveValueFromDataTrieTracker -
func (us *UserAccountStub) RetrieveValueFromDataTrieTracker(key []byte) ([]byte, error) {
	return nil, nil
}
