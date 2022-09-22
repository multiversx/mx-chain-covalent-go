package api

import (
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-go/api/shared"
	"github.com/gin-gonic/gin"
)

// ReturnCode defines the type defines to identify return codes
type ReturnCode string

const (
	// ReturnCodeSuccess defines a successful request
	ReturnCodeSuccess ReturnCode = "successful"

	// ReturnCodeInternalError defines a request which hasn't been executed successfully due to an internal error
	ReturnCodeInternalError ReturnCode = "internal_issue"

	// ReturnCodeRequestError defines a request which hasn't been executed successfully due to a bad request received
	ReturnCodeRequestError ReturnCode = "bad_request"
)

type hyperBlockProxy struct {
	hyperBlockFacade HyperBlockFacadeHandler
}

func NewHyperBlockProxy(hyperBlockFacade HyperBlockFacadeHandler) *hyperBlockProxy {
	return &hyperBlockProxy{
		hyperBlockFacade: hyperBlockFacade,
	}
}

func (hbp *hyperBlockProxy) GetHyperBlockByNonce(c *gin.Context) {
	nonce, err := FetchNonceFromRequest(c)
	if err != nil {
		RespondWithBadRequest(c, errors.New("cannot parse nonce").Error())
		return
	}

	options, err := parseHyperblockQueryOptions(c)
	if err != nil {
		shared.RespondWithValidationError(c, errors.New("bad url parameter(s)"), err)
		return
	}

	blockByNonceResponse, err := hbp.hyperBlockFacade.GetHyperBlockByNonce(nonce, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByNonceResponse)

}

// FetchNonceFromRequest will try to fetch the nonce from the request
func FetchNonceFromRequest(c *gin.Context) (uint64, error) {
	nonceStr := c.Param("nonce")
	if nonceStr == "" {
		return 0, errors.New("invalid block nonce parameter")
	}

	return strconv.ParseUint(nonceStr, 10, 64)
}

func (hbp *hyperBlockProxy) GetHyperBlockByHash(c *gin.Context) {
	hash := c.Param("hash")
	_, err := hex.DecodeString(hash)
	if err != nil {
		RespondWithBadRequest(c, errors.New("invalid block hash parameter").Error())
		return
	}

	options, err := parseHyperblockQueryOptions(c)
	if err != nil {
		shared.RespondWithValidationError(c, errors.New("bad url parameter(s)"), err)
		return
	}

	blockByHashResponse, err := hbp.hyperBlockFacade.GetHyperBlockByHash(hash, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
}

func parseHyperblockQueryOptions(c *gin.Context) (HyperBlockQueryOptions, error) {
	withLogs, err := parseBoolUrlParam(c, UrlParameterWithLogs)
	if err != nil {
		return HyperBlockQueryOptions{}, err
	}

	options := HyperBlockQueryOptions{WithLogs: withLogs}
	return options, nil
}

func parseBoolUrlParam(c *gin.Context, name string) (bool, error) {
	return parseBoolUrlParamWithDefault(c, name, false)
}

func parseBoolUrlParamWithDefault(c *gin.Context, name string, defaultValue bool) (bool, error) {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return defaultValue, nil
	}

	return strconv.ParseBool(param)
}

// RespondWithBadRequest creates a generic response for bad request
func RespondWithBadRequest(c *gin.Context, errorMessage string) {
	RespondWith(c, http.StatusBadRequest, nil, errorMessage, ReturnCodeRequestError)
}

// GenericAPIResponse defines the structure of all responses on API endpoints
type GenericAPIResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
	Code  ReturnCode  `json:"code"`
}

// RespondWith will respond with the generic API response
func RespondWith(c *gin.Context, status int, dataField interface{}, error string, code ReturnCode) {
	c.JSON(
		status,
		GenericAPIResponse{
			Data:  dataField,
			Error: error,
			Code:  code,
		},
	)
}

const (
	// UrlParameterWithLogs represents the name of an URL parameter
	UrlParameterWithLogs = "withLogs"
)
