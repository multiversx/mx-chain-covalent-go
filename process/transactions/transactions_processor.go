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
			myTx, found := transactions[string(txHash)]
			if !found {
				//something is wrong
			}
			convertedTx := myTx.(*transaction.Transaction)

			tx := &schema.Transaction{
				Hash:             txHash,
				MiniBlockHash:    mbHash,
				BlockHash:        headerHash,
				Nonce:            int64(convertedTx.GetNonce()),
				Round:            int64(header.GetRound()), //from header
				Value:            utility.GetBytes(convertedTx.GetValue()),
				Receiver:         convertedTx.GetRcvAddr(),
				Sender:           convertedTx.GetSndAddr(),
				ReceiverShard:    int32(mb.ReceiverShardID),
				SenderShard:      int32(mb.SenderShardID),
				GasPrice:         int64(convertedTx.GetGasPrice()),
				GasLimit:         int64(convertedTx.GetGasLimit()),
				Signature:        convertedTx.GetSignature(),
				Timestamp:        int64(header.GetTimeStamp()), // from header
				SenderUserName:   convertedTx.GetSndUserName(),
				ReceiverUserName: convertedTx.GetRcvAddr(),
			}
			ret = append(ret, tx)
			//convert transaction
		}
	}

	return ret, nil
}
