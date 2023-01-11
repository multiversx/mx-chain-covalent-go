package transactions

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-covalent-go/process"
	"github.com/multiversx/mx-chain-covalent-go/process/utility"
	"github.com/multiversx/mx-chain-covalent-go/schema"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
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
func (txp *transactionProcessor) ProcessTransactions(apiTransactions []*transaction.ApiTransactionResult) ([]*schema.Transaction, error) {
	allTxs := make([]*schema.Transaction, 0, len(apiTransactions))

	for _, apiTx := range apiTransactions {
		if apiTx == nil {
			continue
		}

		tx, err := txp.processTransaction(apiTx)
		if err != nil {
			return nil, err
		}

		allTxs = append(allTxs, tx)
	}

	return allTxs, nil
}

func (txp *transactionProcessor) processTransaction(apiTx *transaction.ApiTransactionResult) (*schema.Transaction, error) {
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
	esdtValues, err := utility.GetBigIntBytesSliceFromStringSlice(apiTx.ESDTValues)
	if err != nil {
		return nil, err
	}
	initiallyPaidFee, err := utility.GetBigIntBytesFromStr(apiTx.InitiallyPaidFee)
	if err != nil {
		return nil, err
	}

	log := txp.logProcessor.ProcessLog(apiTx.Logs)
	return &schema.Transaction{
		Type:                              apiTx.Type,
		ProcessingTypeOnSource:            apiTx.ProcessingTypeOnSource,
		ProcessingTypeOnDestination:       apiTx.ProcessingTypeOnDestination,
		Hash:                              txHash,
		Nonce:                             int64(apiTx.Nonce),
		Round:                             int64(apiTx.Round),
		Epoch:                             int32(apiTx.Epoch),
		Value:                             value,
		Receiver:                          []byte(apiTx.Receiver),
		Sender:                            utility.GetAddressOrMetachainAddr(apiTx.Sender),
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
		OriginalSender:                    utility.GetAddressOrMetachainAddr(apiTx.OriginalSender),
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
		Receipt:                           receiptOrNil(receipt),
		Log:                               logOrNil(log),
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

func receiptOrNil(receipt *schema.Receipt) *schema.Receipt {
	if receipt == nil {
		return nil
	}

	if isReceiptEmpty(receipt) {
		return nil
	}

	return receipt
}

func isReceiptEmpty(receipt *schema.Receipt) bool {
	return len(receipt.Sender) == 0 &&
		len(receipt.Data) == 0 &&
		len(receipt.TxHash) == 0 &&
		len(receipt.Value) == 0
}

func logOrNil(log *schema.Log) *schema.Log {
	if log == nil {
		return nil
	}

	if isLogEmpty(log) {
		return nil
	}

	return log
}

func isLogEmpty(log *schema.Log) bool {
	return len(log.Address) == 0 && len(log.Events) == 0
}
