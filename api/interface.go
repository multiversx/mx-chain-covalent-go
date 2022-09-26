package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Close() error
}

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// HyperBlockFacadeHandler defines the actions needed for fetching of hyperBlocks from Elrond proxy
type HyperBlockFacadeHandler interface {
	GetHyperBlockByNonce(nonce uint64, options HyperBlockQueryOptions) (*HyperBlockApiResponse, error)
	GetHyperBlockByHash(hash string, options HyperBlockQueryOptions) (*HyperBlockApiResponse, error)
}

type HyperBlockProxyFetcher interface {
	GetHyperBlockByNonce(c *gin.Context)
	GetHyperBlockByHash(c *gin.Context)
}
