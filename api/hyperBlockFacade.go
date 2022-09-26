package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const hyperBlockPathByNonce = "/hyperblock/by-nonce"
const hyperBlockPathByHash = "/hyperblock/by-hash"

var log = logger.GetOrCreate("api")

type HyperBlockFacade struct {
	elrondProxyUrl string
	httpClient     HTTPClient
}

// NewHyperBlockFacade will create a hyper block facade, which can fetch hyper blocks from Elrond proxy
func NewHyperBlockFacade(httpClient HTTPClient, elrondProxyUrl string) *HyperBlockFacade {
	return &HyperBlockFacade{
		httpClient:     httpClient,
		elrondProxyUrl: elrondProxyUrl,
	}
}

// GetHyperBlockByNonce will fetch the hyper block from Elrond proxy with provided nonce and options
func (hpf *HyperBlockFacade) GetHyperBlockByNonce(nonce uint64, options HyperBlockQueryOptions) (*HyperBlockApiResponse, error) {
	blockByNoncePath := fmt.Sprintf("%s/%d", hyperBlockPathByNonce, nonce)
	fullPath := hpf.getFullPathWithOptions(blockByNoncePath, options)

	return hpf.getHyperBlock(fullPath)
}

// GetHyperBlockByHash will fetch the hyper block from Elrond proxy with provided hash and options
func (hpf *HyperBlockFacade) GetHyperBlockByHash(hash string, options HyperBlockQueryOptions) (*HyperBlockApiResponse, error) {
	blockByHashPath := fmt.Sprintf("%s/%s", hyperBlockPathByHash, hash)
	fullPath := hpf.getFullPathWithOptions(blockByHashPath, options)

	return hpf.getHyperBlock(fullPath)
}

func (hpf *HyperBlockFacade) getFullPathWithOptions(path string, options HyperBlockQueryOptions) string {
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

func (hpf *HyperBlockFacade) getHyperBlock(path string) (*HyperBlockApiResponse, error) {
	resp, err := hpf.httpClient.Get(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		errNotCritical := resp.Body.Close()
		if errNotCritical != nil {
			log.Warn("close body", "error", errNotCritical.Error())
		}
	}()

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response HyperBlockApiResponse
	err = json.Unmarshal(responseBodyBytes, &response)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, elrond proxy response error: %s", resp.StatusCode, response.Error)
	}

	return &response, nil
}
