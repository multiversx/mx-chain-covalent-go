package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("api")

type multiversxHyperBlockEndPoint struct {
	httpClient HTTPClient
}

// NewMultiversxHyperBlockEndPoint will create a handler which can fetch hyper blocks from Multiversx gateway
func NewMultiversxHyperBlockEndPoint(httpClient HTTPClient) (*multiversxHyperBlockEndPoint, error) {
	if httpClient == nil {
		return nil, errNilHttpServer
	}

	return &multiversxHyperBlockEndPoint{
		httpClient: httpClient,
	}, nil
}

// GetHyperBlock will fetch an MultiversxHyperBlockApiResponse from provided path
func (hpe *multiversxHyperBlockEndPoint) GetHyperBlock(path string) (*MultiversxHyperBlockApiResponse, error) {
	resp, err := hpe.httpClient.Get(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			errNotCritical := resp.Body.Close()
			if errNotCritical != nil {
				log.Warn("close body", "error", errNotCritical.Error())
			}
		}
	}()

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response MultiversxHyperBlockApiResponse
	err = json.Unmarshal(responseBodyBytes, &response)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, multiversx proxy response error: %s", resp.StatusCode, response.Error)
	}

	return &response, nil
}
