package api

import (
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/elrond-go/api/shared"
	"github.com/gin-gonic/gin"
)

var options = HyperBlockQueryOptions{
	WithLogs:     false,
	WithBalances: false,
}

type hyperBlockProxy struct {
	hyperBlockFacade HyperBlockFacadeHandler
	processor        covalent.HyperBlockProcessor
	marshaller       covalent.AvroMarshaller
}

func NewHyperBlockProxy(
	hyperBlockFacade HyperBlockFacadeHandler,
	marshaller covalent.AvroMarshaller,
	hyperBlockProcessor covalent.HyperBlockProcessor,
) *hyperBlockProxy {
	return &hyperBlockProxy{
		hyperBlockFacade: hyperBlockFacade,
		marshaller:       marshaller,
		processor:        hyperBlockProcessor,
	}
}

func (hbp *hyperBlockProxy) GetHyperBlockByNonce(c *gin.Context) {
	nonce, err := getNonceFromRequest(c)
	if err != nil {
		RespondWithBadRequest(c, errors.New("cannot parse nonce").Error())
		return
	}

	hyperBlockApiResponse, err := hbp.hyperBlockFacade.GetHyperBlockByNonce(nonce, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
		return
	}

	hbp.processHyperBlock(c, &hyperBlockApiResponse.Data.HyperBlock)
}

func getNonceFromRequest(c *gin.Context) (uint64, error) {
	nonceStr := c.Param("nonce")
	if nonceStr == "" {
		return 0, errors.New("invalid block nonce parameter")
	}

	return strconv.ParseUint(nonceStr, 10, 64)
}

func (hbp *hyperBlockProxy) processHyperBlock(c *gin.Context, hyperBlock *hyperBlock.HyperBlock) {
	blockSchema, err := hbp.processor.Process(hyperBlock)
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
}

func (hbp *hyperBlockProxy) GetHyperBlockByHash(c *gin.Context) {
	hash, err := checkHashFromRequest(c)
	if err != nil {
		RespondWithBadRequest(c, errors.New("invalid block hash parameter").Error())
		return
	}

	hyperBlockApiResponse, err := hbp.hyperBlockFacade.GetHyperBlockByHash(hash, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
		return
	}

	hbp.processHyperBlock(c, &hyperBlockApiResponse.Data.HyperBlock)
}

func checkHashFromRequest(c *gin.Context) (string, error) {
	hash := c.Param("hash")
	_, err := hex.DecodeString(hash)
	if err != nil {
		return "", errors.New("invalid block hash parameter")
	}

	return hash, nil
}

func RespondWithBadRequest(c *gin.Context, errorMessage string) {
	c.JSON(
		http.StatusBadRequest,
		shared.GenericAPIResponse{
			Data:  nil,
			Error: errorMessage,
			Code:  shared.ReturnCodeRequestError,
		},
	)
}
