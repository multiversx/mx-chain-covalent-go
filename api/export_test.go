package api

import "github.com/gin-gonic/gin"

var ErrInvalidBlockHash = errInvalidBlockHash

var ErrNilHyperBlockFacade = errNilHyperBlockFacade

var ErrInvalidBlockNonce = errInvalidBlockNonce

var ErrInvalidHyperBlocksBatchSize = errInvalidHyperBlocksBatchSize

var ErrMissingQueryParameter = errMissingQueryParameter

func GetNonceFromRequest(c *gin.Context) (uint64, error) {
	return getNonceFromRequest(c)
}
