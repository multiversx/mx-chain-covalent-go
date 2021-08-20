package factory

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/factory"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
)

// ArgsCovalentIndexerFactory holds all input dependencies required by covalent data indexer factory
// in order to create new instances
type ArgsCovalentIndexerFactory struct {
	PubKeyConverter core.PubkeyConverter
	Accounts        covalent.AccountsAdapter
}

// CreateCovalentIndexer creates a new Driver instance of type covalent data indexer
func CreateCovalentIndexer(args *ArgsCovalentIndexerFactory) (covalent.Driver, error) {
	if check.IfNil(args.PubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}
	if check.IfNil(args.Accounts) {
		return nil, covalent.ErrNilAccountsAdapter
	}

	argsDataProcessor := &factory.ArgsDataProcessor{
		PubKeyConvertor: args.PubKeyConverter,
		Accounts:        args.Accounts,
	}

	dataProcessor, err := factory.CreateDataProcessor(argsDataProcessor)
	if err != nil {
		return nil, err
	}

	return covalent.NewCovalentDataIndexer(dataProcessor)
}
