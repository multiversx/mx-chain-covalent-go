package factory

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/receipts"
)

// ArgsHyperBlockProcessor holds all input dependencies required by hyper block processor factory
// in order to create a new hyper block processor
type ArgsHyperBlockProcessor struct {
}

// CreateHyperBlockProcessor creates a new hyper block processor handler
func CreateHyperBlockProcessor(args *ArgsHyperBlockProcessor) (covalent.HyperBlockProcessor, error) {
	receiptsHandler := receipts.NewReceiptsProcessor()
	_ = receiptsHandler // this will be a subcomp of a tx processor which shall follow in next PRs

	return process.NewHyperBlockProcessor(), nil
}
