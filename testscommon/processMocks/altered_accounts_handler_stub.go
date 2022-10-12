package processMocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go/alteredAccount"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
)

// AlteredAccountsHandlerStub -
type AlteredAccountsHandlerStub struct {
	ProcessAccountsCalled func(apiAlteredAccounts []*alteredAccount.AlteredAccount) ([]*schemaV2.AccountBalanceUpdate, error)
}

// ProcessAccounts -
func (ahs *AlteredAccountsHandlerStub) ProcessAccounts(apiAlteredAccounts []*alteredAccount.AlteredAccount) ([]*schemaV2.AccountBalanceUpdate, error) {
	if ahs.ProcessAccountsCalled != nil {
		return ahs.ProcessAccountsCalled(apiAlteredAccounts)
	}

	return nil, nil
}
