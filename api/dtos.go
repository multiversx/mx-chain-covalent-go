package api

import (
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/elrond-go/api/shared"
)

// HyperBlockApiResponse is a response holding a block
type HyperBlockApiResponse struct {
	Data  HyperBlockApiResponsePayload `json:"data"`
	Error string                       `json:"error"`
	Code  shared.ReturnCode            `json:"code"`
}

// HyperBlockApiResponsePayload wraps a hyperBlock
type HyperBlockApiResponsePayload struct {
	HyperBlock hyperBlock.HyperBlock `json:"hyperblock"`
}
