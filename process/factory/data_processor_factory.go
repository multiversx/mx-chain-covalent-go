package factory

import (
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	blockCovalent "github.com/ElrondNetwork/covalent-indexer-go/process/block"
	"github.com/ElrondNetwork/elrond-go-core/core"
)

type ArgsDataProcessor struct {
	PubKeyConvertor core.PubkeyConverter
}

func CreateDataProcessor(args *ArgsDataProcessor) (*process.DataProcessor, error) {
	blockHandler, _ := blockCovalent.NewBlockProcessor()

	return process.NewDataProcessor(blockHandler)
}
