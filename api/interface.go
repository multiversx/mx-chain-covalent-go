package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-covalent-go/cmd/proxy/config"
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

// MultiversxHyperBlockEndpointHandler should fetch hyper block api responses from Multiversx
type MultiversxHyperBlockEndpointHandler interface {
	GetHyperBlock(path string) (*MultiversxHyperBlockApiResponse, error)
}

// HyperBlockFacadeHandler defines the actions needed for fetching of hyperBlocks from Multiversx proxy in covalent format
type HyperBlockFacadeHandler interface {
	GetHyperBlockByNonce(nonce uint64, options config.HyperBlockQueryOptions) (*CovalentHyperBlockApiResponse, error)
	GetHyperBlockByHash(hash string, options config.HyperBlockQueryOptions) (*CovalentHyperBlockApiResponse, error)
	GetHyperBlocksByInterval(noncesInterval *Interval, options config.HyperBlocksQueryOptions) (*CovalentHyperBlocksApiResponse, error)
}

// HyperBlockProxy is the covalent proxy. It should be able to fetch hyper blocks from
// Multiversx proxy(json format), process them and provide avro schema defined hyper blocks(as byte array).
type HyperBlockProxy interface {
	GetHyperBlockByNonce(c *gin.Context)
	GetHyperBlockByHash(c *gin.Context)
	GetHyperBlocksByInterval(c *gin.Context)
}
