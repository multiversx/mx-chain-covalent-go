package api

import (
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
)

// ReturnCode defines the type defines to identify return codes
type ReturnCode string

// ReturnCodeSuccess defines a successful request
const ReturnCodeSuccess ReturnCode = "successful"

// ReturnCodeInternalError defines a request which hasn't been executed successfully due to an internal error
const ReturnCodeInternalError ReturnCode = "internal_issue"

// ReturnCodeRequestError defines a request which hasn't been executed successfully due to a bad request received
const ReturnCodeRequestError ReturnCode = "bad_request"

// ElrondHyperBlockApiResponse is the expected hyper block dto response from Elrond proxy
type ElrondHyperBlockApiResponse struct {
	Data  ElrondHyperBlockApiResponsePayload `json:"data"`
	Error string                             `json:"error"`
	Code  ReturnCode                         `json:"code"`
}

// ElrondHyperBlockApiResponsePayload wraps a hyperBlock
type ElrondHyperBlockApiResponsePayload struct {
	HyperBlock hyperBlock.HyperBlock `json:"hyperblock"`
}

// CovalentHyperBlockApiResponse is the hyper block dto response for Covalent
type CovalentHyperBlockApiResponse struct {
	Data  []byte     `json:"data"`
	Error string     `json:"error"`
	Code  ReturnCode `json:"code"`
}
