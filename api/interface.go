package api

import (
	"context"

	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/gin-gonic/gin"
)

type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Close() error
}

// TODO: Replace return *schema.BlockResult with a new avro schema for hyperblocks

// HyperBlockFacadeHandler defines the actions needed for fetching of hyperBlocks from Elrond proxy
type HyperBlockFacadeHandler interface {
	GetHyperBlockByNonce(nonce uint64, options HyperBlockQueryOptions) (*schema.BlockResult, error)
	GetHyperBlockByHash(hash string, options HyperBlockQueryOptions) (*schema.BlockResult, error)
}

type HyperBlockProxyFetcher interface {
	GetHyperBlockByNonce(c *gin.Context)
	GetHyperBlockByHash(c *gin.Context)
}
