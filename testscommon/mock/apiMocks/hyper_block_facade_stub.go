package apiMocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/cmd/proxy/config"
)

// HyperBlockFacadeStub -
type HyperBlockFacadeStub struct {
	GetHyperBlockByNonceCalled     func(nonce uint64, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error)
	GetHyperBlockByHashCalled      func(hash string, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error)
	GetHyperBlocksByIntervalCalled func(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) (*api.CovalentHyperBlocksApiResponse, error)
}

// GetHyperBlockByNonce -
func (hbf *HyperBlockFacadeStub) GetHyperBlockByNonce(nonce uint64, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
	if hbf.GetHyperBlockByNonceCalled != nil {
		return hbf.GetHyperBlockByNonceCalled(nonce, options)
	}

	return nil, nil
}

// GetHyperBlockByHash -
func (hbf *HyperBlockFacadeStub) GetHyperBlockByHash(hash string, options config.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
	if hbf.GetHyperBlockByHashCalled != nil {
		return hbf.GetHyperBlockByHashCalled(hash, options)
	}

	return nil, nil
}

func (hbf *HyperBlockFacadeStub) GetHyperBlocksByInterval(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) (*api.CovalentHyperBlocksApiResponse, error) {
	if hbf.GetHyperBlocksByIntervalCalled != nil {
		return hbf.GetHyperBlocksByIntervalCalled(noncesInterval, options)
	}

	return nil, nil
}
