package processMocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
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
