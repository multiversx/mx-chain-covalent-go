package processMocks

import (
	"github.com/multiversx/mx-chain-covalent-go/schema"
	"github.com/multiversx/mx-chain-core-go/data/outport"
)

// AlteredAccountsHandlerStub -
type AlteredAccountsHandlerStub struct {
	ProcessAccountsCalled func(apiAlteredAccounts []*outport.AlteredAccount) ([]*schema.AccountBalanceUpdate, error)
}

// ProcessAccounts -
func (ahs *AlteredAccountsHandlerStub) ProcessAccounts(apiAlteredAccounts []*outport.AlteredAccount) ([]*schema.AccountBalanceUpdate, error) {
	if ahs.ProcessAccountsCalled != nil {
		return ahs.ProcessAccountsCalled(apiAlteredAccounts)
	}

	return nil, nil
}
