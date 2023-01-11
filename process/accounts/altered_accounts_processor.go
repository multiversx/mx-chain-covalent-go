package accounts

import (
	"github.com/multiversx/mx-chain-covalent-go/process/utility"
	"github.com/multiversx/mx-chain-covalent-go/schema"
	"github.com/multiversx/mx-chain-core-go/data/outport"
)

type alteredAccountsProcessor struct {
}

// NewAlteredAccountsProcessor creates a new instance of altered accounts processor
func NewAlteredAccountsProcessor() *alteredAccountsProcessor {
	return &alteredAccountsProcessor{}
}

// ProcessAccounts converts accounts data to a specific structure defined by avro schema
func (ap *alteredAccountsProcessor) ProcessAccounts(apiAlteredAccounts []*outport.AlteredAccount) ([]*schema.AccountBalanceUpdate, error) {
	accounts := make([]*schema.AccountBalanceUpdate, 0, len(apiAlteredAccounts))

	for _, apiAlteredAccount := range apiAlteredAccounts {
		if apiAlteredAccount == nil {
			continue
		}

		balance, err := utility.GetBigIntBytesFromStr(apiAlteredAccount.Balance)
		if err != nil {
			return nil, err
		}
		accountTokenData, err := processAccountsTokenData(apiAlteredAccount.Tokens)
		if err != nil {
			return nil, err
		}

		account := &schema.AccountBalanceUpdate{
			Address: []byte(apiAlteredAccount.Address),
			Balance: balance,
			Nonce:   int64(apiAlteredAccount.Nonce),
			Tokens:  tokensOrNil(accountTokenData),
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func processAccountsTokenData(apiTokens []*outport.AccountTokenData) ([]*schema.AccountTokenData, error) {
	tokens := make([]*schema.AccountTokenData, 0, len(apiTokens))

	for _, apiToken := range apiTokens {
		if apiToken == nil {
			continue
		}

		balance, err := utility.GetBigIntBytesFromStr(apiToken.Balance)
		if err != nil {
			return nil, err
		}

		token := &schema.AccountTokenData{
			Nonce:      int64(apiToken.Nonce),
			Identifier: apiToken.Identifier,
			Balance:    balance,
			Properties: apiToken.Properties,
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func tokensOrNil(accounts []*schema.AccountTokenData) []*schema.AccountTokenData {
	if len(accounts) == 0 {
		return nil
	}

	return accounts
}
