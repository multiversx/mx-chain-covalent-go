package transactions

import (
	"encoding/hex"

	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

type transactionProcessor struct {
	logProcessor   process.LogHandler
	receiptHandler process.ReceiptHandler
}

// NewTransactionProcessor creates a new instance of transactions processor
func NewTransactionProcessor(
	logProcessor process.LogHandler,
	receiptHandler process.ReceiptHandler,
) (*transactionProcessor, error) {
	if logProcessor == nil {
		return nil, errNilLogProcessor
	}
	if receiptHandler == nil {
		return nil, errNilReceiptProcessor
	}

	return &transactionProcessor{
		logProcessor:   logProcessor,
		receiptHandler: receiptHandler,
	}, nil
}

// ProcessTransactions converts transactions data to a specific structure defined by avro schema
func (txp *transactionProcessor) ProcessTransactions(apiTransactions []*transaction.ApiTransactionResult) ([]*schemaV2.Transaction, error) {
	allTxs := make([]*schemaV2.Transaction, 0, len(apiTransactions))

	for _, apiTx := range apiTransactions {
		tx, err := txp.processTransaction(apiTx)
		if err != nil {
			return nil, err
		}

		allTxs = append(allTxs, tx)
	}

	return allTxs, nil
}

func (txp *transactionProcessor) processTransaction(apiTx *transaction.ApiTransactionResult) (*schemaV2.Transaction, error) {
	txHash, err := hex.DecodeString(apiTx.Hash)
	if err != nil {
		return nil, err
	}
	value, err := utility.GetBigIntBytesFromStr(apiTx.Value)
	if err != nil {
		return nil, err
	}
	prevTxHash, err := hex.DecodeString(apiTx.PreviousTransactionHash)
	if err != nil {
		return nil, err
	}
	originalTxHash, err := hex.DecodeString(apiTx.OriginalTransactionHash)
	if err != nil {
		return nil, err
	}
	signature, err := hex.DecodeString(apiTx.Signature)
	if err != nil {
		return nil, err
	}
	blockHash, err := hex.DecodeString(apiTx.BlockHash)
	if err != nil {
		return nil, err
	}
	notarizedAtSourceInMetaHash, err := hex.DecodeString(apiTx.NotarizedAtSourceInMetaHash)
	if err != nil {
		return nil, err
	}
	notarizedAtDestinationInMetaHash, err := hex.DecodeString(apiTx.NotarizedAtDestinationInMetaHash)
	if err != nil {
		return nil, err
	}
	miniBlockHash, err := hex.DecodeString(apiTx.MiniBlockHash)
	if err != nil {
		return nil, err
	}
	hyperBlockHash, err := hex.DecodeString(apiTx.HyperblockHash)
	if err != nil {
		return nil, err
	}
	receipt, err := txp.receiptHandler.ProcessReceipt(apiTx.Receipt)
	if err != nil {
		return nil, err
	}
	esdtValues, err := utility.BigIntBytesSliceFromStringSlice(apiTx.ESDTValues)
	if err != nil {
		return nil, err
	}
	initiallyPaidFee, err := utility.GetBigIntBytesFromStr(apiTx.InitiallyPaidFee)
	if err != nil {
		return nil, err
	}

	return &schemaV2.Transaction{
		Type:                              apiTx.Type,
		ProcessingTypeOnSource:            apiTx.ProcessingTypeOnSource,
		ProcessingTypeOnDestination:       apiTx.ProcessingTypeOnDestination,
		Hash:                              txHash,
		Nonce:                             int64(apiTx.Nonce),
		Round:                             int64(apiTx.Round),
		Epoch:                             int32(apiTx.Epoch),
		Value:                             value,
		Receiver:                          []byte(apiTx.Receiver),
		Sender:                            []byte(apiTx.Sender),
		SenderUserName:                    apiTx.SenderUsername,
		ReceiverUserName:                  apiTx.ReceiverUsername,
		GasPrice:                          int64(apiTx.GasPrice),
		GasLimit:                          int64(apiTx.GasLimit),
		Data:                              apiTx.Data,
		CodeMetadata:                      apiTx.CodeMetadata,
		Code:                              []byte(apiTx.Code),
		PreviousTransactionHash:           prevTxHash,
		OriginalTransactionHash:           originalTxHash,
		ReturnMessage:                     apiTx.ReturnMessage,
		OriginalSender:                    []byte(apiTx.OriginalSender),
		Signature:                         signature,
		SourceShard:                       int32(apiTx.SourceShard),
		DestinationShard:                  int32(apiTx.DestinationShard),
		BlockNonce:                        int64(apiTx.BlockNonce),
		BlockHash:                         blockHash,
		NotarizedAtSourceInMetaNonce:      int64(apiTx.NotarizedAtSourceInMetaNonce),
		NotarizedAtSourceInMetaHash:       notarizedAtSourceInMetaHash,
		NotarizedAtDestinationInMetaNonce: int64(apiTx.NotarizedAtDestinationInMetaNonce),
		NotarizedAtDestinationInMetaHash:  notarizedAtDestinationInMetaHash,
		MiniBlockType:                     apiTx.MiniBlockType,
		MiniBlockHash:                     miniBlockHash,
		HyperBlockNonce:                   int64(apiTx.HyperblockNonce),
		HyperBlockHash:                    hyperBlockHash,
		Timestamp:                         apiTx.Timestamp,
		Receipt:                           receipt,
		Log:                               txp.logProcessor.ProcessLog(apiTx.Logs),
		Status:                            apiTx.Status.String(),
		Tokens:                            apiTx.Tokens,
		ESDTValues:                        esdtValues,
		Receivers:                         utility.StringSliceToByteSlice(apiTx.Receivers),
		ReceiversShardIDs:                 utility.UInt32SliceToInt32Slice(apiTx.ReceiversShardIDs),
		Operation:                         apiTx.Operation,
		Function:                          apiTx.Function,
		InitiallyPaidFee:                  initiallyPaidFee,
		IsRelayed:                         apiTx.IsRelayed,
		IsRefund:                          apiTx.IsRefund,
	}, nil
}
