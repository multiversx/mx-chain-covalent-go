package accounts

import (
	"github.com/ElrondNetwork/covalent-indexer-go/alteredAccount"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
)

type alteredAccountsProcessor struct {
}

// NewAlteredAccountsProcessor creates a new instance of altered accounts processor
func NewAlteredAccountsProcessor() *alteredAccountsProcessor {
	return &alteredAccountsProcessor{}
}

// ProcessAccounts converts accounts data to a specific structure defined by avro schema
func (ap *alteredAccountsProcessor) ProcessAccounts(apiAlteredAccounts []*alteredAccount.AlteredAccount) ([]*schema.AccountBalanceUpdate, error) {
	accounts := make([]*schema.AccountBalanceUpdate, 0, len(apiAlteredAccounts))

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

func processAccount(apiAccount *alteredAccount.AlteredAccount) (*schema.AccountBalanceUpdate, error) {
	balance, err := utility.GetBigIntBytesFromStr(apiAccount.Balance)
	if err != nil {
		return nil, err
	}

	return &schema.AccountBalanceUpdate{
		Address: []byte(apiAccount.Address),
		Balance: balance,
		Nonce:   int64(apiAccount.Nonce),
	}, nil
}
