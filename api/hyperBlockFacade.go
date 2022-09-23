package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const hyperBlockPathByNonce = "/hyperblock/by-nonce"

var log = logger.GetOrCreate("process")

type HyperBlockFacade struct {
	path       string
	httpClient *http.Client
}

func NewHyperBlockFacade(requestTimeoutSec uint64, path string) *HyperBlockFacade {
	httpClient := http.DefaultClient

	var mutHttpClient sync.RWMutex
	mutHttpClient.Lock()
	httpClient.Timeout = time.Duration(requestTimeoutSec) * time.Second
	mutHttpClient.Unlock()

	return &HyperBlockFacade{
		httpClient: httpClient,
		path:       path,
	}
}

func (hpf *HyperBlockFacade) GetHyperBlockByNonce(nonce uint64, options HyperBlockQueryOptions) (*schema.BlockResult, error) {
	log.Info("path", "pp", fmt.Sprintf("%s%s/%d", hpf.path, hyperBlockPathByNonce, nonce))
	pppp := fmt.Sprintf("%s%s/%d", hpf.path, hyperBlockPathByNonce, nonce)
	path := BuildUrlWithBlockQueryOptions(fmt.Sprintf("%s%s/%d", hpf.path, hyperBlockPathByNonce, nonce), options)
	log.Info("path", "pp", path)
	var response HyperblockApiResponse

	_, err := hpf.httpGet(pppp, &response)
	if err != nil {
		return nil, err
	}

	log.Info("respons", "r", response.Data.Hyperblock.PrevBlockHash)

	return &schema.BlockResult{
		Block: &schema.Block{
			Nonce:     int64(nonce),
			Hash:      []byte(response.Data.Hyperblock.Hash),
			Timestamp: response.Data.Hyperblock.Timestamp.Milliseconds(),
		},
	}, nil
}
func (hpf *HyperBlockFacade) GetHyperBlockByHash(hash string, options HyperBlockQueryOptions) (*schema.BlockResult, error) {
	return nil, nil
}

func (hpf *HyperBlockFacade) httpGet(
	path string,
	value *HyperblockApiResponse,
) (int, error) {

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userAgent := "Elrond Proxy / 1.0.0 <Requesting data from nodes>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := hpf.httpClient.Do(req)
	if err != nil {
		return http.StatusNotFound, err
	}

	defer func() {
		errNotCritical := resp.Body.Close()
		if errNotCritical != nil {
			log.Warn("base process GET: close body", "error", errNotCritical.Error())
		}
	}()

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	log.Info("response body", "resp", string(responseBodyBytes))

	err = json.Unmarshal(responseBodyBytes, value)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	log.Info("resp val", "timestamp", value.Data.Hyperblock.Timestamp)

	responseStatusCode := resp.StatusCode
	if responseStatusCode == http.StatusOK { // everything ok, return status ok and the expected response
		return responseStatusCode, nil
	}

	// status response not ok, return the error
	return responseStatusCode, errors.New(string(responseBodyBytes))
}

// BuildUrlWithBlockQueryOptions builds an URL with block query parameters
func BuildUrlWithBlockQueryOptions(path string, options HyperBlockQueryOptions) string {
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

// HyperblockApiResponse is a response holding a block
type HyperblockApiResponse struct {
	Data  HyperblockApiResponsePayload `json:"data"`
	Error string                       `json:"error"`
	Code  ReturnCode                   `json:"code"`
}

// HyperblockApiResponsePayload wraps a hyperblock
type HyperblockApiResponsePayload struct {
	Hyperblock Hyperblock `json:"hyperblock"`
}

// Hyperblock contains all fully executed (both in source and in destination shards) transactions notarized in a given metablock
type Hyperblock struct {
	Hash                   string                              `json:"hash"`
	PrevBlockHash          string                              `json:"prevBlockHash"`
	StateRootHash          string                              `json:"stateRootHash"`
	Nonce                  uint64                              `json:"nonce"`
	Round                  uint64                              `json:"round"`
	Epoch                  uint32                              `json:"epoch"`
	NumTxs                 uint32                              `json:"numTxs"`
	AccumulatedFees        string                              `json:"accumulatedFees,omitempty"`
	DeveloperFees          string                              `json:"developerFees,omitempty"`
	AccumulatedFeesInEpoch string                              `json:"accumulatedFeesInEpoch,omitempty"`
	DeveloperFeesInEpoch   string                              `json:"developerFeesInEpoch,omitempty"`
	Timestamp              time.Duration                       `json:"timestamp,omitempty"`
	EpochStartInfo         *api.EpochStartInfo                 `json:"epochStartInfo,omitempty"`
	EpochStartShardsData   []*api.EpochStartShardData          `json:"epochStartShardsData,omitempty"`
	ShardBlocks            []*api.NotarizedBlock               `json:"shardBlocks"`
	Transactions           []*transaction.ApiTransactionResult `json:"transactions"`
	Status                 string                              `json:"status,omitempty"`
}
