package accounts

import (
	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
)

type alteredAccountsProcessor struct {
}

// NewAlteredAccountsProcessor creates a new instance of altered accounts processor
func NewAlteredAccountsProcessor() *alteredAccountsProcessor {
	return &alteredAccountsProcessor{}
}

// ProcessAccounts converts accounts data to a specific structure defined by avro schema
func (ap *alteredAccountsProcessor) ProcessAccounts(apiAlteredAccounts []*api.AlteredAccount) ([]*schemaV2.AccountBalanceUpdate, error) {
	accounts := make([]*schemaV2.AccountBalanceUpdate, 0, len(apiAlteredAccounts))

	for _, apiAccount := range apiAlteredAccounts {
		if apiAccount == nil {
			continue
		}

		account, err := processAccount(apiAccount)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func processAccount(apiAccount *api.AlteredAccount) (*schemaV2.AccountBalanceUpdate, error) {
	balance, err := utility.GetBigIntBytesFromStr(apiAccount.Balance)
	if err != nil {
		return nil, err
	}

	return &schemaV2.AccountBalanceUpdate{
		Address: []byte(apiAccount.Address),
		Balance: balance,
		Nonce:   int64(apiAccount.Nonce),
	}, nil
}
