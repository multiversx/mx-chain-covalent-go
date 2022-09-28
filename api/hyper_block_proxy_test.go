package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const hyperBlockPath = "/hyperBlock"

func startProxyServer(proxy api.HyperBlockProxy, path string) *gin.Engine {
	ws := gin.New()
	routes := ws.Group(path)
	routes.GET("/by-nonce/:nonce", proxy.GetHyperBlockByNonce)
	routes.GET("/by-hash/:hash", proxy.GetHyperBlockByHash)
	return ws
}

func loadResponse(t *testing.T, rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	require.Nil(t, err)
}

func TestHyperBlockProxy_GetHyperBlockByNonce(t *testing.T) {
	t.Parallel()

	requestedNonce := uint64(4)
	blockResponse := &api.CovalentHyperBlockApiResponse{
		Data:  []byte("abc"),
		Error: "",
		Code:  "success",
	}
	facade := &mock.HyperBlockFacadeStub{
		GetHyperBlockByNonceCalled: func(nonce uint64, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
			require.Equal(t, requestedNonce, nonce)
			return blockResponse, nil
		},
	}
	proxy, _ := api.NewHyperBlockProxy(facade)

	ws := startProxyServer(proxy, hyperBlockPath)

	requestPath := fmt.Sprintf("%s/by-nonce/%d", hyperBlockPath, requestedNonce)
	req, _ := http.NewRequest("GET", requestPath, nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &api.CovalentHyperBlockApiResponse{}
	loadResponse(t, resp.Body, apiResp)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, apiResp, blockResponse)
}
