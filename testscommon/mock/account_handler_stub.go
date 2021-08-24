package mock

import "math/big"

// UserAccountStub -
type UserAccountStub struct {
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
	return big.NewInt(10)
}

// AddressBytes calls a custom AddressBytes function if defined, otherwise returns a dummy byte slice
func (us *UserAccountStub) AddressBytes() []byte {
	if us.AddressBytesCalled != nil {
		return us.AddressBytesCalled()
	}
	return []byte("addr")
}

// GetNonce calls a custom GetNonce function if defined, otherwise returns a dummy value
func (us *UserAccountStub) GetNonce() uint64 {
	if us.GetNonceCalled != nil {
		return us.GetNonceCalled()
	}
	return 14
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (us *UserAccountStub) IsInterfaceNil() bool {
	return us == nil
}
