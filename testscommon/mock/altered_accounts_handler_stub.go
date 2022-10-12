package mock

import (
	covalentApi "github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
)

// AlteredAccountsHandlerStub -
type AlteredAccountsHandlerStub struct {
	ProcessAccountsCalled func(apiAlteredAccounts []*covalentApi.AlteredAccount) ([]*schemaV2.AccountBalanceUpdate, error)
}

// ProcessAccounts -
func (ahs *AlteredAccountsHandlerStub) ProcessAccounts(apiAlteredAccounts []*covalentApi.AlteredAccount) ([]*schemaV2.AccountBalanceUpdate, error) {
	if ahs.ProcessAccountsCalled != nil {
		return ahs.ProcessAccountsCalled(apiAlteredAccounts)
	}

	return nil, nil
}
