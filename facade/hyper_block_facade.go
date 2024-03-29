package facade

import (
	"fmt"
	"net/url"

	"github.com/multiversx/mx-chain-covalent-go"
	"github.com/multiversx/mx-chain-covalent-go/api"
	"github.com/multiversx/mx-chain-covalent-go/cmd/proxy/config"
	logger "github.com/multiversx/mx-chain-logger-go"
)

const hyperBlockPathByNonce = "/hyperblock/by-nonce"
const hyperBlockPathByHash = "/hyperblock/by-hash"

var log = logger.GetOrCreate("facade")

type hyperBlockFacade struct {
	multiversxProxyUrl string
	processor          covalent.HyperBlockProcessor
	multiversxEndpoint api.MultiversxHyperBlockEndpointHandler
	encoder            AvroEncoder
}

// NewHyperBlockFacade will create a hyper block facade, which can fetch hyper blocks from Multiversx proxy
func NewHyperBlockFacade(
	multiversxProxyUrl string,
	avroEncoder AvroEncoder,
	multiversxHyperBlockEndpoint api.MultiversxHyperBlockEndpointHandler,
	hyperBlockProcessor covalent.HyperBlockProcessor,
) (*hyperBlockFacade, error) {
	if len(multiversxProxyUrl) == 0 {
		return nil, errEmptyMultiversxProxyUrl
	}
	if avroEncoder == nil {
		return nil, errNilAvroEncoder
	}
	if hyperBlockProcessor == nil {
		return nil, errNilHyperBlockProcessor
	}
	if multiversxHyperBlockEndpoint == nil {
		return nil, errNilHyperBlockEndpointHandler
	}

	return &hyperBlockFacade{
		multiversxProxyUrl: multiversxProxyUrl,
		processor:          hyperBlockProcessor,
		encoder:            avroEncoder,
		multiversxEndpoint: multiversxHyperBlockEndpoint,
	}, nil
}

// GetHyperBlockByNonce will fetch the hyper block from Multiversx proxy with provided nonce and options in covalent format
func (hbf *hyperBlockFacade) GetHyperBlockByNonce(nonce uint64, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
	fullPath := hbf.getHyperBlockByNonceFullPath(nonce, options)
	return hbf.getHyperBlock(fullPath)
}

// GetHyperBlocksByInterval will fetch the hyper blocks from Multiversx proxy with provided nonces interval and options in covalent format
func (hbf *hyperBlockFacade) GetHyperBlocksByInterval(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) (*api.CovalentHyperBlocksApiResponse, error) {
	encodedHyperBlocks, err := hbf.getHyperBlocksByNonces(noncesInterval, options)
	if err != nil {
		return nil, err
	}

	return &api.CovalentHyperBlocksApiResponse{
		Data:  encodedHyperBlocks,
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}, nil
}

func (hbf *hyperBlockFacade) getHyperBlockByNonceFullPath(nonce uint64, options config.HyperBlockQueryOptions) string {
	blockByNoncePath := fmt.Sprintf("%s/%d", hyperBlockPathByNonce, nonce)
	return hbf.getFullPathWithOptions(blockByNoncePath, options)
}

func (hbf *hyperBlockFacade) getFullPathWithOptions(path string, options config.HyperBlockQueryOptions) string {
	pathWithOptions := buildUrlWithBlockQueryOptions(path, options)
	return fmt.Sprintf("%s%s", hbf.multiversxProxyUrl, pathWithOptions)
}

func buildUrlWithBlockQueryOptions(path string, options config.HyperBlockQueryOptions) string {
	u := url.URL{Path: path}
	query := u.Query()

	setQueryParamIfTrue(query, options.WithLogs, api.UrlParameterWithLogs)
	setQueryParamIfTrue(query, options.NotarizedAtSource, api.UrlParameterNotarizedAtSource)
	setQueryParamIfTrue(query, options.WithAlteredAccounts, api.UrlParameterWithAlteredAccounts)
	setQueryParamIfNotEmpty(query, options.Tokens, api.UrlParameterTokens)

	u.RawQuery = query.Encode()
	return u.String()
}

func setQueryParamIfTrue(query url.Values, option bool, urlParam string) {
	if option {
		query.Set(urlParam, "true")
	}
}

func setQueryParamIfNotEmpty(query url.Values, option string, urlParam string) {
	if len(option) > 0 {
		query.Set(urlParam, option)
	}
}

func (hbf *hyperBlockFacade) getHyperBlockAvroBytes(path string) ([]byte, error) {
	multiversxHyperBlock, err := hbf.multiversxEndpoint.GetHyperBlock(path)
	if err != nil {
		return nil, err
	}

	hyperBlockSchema, err := hbf.processor.Process(&multiversxHyperBlock.Data.HyperBlock)
	if err != nil {
		return nil, err
	}

	return hbf.encoder.Encode(hyperBlockSchema)
}

func (hbf *hyperBlockFacade) getHyperBlock(path string) (*api.CovalentHyperBlockApiResponse, error) {
	hyperBlockSchemaAvroBytes, err := hbf.getHyperBlockAvroBytes(path)
	if err != nil {
		return nil, err
	}

	return &api.CovalentHyperBlockApiResponse{
		Data:  hyperBlockSchemaAvroBytes,
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}, nil
}

// GetHyperBlockByHash will fetch the hyper block from Multiversx proxy with provided hash and options in covalent format
func (hbf *hyperBlockFacade) GetHyperBlockByHash(hash string, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
	blockByHashPath := fmt.Sprintf("%s/%s", hyperBlockPathByHash, hash)
	fullPath := hbf.getFullPathWithOptions(blockByHashPath, options)

	return hbf.getHyperBlock(fullPath)
}
