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
)

type transactionProcessor struct {
	hasher     hashing.Hasher
	marshaller marshal.Marshalizer
}

// NewTransactionProcessor creates a new instance of transactions processor
func NewTransactionProcessor(hasher hashing.Hasher, marshaller marshal.Marshalizer) (*transactionProcessor, error) {
	if check.IfNil(marshaller) {
		return nil, covalent.ErrNilMarshaller
	}
	if check.IfNil(hasher) {
		return nil, covalent.ErrNilHasher
	}

	return &transactionProcessor{
		hasher:     hasher,
		marshaller: marshaller,
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

	ret := make([]*schema.Transaction, 0)

	for _, mb := range body.MiniBlocks {
		if mb.Type != block.TxBlock {
			continue
		}

		mbHash, err := core.CalculateHash(txp.marshaller, txp.hasher, mb)
		if err != nil {
			return nil, err
		}

		for _, txHash := range mb.TxHashes {
			tx, found := txp.findTransactionInPool(txHash, transactions)
			if !found {
				continue
			}

			convertedTx := txp.convertTransaction(tx, txHash, mbHash, headerHash, mb, header)
			ret = append(ret, convertedTx)
		}
	}

	return ret, nil
}

func (txp *transactionProcessor) findTransactionInPool(
	txHash []byte,
	transactions map[string]data.TransactionHandler) (*transaction.Transaction, bool) {

	tx, found := transactions[string(txHash)]
	if !found {
		return nil, false
	}

	convertedTx, castOk := tx.(*transaction.Transaction)
	if !castOk {
		return nil, false
	}

	return convertedTx, true
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
		Receiver:         tx.GetRcvAddr(),
		Sender:           tx.GetSndAddr(),
		ReceiverShard:    int32(miniBlock.ReceiverShardID),
		SenderShard:      int32(miniBlock.SenderShardID),
		GasPrice:         int64(tx.GetGasPrice()),
		GasLimit:         int64(tx.GetGasLimit()),
		Signature:        tx.GetSignature(),
		Timestamp:        int64(header.GetTimeStamp()),
		SenderUserName:   tx.GetSndUserName(),
		ReceiverUserName: tx.GetRcvAddr(),
	}
}
