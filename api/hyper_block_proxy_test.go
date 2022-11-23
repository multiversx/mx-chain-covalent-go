package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/cmd/proxy/config"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock/apiMocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const (
	hyperBlockPath  = "/hyperblock"
	hyperBlocksPath = "/hyperblocks"
)

func startProxyServer(proxy api.HyperBlockProxy) *gin.Engine {
	ws := gin.New()

	routes := ws.Group(hyperBlockPath)
	routes.GET("/by-nonce/:nonce", proxy.GetHyperBlockByNonce)
	routes.GET("/by-hash/:hash", proxy.GetHyperBlockByHash)

	ws.Group(hyperBlocksPath).GET("", proxy.GetHyperBlocksByInterval)

	return ws
}

func loadResponse(t *testing.T, rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	require.Nil(t, err)
}

func serveHTTPRequest(t *testing.T, ws *gin.Engine, path string, expectedStatus int) *bytes.Buffer {
	req, err := http.NewRequest("GET", path, nil)
	require.Nil(t, err)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	require.Equal(t, expectedStatus, resp.Code)

	return resp.Body
}

func sendRequest(t *testing.T, ws *gin.Engine, path string, expectedStatus int) *api.CovalentHyperBlockApiResponse {
	body := serveHTTPRequest(t, ws, path, expectedStatus)

	apiResp := &api.CovalentHyperBlockApiResponse{}
	loadResponse(t, body, apiResp)

	return apiResp
}

func sendHyperBlocksRequest(t *testing.T, ws *gin.Engine, path string, expectedStatus int) *api.CovalentHyperBlocksApiResponse {
	body := serveHTTPRequest(t, ws, path, expectedStatus)

	apiResp := &api.CovalentHyperBlocksApiResponse{}
	loadResponse(t, body, apiResp)

	return apiResp
}

func getConfig() config.Config {
	return config.Config{
		HyperBlocksBatchSize: 10,
	}
}

func TestNewHyperBlockProxy(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		proxy, err := api.NewHyperBlockProxy(&apiMocks.HyperBlockFacadeStub{}, getConfig())
		require.Nil(t, err)
		require.NotNil(t, proxy)
	})

	t.Run("nil facade, should return error", func(t *testing.T) {
		t.Parallel()

		proxy, err := api.NewHyperBlockProxy(nil, getConfig())
		require.Nil(t, proxy)
		require.Equal(t, api.ErrNilHyperBlockFacade, err)
	})

	t.Run("invalid hyper blocks batch size, should return error", func(t *testing.T) {
		t.Parallel()

		proxy, err := api.NewHyperBlockProxy(&apiMocks.HyperBlockFacadeStub{}, config.Config{})
		require.Nil(t, proxy)
		require.ErrorIs(t, err, api.ErrInvalidHyperBlocksBatchSize)
	})
}

func TestGetNonceFromRequest_MissingNonce_ShouldReturnError(t *testing.T) {
	t.Parallel()

	context := &gin.Context{
		Params: []gin.Param{
			{
				Key:   "nonce",
				Value: "",
			},
		},
	}

	nonce, err := api.GetNonceFromRequest(context)
	require.Equal(t, uint64(0), nonce)
	require.Equal(t, api.ErrInvalidBlockNonce, err)
}

func TestHyperBlockProxy_GetHyperBlockByNonce(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		requestedNonce := uint64(4)
		blockResponse := &api.CovalentHyperBlockApiResponse{
			Data:  []byte("abc"),
			Error: "",
			Code:  "success",
		}
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlockByNonceCalled: func(nonce uint64, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				require.Equal(t, requestedNonce, nonce)
				return blockResponse, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())
		ws := startProxyServer(proxy)
		requestPath := fmt.Sprintf("%s/by-nonce/%d", hyperBlockPath, requestedNonce)
		apiResp := sendRequest(t, ws, requestPath, http.StatusOK)
		require.Equal(t, apiResp, blockResponse)
	})

	t.Run("invalid nonce, should error", func(t *testing.T) {
		t.Parallel()

		getHyperBlockFromFacadeCalled := false
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlockByNonceCalled: func(nonce uint64, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				getHyperBlockFromFacadeCalled = true
				return nil, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s/by-nonce/abc", hyperBlockPath)
		apiResp := sendRequest(t, ws, requestPath, http.StatusBadRequest)
		require.False(t, getHyperBlockFromFacadeCalled)
		require.Empty(t, apiResp.Data)
		require.Equal(t, apiResp.Code, api.ReturnCodeRequestError)
		require.True(t, strings.Contains(apiResp.Error, "abc"))
	})

	t.Run("could not get hyper block from facade, should error", func(t *testing.T) {
		t.Parallel()

		errFacade := errors.New("error getting hyper block from facade")
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlockByNonceCalled: func(nonce uint64, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				return nil, errFacade
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s/by-nonce/4", hyperBlockPath)
		apiResp := sendRequest(t, ws, requestPath, http.StatusInternalServerError)
		require.Equal(t, &api.CovalentHyperBlockApiResponse{
			Data:  nil,
			Error: errFacade.Error(),
			Code:  api.ReturnCodeInternalError,
		}, apiResp)
	})
}

func TestHyperBlockProxy_GetHyperBlocksByInterval(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		startNonce := uint64(4)
		endNonce := uint64(8)
		blockResponse := &api.CovalentHyperBlocksApiResponse{
			Data:  [][]byte{[]byte("first"), []byte("second")},
			Error: "",
			Code:  "success",
		}
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlocksByIntervalCalled: func(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) (*api.CovalentHyperBlocksApiResponse, error) {
				require.Equal(t, &api.Interval{
					Start: startNonce,
					End:   endNonce,
				}, noncesInterval)
				return blockResponse, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())

		ws := startProxyServer(proxy)
		requestPath := fmt.Sprintf("%s?startNonce=%d&endNonce=%d", hyperBlocksPath, startNonce, endNonce)

		apiResp := sendHyperBlocksRequest(t, ws, requestPath, http.StatusOK)
		require.Equal(t, apiResp, blockResponse)
	})

	t.Run("invalid start nonce, should error", func(t *testing.T) {
		t.Parallel()

		getHyperBlockFromFacadeCalled := false
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlocksByIntervalCalled: func(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) (*api.CovalentHyperBlocksApiResponse, error) {
				getHyperBlockFromFacadeCalled = true
				return nil, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s?startNonce=abc&endNonce=8", hyperBlocksPath)
		apiResp := sendHyperBlocksRequest(t, ws, requestPath, http.StatusBadRequest)
		require.False(t, getHyperBlockFromFacadeCalled)
		require.Empty(t, apiResp.Data)
		require.Equal(t, apiResp.Code, api.ReturnCodeRequestError)
		require.True(t, strings.Contains(apiResp.Error, "abc"))
	})

	t.Run("invalid end nonce, should error", func(t *testing.T) {
		t.Parallel()

		getHyperBlockFromFacadeCalled := false
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlocksByIntervalCalled: func(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) (*api.CovalentHyperBlocksApiResponse, error) {
				getHyperBlockFromFacadeCalled = true
				return nil, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s?startNonce=4&endNonce=cde", hyperBlocksPath)
		apiResp := sendHyperBlocksRequest(t, ws, requestPath, http.StatusBadRequest)
		require.False(t, getHyperBlockFromFacadeCalled)
		require.Empty(t, apiResp.Data)
		require.Equal(t, apiResp.Code, api.ReturnCodeRequestError)
		require.True(t, strings.Contains(apiResp.Error, "cde"))
	})

	t.Run("missing nonce, should error", func(t *testing.T) {
		t.Parallel()

		getHyperBlockFromFacadeCalled := false
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlocksByIntervalCalled: func(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) (*api.CovalentHyperBlocksApiResponse, error) {
				getHyperBlockFromFacadeCalled = true
				return nil, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s?startNonce=4", hyperBlocksPath)
		apiResp := sendHyperBlocksRequest(t, ws, requestPath, http.StatusBadRequest)
		require.False(t, getHyperBlockFromFacadeCalled)
		require.Empty(t, apiResp.Data)
		require.Equal(t, apiResp.Code, api.ReturnCodeRequestError)
		require.True(t, strings.Contains(apiResp.Error, api.ErrMissingQueryParameter.Error()))
		require.True(t, strings.Contains(apiResp.Error, "endNonce"))
	})

	t.Run("could not get hyper blocks from facade, should error", func(t *testing.T) {
		t.Parallel()

		errFacade := errors.New("error getting hyper block from facade")
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlocksByIntervalCalled: func(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) (*api.CovalentHyperBlocksApiResponse, error) {
				return nil, errFacade
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s?startNonce=4&endNonce=8", hyperBlocksPath)
		apiResp := sendHyperBlocksRequest(t, ws, requestPath, http.StatusInternalServerError)
		require.Equal(t, &api.CovalentHyperBlocksApiResponse{
			Data:  nil,
			Error: errFacade.Error(),
			Code:  api.ReturnCodeInternalError,
		}, apiResp)
	})
}

func TestHyperBlockProxy_GetHyperBlockByHash(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		requestedHash := "ff"
		blockResponse := &api.CovalentHyperBlockApiResponse{
			Data:  []byte("abc"),
			Error: "",
			Code:  "success",
		}
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlockByHashCalled: func(hash string, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				require.Equal(t, requestedHash, hash)
				return blockResponse, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())
		ws := startProxyServer(proxy)
		requestPath := fmt.Sprintf("%s/by-hash/%s", hyperBlockPath, requestedHash)
		apiResp := sendRequest(t, ws, requestPath, http.StatusOK)
		require.Equal(t, apiResp, blockResponse)
	})

	t.Run("invalid hash, should error", func(t *testing.T) {
		t.Parallel()

		getHyperBlockFromFacadeCalled := false
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlockByHashCalled: func(hash string, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				getHyperBlockFromFacadeCalled = true
				return nil, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s/by-hash/zx", hyperBlockPath)
		apiResp := sendRequest(t, ws, requestPath, http.StatusBadRequest)
		require.False(t, getHyperBlockFromFacadeCalled)
		require.Empty(t, apiResp.Data)
		require.Equal(t, apiResp.Code, api.ReturnCodeRequestError)
		require.Equal(t, apiResp.Error, api.ErrInvalidBlockHash.Error())
	})

	t.Run("could not get hyper block from facade, should error", func(t *testing.T) {
		t.Parallel()

		errFacade := errors.New("error getting hyper block from facade")
		facade := &apiMocks.HyperBlockFacadeStub{
			GetHyperBlockByHashCalled: func(hash string, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				return nil, errFacade
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade, getConfig())
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s/by-hash/ff", hyperBlockPath)
		apiResp := sendRequest(t, ws, requestPath, http.StatusInternalServerError)
		require.Equal(t, &api.CovalentHyperBlockApiResponse{
			Data:  nil,
			Error: errFacade.Error(),
			Code:  api.ReturnCodeInternalError,
		}, apiResp)
	})
}
