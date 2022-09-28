package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPServer defines what an http server should do
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Close() error
}

// HTTPClient defines what a client which should be able to GET requests should do
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// ElrondHyperBlockEndpointHandler should fetch hyper block api responses from elrond
type ElrondHyperBlockEndpointHandler interface {
	GetHyperBlock(path string) (*ElrondHyperBlockApiResponse, error)
}

// HyperBlockFacadeHandler defines the actions needed for fetching of hyperBlocks from Elrond proxy in covalent format
type HyperBlockFacadeHandler interface {
	GetHyperBlockByNonce(nonce uint64, options HyperBlockQueryOptions) (*CovalentHyperBlockApiResponse, error)
	GetHyperBlockByHash(hash string, options HyperBlockQueryOptions) (*CovalentHyperBlockApiResponse, error)
}

// HyperBlockProxy is the covalent proxy. It should be able to fetch hyper blocks from
// Elrond proxy(json format), process them and provide avro schema defined hyper blocks(as byte array).
type HyperBlockProxy interface {
	GetHyperBlockByNonce(c *gin.Context)
	GetHyperBlockByHash(c *gin.Context)
}
