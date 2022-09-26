package api

import (
	"context"

	"github.com/gin-gonic/gin"
)

type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Close() error
}

// HyperBlockFacadeHandler defines the actions needed for fetching of hyperBlocks from Elrond proxy
type HyperBlockFacadeHandler interface {
	GetHyperBlockByNonce(nonce uint64, options HyperBlockQueryOptions) (*HyperblockApiResponse, error)
	GetHyperBlockByHash(hash string, options HyperBlockQueryOptions) (*HyperblockApiResponse, error)
}

type HyperBlockProxyFetcher interface {
	GetHyperBlockByNonce(c *gin.Context)
	GetHyperBlockByHash(c *gin.Context)
}
