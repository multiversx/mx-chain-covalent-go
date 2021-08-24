package mock

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// AccountsAdapterStub -
type AccountsAdapterStub struct {
	UserAccountHandler vmcommon.AccountHandler
	LoadAccountCalled  func(address []byte) (vmcommon.AccountHandler, error)
}

// LoadAccount calls a custom load account function if defined, otherwise returns UserAccountStub, nil
func (as *AccountsAdapterStub) LoadAccount(address []byte) (vmcommon.AccountHandler, error) {
	if as.LoadAccountCalled != nil {
		return as.LoadAccountCalled(address)
	}
	return as.UserAccountHandler, nil
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (as *AccountsAdapterStub) IsInterfaceNil() bool {
	return as == nil
}
