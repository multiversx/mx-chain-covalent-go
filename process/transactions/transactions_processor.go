package transactions

import (
	"fmt"

	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	erdBlock "github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/elrond-go-core/data/rewardTx"
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
	marshaller marshal.Marshalizer,
) (*transactionProcessor, error) {
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
	pool *indexer.Pool,
) ([]*schema.Transaction, error) { //TODO: MAYBE NEVER RETURN ERROR???????????????????????????????????
	body, ok := bodyHandler.(*erdBlock.Body)
	if !ok {
		return nil, covalent.ErrBlockBodyAssertion
	}

	transactions := getTransactions(pool)
	allTxs := make([]*schema.Transaction, 0)
	for _, currMiniBlock := range body.MiniBlocks {
		var txsInCurrMB []*schema.Transaction
		var err error

		switch currMiniBlock.Type {
		case block.TxBlock:
			txsInCurrMB, err = txp.processTxsFromMiniBlock(transactions, currMiniBlock, header, headerHash, block.TxBlock)
		case block.RewardsBlock:
			txsInCurrMB, err = txp.processTxsFromMiniBlock(transactions, currMiniBlock, header, headerHash, block.RewardsBlock)
		default:
			continue
		}

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
	blockHash []byte,
	mbType block.Type,
) ([]*schema.Transaction, error) {

	miniBlockHash, err := core.CalculateHash(txp.marshaller, txp.hasher, miniBlock)
	if err != nil {
		return nil, err
	}

	txsInMiniBlock := make([]*schema.Transaction, 0, len(miniBlock.TxHashes))
	for _, txHash := range miniBlock.TxHashes {
		tx, isTxInPool := transactions[string(txHash)]
		if !isTxInPool {
			log.Warn("transactionProcessor.processTxsFromMiniBlock tx hash not found in tx pool", "hash", txHash)
			continue
		}

		convertedTx := txp.convertTransaction(tx, txHash, miniBlockHash, blockHash, miniBlock, header, mbType)
		if convertedTx != nil {
			txsInMiniBlock = append(txsInMiniBlock, convertedTx)
		}
	}

	return txsInMiniBlock, nil
}

func (txp *transactionProcessor) convertRewardTransaction(
	transaction data.TransactionHandler,
	txHash []byte,
	miniBlockHash []byte,
	blockHash []byte,
	miniBlock *erdBlock.MiniBlock,
	header data.HeaderHandler,
) *schema.Transaction {
	tx, castOk := transaction.(*rewardTx.RewardTx)
	if !castOk {
		return nil
	}
	//TODO: ADD DATAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
	return &schema.Transaction{
		Hash:             txHash,
		MiniBlockHash:    miniBlockHash,
		BlockHash:        blockHash,
		Nonce:            0,
		Round:            int64(tx.GetNonce()),
		Value:            utility.GetBytes(tx.GetValue()),
		Receiver:         utility.EncodePubKey(txp.pubKeyConverter, tx.GetRcvAddr()),
		Sender:           []byte(fmt.Sprintf("%d", core.MetachainShardId)),
		ReceiverShard:    int32(miniBlock.ReceiverShardID),
		SenderShard:      int32(miniBlock.SenderShardID),
		GasPrice:         0,
		GasLimit:         0,
		Signature:        make([]byte, 0),
		Timestamp:        int64(header.GetTimeStamp()),
		SenderUserName:   make([]byte, 0),
		ReceiverUserName: make([]byte, 0),
	}
}

func (txp *transactionProcessor) convertNormalTransaction(
	ntransaction data.TransactionHandler,
	txHash []byte,
	miniBlockHash []byte,
	blockHash []byte,
	miniBlock *erdBlock.MiniBlock,
	header data.HeaderHandler,
) *schema.Transaction {
	tx, castOk := ntransaction.(*transaction.Transaction)
	if !castOk {
		return nil
	}

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

func (txp *transactionProcessor) convertTransaction(
	tx data.TransactionHandler,
	txHash []byte,
	miniBlockHash []byte,
	blockHash []byte,
	miniBlock *erdBlock.MiniBlock,
	header data.HeaderHandler,
	mbType block.Type,
) *schema.Transaction {
	var ret *schema.Transaction

	switch mbType {
	case block.TxBlock:
		ret = txp.convertNormalTransaction(tx, txHash, miniBlockHash, blockHash, miniBlock, header)
	case block.RewardsBlock:
		ret = txp.convertRewardTransaction(tx, txHash, miniBlockHash, blockHash, miniBlock, header)
	default:
		return nil
	}

	return ret
}

func getTransactions(pool *indexer.Pool) map[string]data.TransactionHandler {
	transactions := make(map[string]data.TransactionHandler, len(pool.Txs)+len(pool.Rewards)+len(pool.Invalid))

	mergeTxsMaps(transactions, pool.Txs)
	mergeTxsMaps(transactions, pool.Rewards)
	mergeTxsMaps(transactions, pool.Invalid)

	return transactions
}

func mergeTxsMaps(dst, src map[string]data.TransactionHandler) {
	for key, value := range src {
		dst[key] = value
	}
}
