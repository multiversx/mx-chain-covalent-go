package api

import (
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/covalent-indexer-go/cmd/proxy/config"
	"github.com/gin-gonic/gin"
)

type hyperBlockProxy struct {
	hyperBlockFacade HyperBlockFacadeHandler
	options          config.HyperBlockQueryOptions
}

// NewHyperBlockProxy will create a covalent proxy, able to fetch hyper block requests
// from Elrond and return them in covalent format
func NewHyperBlockProxy(
	hyperBlockFacade HyperBlockFacadeHandler,
	hyperBlockQueryOptions config.HyperBlockQueryOptions,
) (*hyperBlockProxy, error) {
	if hyperBlockFacade == nil {
		return nil, errNilHyperBlockFacade
	}

	return &hyperBlockProxy{
		hyperBlockFacade: hyperBlockFacade,
		options:          hyperBlockQueryOptions,
	}, nil
}

// GetHyperBlockByNonce will fetch requested hyper block request by nonce
func (hbp *hyperBlockProxy) GetHyperBlockByNonce(c *gin.Context) {
	nonce, err := getNonceFromRequest(c)
	if err != nil {
		respondWithBadRequest(c, err)
		return
	}

	hyperBlockApiResponse, err := hbp.hyperBlockFacade.GetHyperBlockByNonce(nonce, hbp.options)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, hyperBlockApiResponse)
}

func getNonceFromRequest(c *gin.Context) (uint64, error) {
	nonceStr := c.Param("nonce")
	if nonceStr == "" {
		return 0, errInvalidBlockNonce
	}

	return strconv.ParseUint(nonceStr, 10, 64)
}

// GetHyperBlockByHash will fetch requested hyper block request by hash
func (hbp *hyperBlockProxy) GetHyperBlockByHash(c *gin.Context) {
	hash, err := getHashFromRequest(c)
	if err != nil {
		respondWithBadRequest(c, err)
		return
	}

	hyperBlockApiResponse, err := hbp.hyperBlockFacade.GetHyperBlockByHash(hash, hbp.options)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, hyperBlockApiResponse)
}

func getHashFromRequest(c *gin.Context) (string, error) {
	hash := c.Param("hash")
	_, err := hex.DecodeString(hash)
	if err != nil {
		return "", errInvalidBlockHash
	}

	return hash, nil
}

func respondWithInternalError(c *gin.Context, err error) {
	c.JSON(
		http.StatusInternalServerError,
		CovalentHyperBlockApiResponse{
			Data:  nil,
			Error: err.Error(),
			Code:  ReturnCodeInternalError,
		},
	)
}

func respondWithBadRequest(c *gin.Context, err error) {
	c.JSON(
		http.StatusBadRequest,
		CovalentHyperBlockApiResponse{
			Data:  nil,
			Error: err.Error(),
			Code:  ReturnCodeRequestError,
		},
	)
}
