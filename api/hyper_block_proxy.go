package api

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/covalent-indexer-go/cmd/proxy/config"
	"github.com/gin-gonic/gin"
)

const (
	startNonce = "startNonce"
	endNonce   = "endNonce"
)

type hyperBlockProxy struct {
	hyperBlockFacade HyperBlockFacadeHandler
	options          config.HyperBlockQueryOptions
	batchSize        uint32
}

// NewHyperBlockProxy will create a covalent proxy, able to fetch hyper block requests
// from Elrond and return them in covalent format
func NewHyperBlockProxy(
	hyperBlockFacade HyperBlockFacadeHandler,
	cfg config.Config,
) (*hyperBlockProxy, error) {
	if hyperBlockFacade == nil {
		return nil, errNilHyperBlockFacade
	}
	if cfg.HyperBlocksBatchSize == 0 {
		return nil, fmt.Errorf("%w; expected non zero value", errInvalidHyperBlocksBatchSize)
	}

	return &hyperBlockProxy{
		hyperBlockFacade: hyperBlockFacade,
		options:          cfg.HyperBlockQueryOptions,
		batchSize:        cfg.HyperBlocksBatchSize,
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

// GetHyperBlocksByInterval will fetch requested hyper blocks from start to end nonce
func (hbp *hyperBlockProxy) GetHyperBlocksByInterval(c *gin.Context) {
	noncesInterval, err := getIntervalFromRequest(c)
	if err != nil {
		respondWithBadRequest(c, err)
		return
	}

	options := config.HyperBlocksQueryOptions{
		QueryOptions: hbp.options,
		BatchSize:    hbp.batchSize,
	}
	hyperBlockApiResponse, err := hbp.hyperBlockFacade.GetHyperBlocksByInterval(noncesInterval, options)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, hyperBlockApiResponse)
}

func getIntervalFromRequest(c *gin.Context) (*Interval, error) {
	start, err := getUIntUrlParam(c, startNonce)
	if err != nil {
		return nil, err
	}

	end, err := getUIntUrlParam(c, endNonce)
	if err != nil {
		return nil, err
	}

	return &Interval{
		Start: start,
		End:   end,
	}, nil
}

func getUIntUrlParam(c *gin.Context, name string) (uint64, error) {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return 0, fmt.Errorf("%w: %s", errMissingQueryParameter, name)
	}

	return strconv.ParseUint(param, 10, 32)
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
