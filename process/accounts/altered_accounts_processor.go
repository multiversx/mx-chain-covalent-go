package accounts

import (
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/esdt"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
)

type alteredAccountsProcessor struct {
}

// NewAlteredAccountsProcessor creates a new instance of altered accounts processor
func NewAlteredAccountsProcessor() *alteredAccountsProcessor {
	return &alteredAccountsProcessor{}
}

// ProcessAccounts converts accounts data to a specific structure defined by avro schema
func (ap *alteredAccountsProcessor) ProcessAccounts(notarizedBlocks []*api.NotarizedBlock) ([]*schema.AccountBalanceUpdate, error) {
	accounts := make([]*schema.AccountBalanceUpdate, 0, len(notarizedBlocks))

	for _, block := range notarizedBlocks {
		if block == nil {
			continue
		}

		accountsInBlock, err := processAlteredAccounts(block.AlteredAccounts)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, accountsInBlock...)
	}

	return accounts, nil
}

func processAlteredAccounts(apiAlteredAccounts []*outport.AlteredAccount) ([]*schema.AccountBalanceUpdate, error) {
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
			Address:          []byte(apiAlteredAccount.Address),
			Balance:          balance,
			Nonce:            int64(apiAlteredAccount.Nonce),
			AccountTokenData: accountTokenData,
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
			MetaData:   processMetaData(apiToken.MetaData),
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func processMetaData(apiMetaData *esdt.MetaData) *schema.MetaData {
	return &schema.MetaData{
		Nonce:      int64(apiMetaData.Nonce),
		Name:       apiMetaData.Name,
		Creator:    apiMetaData.Creator,
		Royalties:  int32(apiMetaData.Royalties),
		Hash:       apiMetaData.Hash,
		URIs:       apiMetaData.URIs,
		Attributes: apiMetaData.Attributes,
	}
}
