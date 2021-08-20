package accounts

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
)

type accountsProcessor struct {
	pubKeyConverter core.PubkeyConverter
	accounts        covalent.AccountsAdapter
}

// NewAccountsProcessor creates a new instance of accounts processor
func NewAccountsProcessor(accounts covalent.AccountsAdapter, pubKeyConverter core.PubkeyConverter) (*accountsProcessor, error) {
	return &accountsProcessor{
		accounts:        accounts,
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// ProcessAccounts converts accounts data to a specific structure defined by avro schema
func (ap *accountsProcessor) ProcessAccounts() ([]*schema.AccountBalanceUpdate, error) {
	return nil, nil
}
