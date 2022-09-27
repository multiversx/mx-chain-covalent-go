package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type elrondHyperBlockEndPoint struct {
	httpClient HTTPClient
}

// NewElrondHyperBlockEndPoint will create a handler which can fetch hyper blocks from Elrond
func NewElrondHyperBlockEndPoint(httpClient HTTPClient) (*elrondHyperBlockEndPoint, error) {
	if httpClient == nil {
		return nil, errNilHttpServer
	}

	return &elrondHyperBlockEndPoint{
		httpClient: httpClient,
	}, nil
}

// GetHyperBlock will fetch an ElrondHyperBlockApiResponse from provided path
func (hpe *elrondHyperBlockEndPoint) GetHyperBlock(path string) (*ElrondHyperBlockApiResponse, error) {
	resp, err := hpe.httpClient.Get(path)
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

	var response ElrondHyperBlockApiResponse
	err = json.Unmarshal(responseBodyBytes, &response)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, elrond proxy response error: %s", resp.StatusCode, response.Error)
	}

	return &response, nil
}
