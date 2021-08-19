package transactions

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/smartContractResult"
)

type scProcessor struct {
	pubKeyConverter core.PubkeyConverter
}

// NewSCProcessor creates a new instance of smart contracts processor
func NewSCProcessor(pubKeyConverter core.PubkeyConverter) (*scProcessor, error) {
	if check.IfNil(pubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}

	return &scProcessor{pubKeyConverter: pubKeyConverter}, nil
}

// ProcessSCs converts smart contracts data to a specific structure defined by avro schema
func (scp *scProcessor) ProcessSCs(transactions map[string]data.TransactionHandler, timeStamp uint64) ([]*schema.SCResult, error) {
	allScrTxs := make([]*schema.SCResult, 0)
	for currHash, currTx := range transactions {
		scrTx := currTx.(*smartContractResult.SmartContractResult)
		scrCovalent := &schema.SCResult{
			Hash:           []byte(currHash),
			Nonce:          int64(scrTx.GetNonce()),
			GasLimit:       int64(scrTx.GetGasLimit()),
			GasPrice:       int64(scrTx.GetGasPrice()),
			Value:          utility.GetBytes(scrTx.GetValue()),
			Sender:         utility.EncodePubKey(scp.pubKeyConverter, scrTx.GetSndAddr()),
			Receiver:       utility.EncodePubKey(scp.pubKeyConverter, scrTx.GetRcvAddr()),
			RelayerAddr:    utility.EncodePubKey(scp.pubKeyConverter, scrTx.GetRelayerAddr()),
			RelayedValue:   utility.GetBytes(scrTx.GetRelayedValue()),
			Code:           scrTx.GetCode(),
			Data:           scrTx.GetData(),
			PrevTxHash:     scrTx.GetPrevTxHash(),
			OriginalTxHash: scrTx.GetOriginalTxHash(),
			CallType:       int32(scrTx.GetCallType()),
			CodeMetadata:   scrTx.GetCodeMetadata(),
			ReturnMessage:  scrTx.GetReturnMessage(),
			Timestamp:      int64(timeStamp),
		}

		allScrTxs = append(allScrTxs, scrCovalent)
	}
	return allScrTxs, nil
}
