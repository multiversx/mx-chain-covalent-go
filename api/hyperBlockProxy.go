package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/covalent-indexer-go"
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

var options = HyperBlockQueryOptions{
	WithLogs:     false,
	WithBalances: false,
}

type hyperBlockProxy struct {
	hyperBlockFacade HyperBlockFacadeHandler
	processor        covalent.DataHandler
	marshaller       covalent.AvroMarshaller
}

func NewHyperBlockProxy(
	hyperBlockFacade HyperBlockFacadeHandler,
	marshaller covalent.AvroMarshaller,
	blockProcessor covalent.DataHandler,
) *hyperBlockProxy {
	return &hyperBlockProxy{
		hyperBlockFacade: hyperBlockFacade,
		marshaller:       marshaller,
		processor:        blockProcessor,
	}
}

func (hbp *hyperBlockProxy) GetHyperBlockByNonce(c *gin.Context) {
	nonce, err := FetchNonceFromRequest(c)
	if err != nil {
		RespondWithBadRequest(c, errors.New("cannot parse nonce").Error())
		return
	}

	blockByNonceResponse, err := hbp.hyperBlockFacade.GetHyperBlockByNonce(nonce, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
		return
	}
	//hash, err := hex.DecodeString(blockByNonceResponse.Data.Hyperblock.Hash)
	//if err != nil {
	//	shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
	//	return
	//}
	_ = blockByNonceResponse
	x, _ := json.Marshal(blockByNonceResponse)
	log.Info("dsadsa", "dsada", string(x))

	//blockRes := &schema.BlockResult{
	//	Block: &schema.Block{
	//		Hash: hash,
	//	},
	//}
	//rrr, err := utility.Encode(blockRes)
	//rrr, err := json.Marshal(blockRes)

	//if err != nil {
	//	shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
	//	return
	//}
	blockSchema, err := hbp.processor.ProcessData(nil)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
		return
	}

	blockSchemaAvroBytes, err := hbp.marshaller.Encode(blockSchema)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
		return
	}

	shared.RespondWithSuccess(c, blockSchemaAvroBytes)
	//c.Data(http.StatusOK,"" ,rrr)
}

// FetchNonceFromRequest will try to fetch the nonce from the request
func FetchNonceFromRequest(c *gin.Context) (uint64, error) {
	nonceStr := c.Param("nonce")
	if nonceStr == "" {
		return 0, errors.New("invalid block nonce parameter")
	}

	return strconv.ParseUint(nonceStr, 10, 64)
}

func checkHashFromRequest(c *gin.Context) (string, error) {
	hash := c.Param("hash")
	_, err := hex.DecodeString(hash)
	if err != nil {
		return "", errors.New("invalid block hash parameter")
	}

	return hash, nil
}

func (hbp *hyperBlockProxy) GetHyperBlockByHash(c *gin.Context) {
	hash, err := checkHashFromRequest(c)
	if err != nil {
		RespondWithBadRequest(c, errors.New("invalid block hash parameter").Error())
		return
	}

	blockByHashResponse, err := hbp.hyperBlockFacade.GetHyperBlockByHash(hash, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
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
	// UrlParameterWithBalances represents the name of an URL parameter to query balances per hyperBlock
	UrlParameterWithBalances = "withBalances"
	// UrlParameterWithLogs represents the name of an URL parameter to query logs per hyperBlock
	UrlParameterWithLogs = "withLogs"
)
