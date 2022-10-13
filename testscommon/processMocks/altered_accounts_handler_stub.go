package processMocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go/alteredAccount"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
)

// AlteredAccountsHandlerStub -
type AlteredAccountsHandlerStub struct {
	ProcessAccountsCalled func(apiAlteredAccounts []*alteredAccount.AlteredAccount) ([]*schema.AccountBalanceUpdate, error)
}

// ProcessAccounts -
func (ahs *AlteredAccountsHandlerStub) ProcessAccounts(apiAlteredAccounts []*alteredAccount.AlteredAccount) ([]*schema.AccountBalanceUpdate, error) {
	if ahs.ProcessAccountsCalled != nil {
		return ahs.ProcessAccountsCalled(apiAlteredAccounts)
	}

	return nil, nil
}
