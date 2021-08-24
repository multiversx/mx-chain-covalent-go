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

// GetBalance -
func (us *UserAccountStub) GetBalance() *big.Int {
	if us.GetBalanceCalled != nil {
		return us.GetBalanceCalled()
	}
	return big.NewInt(10)
}

// AddressBytes -
func (us *UserAccountStub) AddressBytes() []byte {
	if us.AddressBytesCalled != nil {
		return us.AddressBytesCalled()
	}
	return []byte("addr")
}

// GetNonce -
func (us *UserAccountStub) GetNonce() uint64 {
	if us.GetNonceCalled != nil {
		return us.GetNonceCalled()
	}
	return 14
}

// IsInterfaceNil -
func (us *UserAccountStub) IsInterfaceNil() bool {
	return us == nil
}
