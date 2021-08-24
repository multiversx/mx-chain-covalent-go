package mock

import vmcommon "github.com/ElrondNetwork/elrond-vm-common"

type AccountsAdapterStub struct {
	LoadAccountCalled func(address []byte) (vmcommon.AccountHandler, error)
}

func (as *AccountsAdapterStub) LoadAccount(address []byte) (vmcommon.AccountHandler, error) {
	if as.LoadAccountCalled != nil {
		return as.LoadAccountCalled(address)
	}
	return nil, nil
}

func (as *AccountsAdapterStub) IsInterfaceNil() bool {
	return as == nil
}
