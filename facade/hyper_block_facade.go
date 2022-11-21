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
	fullPath := hpf.getHyperBlockByNonceFullPath(nonce, options)
	return hpf.getHyperBlock(fullPath)
}

// GetHyperBlocksByInterval will fetch the hyper blocks from Elrond proxy with provided nonces interval and options in covalent format
func (hpf *hyperBlockFacade) GetHyperBlocksByInterval(noncesInterval *api.Interval, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlocksApiResponse, error) {
	if noncesInterval.Start > noncesInterval.End {
		return nil, errInvalidNoncesInterval
	}

	// Dummy implementation with no parallel bulk requests. This implementation will follow in next PR
	encodedHyperBlocks := make([][]byte, 0, noncesInterval.End-noncesInterval.Start)
	for nonce := noncesInterval.Start; nonce <= noncesInterval.End; nonce++ {
		fullPath := hpf.getHyperBlockByNonceFullPath(nonce, options)
		encodedHyperBlock, err := hpf.getHyperBlockAvroBytes(fullPath)
		if err != nil {
			return nil, err
		}

		encodedHyperBlocks = append(encodedHyperBlocks, encodedHyperBlock)
	}

	return &api.CovalentHyperBlocksApiResponse{
		Data:  encodedHyperBlocks,
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}, nil
}

func (hpf *hyperBlockFacade) getHyperBlockByNonceFullPath(nonce uint64, options config.HyperBlockQueryOptions) string {
	blockByNoncePath := fmt.Sprintf("%s/%d", hyperBlockPathByNonce, nonce)
	return hpf.getFullPathWithOptions(blockByNoncePath, options)
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

func (hpf *hyperBlockFacade) getHyperBlockAvroBytes(path string) ([]byte, error) {
	elrondHyperBlock, err := hpf.elrondEndpoint.GetHyperBlock(path)
	if err != nil {
		return nil, err
	}

	hyperBlockSchema, err := hpf.processor.Process(&elrondHyperBlock.Data.HyperBlock)
	if err != nil {
		return nil, err
	}

	return hpf.encoder.Encode(hyperBlockSchema)
}

func (hpf *hyperBlockFacade) getHyperBlock(path string) (*api.CovalentHyperBlockApiResponse, error) {
	hyperBlockSchemaAvroBytes, err := hpf.getHyperBlockAvroBytes(path)
	if err != nil {
		return nil, err
	}

	return &api.CovalentHyperBlockApiResponse{
		Data:  hyperBlockSchemaAvroBytes,
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}, nil
}

// GetHyperBlockByHash will fetch the hyper block from Elrond proxy with provided hash and options in covalent format
func (hpf *hyperBlockFacade) GetHyperBlockByHash(hash string, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
	blockByHashPath := fmt.Sprintf("%s/%s", hyperBlockPathByHash, hash)
	fullPath := hpf.getFullPathWithOptions(blockByHashPath, options)

	return hpf.getHyperBlock(fullPath)
}
