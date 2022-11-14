package facade

import (
	"fmt"
	"net/url"

	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/cmd/proxy/config"
)

const hyperBlockPathByNonce = "/hyperblock/by-nonce"
const hyperBlockPathByHash = "/hyperblock/by-hash"

type hyperBlockFacade struct {
	elrondProxyUrl string
	processor      covalent.HyperBlockProcessor
	elrondEndpoint api.ElrondHyperBlockEndpointHandler
	encoder        AvroEncoder
}

// NewHyperBlockFacade will create a hyper block facade, which can fetch hyper blocks from Elrond proxy
func NewHyperBlockFacade(
	elrondProxyUrl string,
	avroEncoder AvroEncoder,
	elrondHyperBlockEndpoint api.ElrondHyperBlockEndpointHandler,
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
func (hpf *hyperBlockFacade) GetHyperBlockByNonce(nonce uint64, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
	blockByNoncePath := fmt.Sprintf("%s/%d", hyperBlockPathByNonce, nonce)
	fullPath := hpf.getFullPathWithOptions(blockByNoncePath, options)

	return hpf.getHyperBlock(fullPath)
}

// GetHyperBlockByHash will fetch the hyper block from Elrond proxy with provided hash and options in covalent format
func (hpf *hyperBlockFacade) GetHyperBlockByHash(hash string, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
	blockByHashPath := fmt.Sprintf("%s/%s", hyperBlockPathByHash, hash)
	fullPath := hpf.getFullPathWithOptions(blockByHashPath, options)

	return hpf.getHyperBlock(fullPath)
}

func (hpf *hyperBlockFacade) getFullPathWithOptions(path string, options config.HyperBlockQueryOptions) string {
	pathWithOptions := buildUrlWithBlockQueryOptions(path, options)
	return fmt.Sprintf("%s%s", hpf.elrondProxyUrl, pathWithOptions)
}

func buildUrlWithBlockQueryOptions(path string, options config.HyperBlockQueryOptions) string {
	u := url.URL{Path: path}
	query := u.Query()

	setQueryParamIfTrue(query, options.WithLogs, api.UrlParameterWithLogs)
	setQueryParamIfTrue(query, options.NotarizedAtSource, api.UrlParameterNotarizedAtSource)
	setQueryParamIfTrue(query, options.WithAlteredAccounts, api.UrlParameterWithAlteredAccounts)
	setQueryParamIfNotEmpty(query, options.Tokens, api.UrlParameterTokens)
	setQueryParamIfTrue(query, options.WithMetaData, api.UrlParameterWithMetaData)

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

func (hpf *hyperBlockFacade) getHyperBlock(path string) (*api.CovalentHyperBlockApiResponse, error) {
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

	return &api.CovalentHyperBlockApiResponse{
		Data:  hyperBlockSchemaAvroBytes,
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}, nil
}
