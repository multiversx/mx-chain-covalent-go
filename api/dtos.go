package api

import (
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/elrond-go/api/shared"
)

// ElrondHyperBlockApiResponse is the expected hyper block dto response from Elrond proxy
type ElrondHyperBlockApiResponse struct {
	Data  ElrondHyperBlockApiResponsePayload `json:"data"`
	Error string                             `json:"error"`
	Code  shared.ReturnCode                  `json:"code"`
}

// ElrondHyperBlockApiResponsePayload wraps a hyperBlock
type ElrondHyperBlockApiResponsePayload struct {
	HyperBlock hyperBlock.HyperBlock `json:"hyperblock"`
}

// CovalentHyperBlockApiResponse is the hyper block dto response for Covalent
type CovalentHyperBlockApiResponse struct {
	Data  []byte            `json:"data"`
	Error string            `json:"error"`
	Code  shared.ReturnCode `json:"code"`
}

// TODO: Remove this from here and use the one from core, once feat/altered-account from elrond-go & elrond-go-core are finished

// AlteredAccount is the altered account dto response from Elrond proxy
type AlteredAccount struct {
	Address string `json:"address"`
	Balance string `json:"balance,omitempty"`
	Nonce   uint64 `json:"nonce"`
}
