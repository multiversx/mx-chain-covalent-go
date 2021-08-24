package mock

import (
	"math/big"
	"strconv"
)

// UserAccountStub -
type UserAccountStub struct {
	CurrentBalance int64
	CurrentNonce   uint64
}

// IncreaseNonce -
func (uas *UserAccountStub) IncreaseNonce(_ uint64) {
}

// GetBalance increments CurrentBalance and returns it as a big int
func (uas *UserAccountStub) GetBalance() *big.Int {
	uas.CurrentBalance++
	return big.NewInt(uas.CurrentBalance)
}

// AddressBytes returns a byte slice of ("addr" + CurrentBalance)
func (uas *UserAccountStub) AddressBytes() []byte {
	return []byte("addr" + strconv.Itoa(int(uas.CurrentBalance)))
}

// GetNonce increments CurrentNonce and returns it
func (uas *UserAccountStub) GetNonce() uint64 {
	uas.CurrentNonce++
	return uas.CurrentNonce
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (uas *UserAccountStub) IsInterfaceNil() bool {
	return uas == nil
}

// RetrieveValueFromDataTrieTracker -
func (uas *UserAccountStub) RetrieveValueFromDataTrieTracker(key []byte) ([]byte, error) {
	return nil, nil
}
