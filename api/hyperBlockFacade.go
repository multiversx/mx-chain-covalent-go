package api

import (
	"fmt"
	"net/url"

	"github.com/ElrondNetwork/covalent-indexer-go"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/api/shared"
)

const hyperBlockPathByNonce = "/hyperblock/by-nonce"
const hyperBlockPathByHash = "/hyperblock/by-hash"

var log = logger.GetOrCreate("api")

type hyperBlockFacade struct {
	elrondProxyUrl string
	processor      covalent.HyperBlockProcessor
	elrondEndpoint ElrondHyperBlockEndpointHandler
	encoder        AvroEncoder
}

// NewHyperBlockFacade will create a hyper block facade, which can fetch hyper blocks from Elrond proxy
func NewHyperBlockFacade(
	elrondProxyUrl string,
	avroEncoder AvroEncoder,
	elrondHyperBlockEndpoint ElrondHyperBlockEndpointHandler,
	hyperBlockProcessor covalent.HyperBlockProcessor,
) (*hyperBlockFacade, error) {
	if len(elrondProxyUrl) == 0 {
		return nil, errEmptyElrondProxyUrl
	}
	if avroEncoder == nil {
		return nil, errNilAvroEncoder
	}
	if hyperBlockProcessor == nil {
		return nil, errNilHyperBlockProcessor
	}
	if elrondHyperBlockEndpoint == nil {
		return nil, errNilHyperBlockEndpointHandler
	}

	return &hyperBlockFacade{
		elrondProxyUrl: elrondProxyUrl,
		processor:      hyperBlockProcessor,
		encoder:        avroEncoder,
		elrondEndpoint: elrondHyperBlockEndpoint,
	}, nil
}

// GetHyperBlockByNonce will fetch the hyper block from Elrond proxy with provided nonce and options in covalent format
func (hpf *hyperBlockFacade) GetHyperBlockByNonce(nonce uint64, options HyperBlockQueryOptions) (*CovalentHyperBlockApiResponse, error) {
	blockByNoncePath := fmt.Sprintf("%s/%d", hyperBlockPathByNonce, nonce)
	fullPath := hpf.getFullPathWithOptions(blockByNoncePath, options)

	return hpf.getHyperBlock(fullPath)
}

// GetHyperBlockByHash will fetch the hyper block from Elrond proxy with provided hash and options in covalent format
func (hpf *hyperBlockFacade) GetHyperBlockByHash(hash string, options HyperBlockQueryOptions) (*CovalentHyperBlockApiResponse, error) {
	blockByHashPath := fmt.Sprintf("%s/%s", hyperBlockPathByHash, hash)
	fullPath := hpf.getFullPathWithOptions(blockByHashPath, options)

	return hpf.getHyperBlock(fullPath)
}

func (hpf *hyperBlockFacade) getFullPathWithOptions(path string, options HyperBlockQueryOptions) string {
	pathWithOptions := buildUrlWithBlockQueryOptions(path, options)
	return fmt.Sprintf("%s%s", hpf.elrondProxyUrl, pathWithOptions)
}

func buildUrlWithBlockQueryOptions(path string, options HyperBlockQueryOptions) string {
	u := url.URL{Path: path}
	query := u.Query()

	setQueryParamIfTrue(query, options.WithLogs, UrlParameterWithLogs)
	setQueryParamIfTrue(query, options.WithBalances, UrlParameterWithBalances)

	u.RawQuery = query.Encode()
	return u.String()
}

func setQueryParamIfTrue(query url.Values, option bool, urlParam string) {
	if option {
		query.Set(urlParam, "true")
	}
}

func (hpf *hyperBlockFacade) getHyperBlock(path string) (*CovalentHyperBlockApiResponse, error) {
	elrondHyperBlock, err := hpf.elrondEndpoint.GetHyperBlock(path)
	if err != nil {
		return nil, err
	}

	hyperBlockSchema, err := hpf.processor.Process(&elrondHyperBlock.Data.HyperBlock)
	if err != nil {
		return nil, err
	}

	hyperBlockSchemaAvroBytes, err := hpf.encoder.Encode(hyperBlockSchema)
	if err != nil {
		return nil, err
	}

	return &CovalentHyperBlockApiResponse{
		Data:  hyperBlockSchemaAvroBytes,
		Error: "",
		Code:  shared.ReturnCodeSuccess,
	}, nil
}
