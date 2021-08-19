package transactions

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	erdBlock "github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("process/transactions/transactionProcessor")

type transactionProcessor struct {
	hasher          hashing.Hasher
	marshaller      marshal.Marshalizer
	pubKeyConverter core.PubkeyConverter
}

// NewTransactionProcessor creates a new instance of transactions processor
func NewTransactionProcessor(
	pubKeyConverter core.PubkeyConverter,
	hasher hashing.Hasher,
	marshaller marshal.Marshalizer) (*transactionProcessor, error) {
	if check.IfNil(pubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}
	if check.IfNil(marshaller) {
		return nil, covalent.ErrNilMarshaller
	}
	if check.IfNil(hasher) {
		return nil, covalent.ErrNilHasher
	}

	return &transactionProcessor{
		pubKeyConverter: pubKeyConverter,
		hasher:          hasher,
		marshaller:      marshaller,
	}, nil
}

// ProcessTransactions converts transactions data to a specific structure defined by avro schema
func (txp *transactionProcessor) ProcessTransactions(
	header data.HeaderHandler,
	headerHash []byte,
	bodyHandler data.BodyHandler,
	transactions map[string]data.TransactionHandler) ([]*schema.Transaction, error) {

	body, ok := bodyHandler.(*erdBlock.Body)
	if !ok {
		return nil, covalent.ErrBlockBodyAssertion
	}

	allTxs := make([]*schema.Transaction, 0)
	for _, currMiniBlock := range body.MiniBlocks {
		if currMiniBlock.Type != block.TxBlock {
			continue
		}

		txsInCurrMB, err := txp.processTxsFromMiniBlock(transactions, currMiniBlock, header, headerHash)
		if err != nil {
			return nil, err
		}
		allTxs = append(allTxs, txsInCurrMB...)
	}

	return allTxs, nil
}

func (txp *transactionProcessor) processTxsFromMiniBlock(
	transactions map[string]data.TransactionHandler,
	miniBlock *erdBlock.MiniBlock,
	header data.HeaderHandler,
	blockHash []byte) ([]*schema.Transaction, error) {

	miniBlockHash, err := core.CalculateHash(txp.marshaller, txp.hasher, miniBlock)
	if err != nil {
		return nil, err
	}

	txsInMiniBlock := make([]*schema.Transaction, 0)
	for _, txHash := range miniBlock.TxHashes {
		tx, found := findTransactionInPool(txHash, transactions)
		if !found {
			log.Warn("transactionProcessor.processTxsFromMiniBlock tx hash not found in tx pool", "hash", txHash)
			continue
		}

		convertedTx := txp.convertTransaction(tx, txHash, miniBlockHash, blockHash, miniBlock, header)
		txsInMiniBlock = append(txsInMiniBlock, convertedTx)
	}

	return txsInMiniBlock, nil
}

func findTransactionInPool(txHash []byte, transactions map[string]data.TransactionHandler) (*transaction.Transaction, bool) {
	tx, isInTxPool := transactions[string(txHash)]
	if !isInTxPool {
		return nil, false
	}

	castedTx, castOk := tx.(*transaction.Transaction)
	if !castOk {
		return nil, false
	}

	return castedTx, true
}

func (txp *transactionProcessor) convertTransaction(
	tx *transaction.Transaction,
	txHash []byte,
	miniBlockHash []byte,
	blockHash []byte,
	miniBlock *erdBlock.MiniBlock,
	header data.HeaderHandler) *schema.Transaction {

	return &schema.Transaction{
		Hash:             txHash,
		MiniBlockHash:    miniBlockHash,
		BlockHash:        blockHash,
		Nonce:            int64(tx.GetNonce()),
		Round:            int64(header.GetRound()),
		Value:            utility.GetBytes(tx.GetValue()),
		Receiver:         utility.EncodePubKey(txp.pubKeyConverter, tx.GetRcvAddr()),
		Sender:           utility.EncodePubKey(txp.pubKeyConverter, tx.GetSndAddr()),
		ReceiverShard:    int32(miniBlock.ReceiverShardID),
		SenderShard:      int32(miniBlock.SenderShardID),
		GasPrice:         int64(tx.GetGasPrice()),
		GasLimit:         int64(tx.GetGasLimit()),
		Signature:        tx.GetSignature(),
		Timestamp:        int64(header.GetTimeStamp()),
		SenderUserName:   tx.GetSndUserName(),
		ReceiverUserName: tx.GetRcvUserName(),
	}
}
