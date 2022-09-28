package api_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go/api/shared"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const hyperBlockPath = "/hyperBlock"

func startProxyServer(proxy api.HyperBlockProxy) *gin.Engine {
	ws := gin.New()

	routes := ws.Group(hyperBlockPath)
	routes.GET("/by-nonce/:nonce", proxy.GetHyperBlockByNonce)
	routes.GET("/by-hash/:hash", proxy.GetHyperBlockByHash)

	return ws
}

func loadResponse(t *testing.T, rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	require.Nil(t, err)
}

func sendRequest(t *testing.T, ws *gin.Engine, path string, expectedStatus int) *api.CovalentHyperBlockApiResponse {
	req, err := http.NewRequest("GET", path, nil)
	require.Nil(t, err)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	require.Equal(t, expectedStatus, resp.Code)

	apiResp := &api.CovalentHyperBlockApiResponse{}
	loadResponse(t, resp.Body, apiResp)

	return apiResp
}

func TestNewHyperBlockProxy(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		proxy, err := api.NewHyperBlockProxy(&mock.HyperBlockFacadeStub{})
		require.Nil(t, err)
		require.NotNil(t, proxy)
	})

	t.Run("nil facade, should return error", func(t *testing.T) {
		t.Parallel()

		proxy, err := api.NewHyperBlockProxy(nil)
		require.Nil(t, proxy)
		require.Equal(t, api.ErrNilHyperBlockFacade, err)
	})
}

func TestGetNonceFromRequest_MissingNonce_ShouldReturnError(t *testing.T) {
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
		facade := &mock.HyperBlockFacadeStub{
			GetHyperBlockByNonceCalled: func(nonce uint64, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				require.Equal(t, requestedNonce, nonce)
				return blockResponse, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade)
		ws := startProxyServer(proxy)
		requestPath := fmt.Sprintf("%s/by-nonce/%d", hyperBlockPath, requestedNonce)
		apiResp := sendRequest(t, ws, requestPath, http.StatusOK)
		require.Equal(t, apiResp, blockResponse)
	})

	t.Run("invalid nonce, should error", func(t *testing.T) {
		t.Parallel()

		getHyperBlockFromFacadeCalled := false
		facade := &mock.HyperBlockFacadeStub{
			GetHyperBlockByNonceCalled: func(nonce uint64, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				getHyperBlockFromFacadeCalled = true
				return nil, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade)
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s/by-nonce/abc", hyperBlockPath)
		apiResp := sendRequest(t, ws, requestPath, http.StatusBadRequest)
		require.False(t, getHyperBlockFromFacadeCalled)
		require.Empty(t, apiResp.Data)
		require.Equal(t, apiResp.Code, shared.ReturnCodeRequestError)
		require.True(t, strings.Contains(apiResp.Error, "abc"))
	})

	t.Run("could not get hyper block from facade, should error", func(t *testing.T) {
		t.Parallel()

		errFacade := errors.New("error getting hyper block from facade")
		facade := &mock.HyperBlockFacadeStub{
			GetHyperBlockByNonceCalled: func(nonce uint64, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				return nil, errFacade
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade)
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s/by-nonce/4", hyperBlockPath)
		apiResp := sendRequest(t, ws, requestPath, http.StatusInternalServerError)
		require.Equal(t, &api.CovalentHyperBlockApiResponse{
			Data:  nil,
			Error: errFacade.Error(),
			Code:  shared.ReturnCodeInternalError,
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
		facade := &mock.HyperBlockFacadeStub{
			GetHyperBlockByHashCalled: func(hash string, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				require.Equal(t, requestedHash, hash)
				return blockResponse, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade)
		ws := startProxyServer(proxy)
		requestPath := fmt.Sprintf("%s/by-hash/%s", hyperBlockPath, requestedHash)
		apiResp := sendRequest(t, ws, requestPath, http.StatusOK)
		require.Equal(t, apiResp, blockResponse)
	})

	t.Run("invalid hash, should error", func(t *testing.T) {
		t.Parallel()

		getHyperBlockFromFacadeCalled := false
		facade := &mock.HyperBlockFacadeStub{
			GetHyperBlockByHashCalled: func(hash string, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				getHyperBlockFromFacadeCalled = true
				return nil, nil
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade)
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s/by-hash/zx", hyperBlockPath)
		apiResp := sendRequest(t, ws, requestPath, http.StatusBadRequest)
		require.False(t, getHyperBlockFromFacadeCalled)
		require.Empty(t, apiResp.Data)
		require.Equal(t, apiResp.Code, shared.ReturnCodeRequestError)
		require.Equal(t, apiResp.Error, api.ErrInvalidBlockHash.Error())
	})

	t.Run("could not get hyper block from facade, should error", func(t *testing.T) {
		t.Parallel()

		errFacade := errors.New("error getting hyper block from facade")
		facade := &mock.HyperBlockFacadeStub{
			GetHyperBlockByHashCalled: func(hash string, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
				return nil, errFacade
			},
		}
		proxy, _ := api.NewHyperBlockProxy(facade)
		ws := startProxyServer(proxy)

		requestPath := fmt.Sprintf("%s/by-hash/ff", hyperBlockPath)
		apiResp := sendRequest(t, ws, requestPath, http.StatusInternalServerError)
		require.Equal(t, &api.CovalentHyperBlockApiResponse{
			Data:  nil,
			Error: errFacade.Error(),
			Code:  shared.ReturnCodeInternalError,
		}, apiResp)
	})
}
