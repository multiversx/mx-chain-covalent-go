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

type hyperBlockProxy struct {
	hyperBlockFacade HyperBlockFacadeHandler
	processor        covalent.HyperBlockProcessor
	encoder          AvroEncoder
}

// NewHyperBlockProxy will create a covalent proxy, able to process hyper block requests from Elrond
func NewHyperBlockProxy(
	hyperBlockFacade HyperBlockFacadeHandler,
	avroEncoder AvroEncoder,
	hyperBlockProcessor covalent.HyperBlockProcessor,
) *hyperBlockProxy {
	return &hyperBlockProxy{
		hyperBlockFacade: hyperBlockFacade,
		encoder:          avroEncoder,
		processor:        hyperBlockProcessor,
	}
}

// GetHyperBlockByNonce will process given hyper block request by nonce from Elrond
func (hbp *hyperBlockProxy) GetHyperBlockByNonce(c *gin.Context) {
	nonce, err := getNonceFromRequest(c)
	if err != nil {
		respondWithBadRequest(c, errors.New("cannot parse nonce").Error())
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

	blockSchemaAvroBytes, err := hbp.encoder.Encode(blockSchema)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
		return
	}

	shared.RespondWithSuccess(c, blockSchemaAvroBytes)
}

// GetHyperBlockByHash will process given hyper block request by hash from Elrond
func (hbp *hyperBlockProxy) GetHyperBlockByHash(c *gin.Context) {
	hash, err := getHashFromRequest(c)
	if err != nil {
		respondWithBadRequest(c, errors.New("invalid block hash parameter").Error())
		return
	}

	hyperBlockApiResponse, err := hbp.hyperBlockFacade.GetHyperBlockByHash(hash, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), shared.ReturnCodeInternalError)
		return
	}

	hbp.processHyperBlock(c, &hyperBlockApiResponse.Data.HyperBlock)
}

func getHashFromRequest(c *gin.Context) (string, error) {
	hash := c.Param("hash")
	_, err := hex.DecodeString(hash)
	if err != nil {
		return "", errors.New("invalid block hash parameter")
	}

	return hash, nil
}

func respondWithBadRequest(c *gin.Context, errorMessage string) {
	c.JSON(
		http.StatusBadRequest,
		shared.GenericAPIResponse{
			Data:  nil,
			Error: errorMessage,
			Code:  shared.ReturnCodeRequestError,
		},
	)
}
