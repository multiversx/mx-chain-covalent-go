package processMocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
)

// AlteredAccountsHandlerStub -
type AlteredAccountsHandlerStub struct {
	ProcessAccountsCalled func(apiNotarizedBlocks []*api.NotarizedBlock) ([]*schema.AccountBalanceUpdate, error)
}

// ProcessAccounts -
func (ahs *AlteredAccountsHandlerStub) ProcessAccounts(apiNotarizedBlocks []*api.NotarizedBlock) ([]*schema.AccountBalanceUpdate, error) {
	if ahs.ProcessAccountsCalled != nil {
		return ahs.ProcessAccountsCalled(apiNotarizedBlocks)
	}

	return nil, nil
}
