package process

import (
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
)

type hyperBlockProcessor struct {
	transactionProcessor TransactionHandler
	shardBlocksProcessor ShardBlocksHandler
}

// NewHyperBlockProcessor will create a new instance of an hyper block processor
func NewHyperBlockProcessor(
	transactionHandler TransactionHandler,
	shardBlockHandler ShardBlocksHandler,
) (*hyperBlockProcessor, error) {
	if transactionHandler == nil {
		return nil, errNilTransactionHandler
	}
	if shardBlockHandler == nil {
		return nil, errNilShardBlocksHandler
	}

	return &hyperBlockProcessor{
		transactionProcessor: transactionHandler,
		shardBlocksProcessor: shardBlockHandler,
	}, nil
}

// Process will process current hyper block and convert it to an avro schema block result
func (hbp *hyperBlockProcessor) Process(hyperBlock *hyperBlock.HyperBlock) (*schemaV2.HyperBlock, error) {
	avroHyperBlock := schemaV2.NewHyperBlock()

	txs, err := hbp.transactionProcessor.ProcessTransactions(hyperBlock.Transactions)
	if err != nil {
		return nil, err
	}
	shardBlocks, err := hbp.shardBlocksProcessor.ProcessShardBlocks(hyperBlock.ShardBlocks)
	if err != nil {
		return nil, err
	}

	avroHyperBlock.ShardBlocks = shardBlocks
	avroHyperBlock.Transactions = txs
	return avroHyperBlock, nil
}
