package api

import (
	"encoding/hex"
	"fmt"
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

// GetHyperBlocksByNonce will fetch requested hyper blocks from start to end nonce
func (hbp *hyperBlockProxy) GetHyperBlocksByInterval(c *gin.Context) {
	nonceInterval, err := getIntervalFromRequest(c)
	if err != nil {
		respondWithBadRequest(c, err)
		return
	}
	fmt.Println(fmt.Sprintf("%v", nonceInterval))

	hyperBlockApiResponse, err := hbp.hyperBlockFacade.GetHyperBlocksByInterval(nonceInterval, hbp.options)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, hyperBlockApiResponse)
}

func getIntervalFromRequest(c *gin.Context) (*Interval, error) {
	startNonce, err := parseUintUrlParam(c, "startNonce")
	if err != nil {
		return nil, err
	}

	endNonce, err := parseUintUrlParam(c, "endNonce")
	if err != nil {
		return nil, err
	}

	return &Interval{
		Start: startNonce,
		End:   endNonce,
	}, nil
}

func parseUintUrlParam(c *gin.Context, name string) (uint64, error) {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return 0, nil
	}

	value, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return 0, err
	}

	return value, nil
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
