package factory

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/factory"
	"github.com/ElrondNetwork/elrond-go-core/core"
)

type ArgsCovalentIndexerFactory struct {
	PubKeyConverter core.PubkeyConverter
}

func CreateCovalentIndexer(args *ArgsCovalentIndexerFactory) (covalent.Driver, error) {

	argsDataProcessor := &factory.ArgsDataProcessor{
		PubKeyConvertor: args.PubKeyConverter,
	}

	dataProcessor, err := factory.CreateDataProcessor(argsDataProcessor)
	if err != nil {
		return nil, err
	}

	return covalent.NewCovalentDataIndexer(dataProcessor)
}
