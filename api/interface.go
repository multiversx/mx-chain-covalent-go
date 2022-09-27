package api

import (
	"context"
	"net/http"

	"github.com/elodina/go-avro"
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

// HyperBlockFacadeHandler defines the actions needed for fetching of hyperBlocks from Elrond proxy
type HyperBlockFacadeHandler interface {
	GetHyperBlockByNonce(nonce uint64, options HyperBlockQueryOptions) (*HyperBlockApiResponse, error)
	GetHyperBlockByHash(hash string, options HyperBlockQueryOptions) (*HyperBlockApiResponse, error)
}

// HyperBlockProxy is the covalent proxy. It should be able to fetch hyper blocks from
// Elrond proxy(json format), process them and provide avro schema defined hyper blocks(as byte array).
type HyperBlockProxy interface {
	GetHyperBlockByNonce(c *gin.Context)
	GetHyperBlockByHash(c *gin.Context)
}

// AvroEncoder should be able to encode any avro schema in a byte array
type AvroEncoder interface {
	Encode(record avro.AvroRecord) ([]byte, error)
}
