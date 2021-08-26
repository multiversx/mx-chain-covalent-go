package factory

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/factory"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
)

// ArgsCovalentIndexerFactory holds all input dependencies required by covalent data indexer factory
// in order to create new instances
type ArgsCovalentIndexerFactory struct {
	Enabled          bool
	PubKeyConverter  core.PubkeyConverter
	Accounts         covalent.AccountsAdapter
	Hasher           hashing.Hasher
	Marshaller       marshal.Marshalizer
	ShardCoordinator process.ShardCoordinator
}

// CreateCovalentIndexer creates a new Driver instance of type covalent data indexer
func CreateCovalentIndexer(args *ArgsCovalentIndexerFactory) (covalent.Driver, error) {
	if check.IfNil(args.PubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}
	if check.IfNil(args.Accounts) {
		return nil, covalent.ErrNilAccountsAdapter
	}
	if check.IfNil(args.Hasher) {
		return nil, covalent.ErrNilHasher
	}
	if check.IfNil(args.Marshaller) {
		return nil, covalent.ErrNilMarshaller
	}

	argsDataProcessor := &factory.ArgsDataProcessor{
		PubKeyConvertor:  args.PubKeyConverter,
		Accounts:         args.Accounts,
		Hasher:           args.Hasher,
		Marshaller:       args.Marshaller,
		ShardCoordinator: args.ShardCoordinator,
	}

	dataProcessor, err := factory.CreateDataProcessor(argsDataProcessor)
	if err != nil {
		return nil, err
	}

	return covalent.NewCovalentDataIndexer(dataProcessor)
}
