package apiMocks

import "github.com/ElrondNetwork/covalent-indexer-go/api"

// HyperBlockFacadeStub -
type HyperBlockFacadeStub struct {
	GetHyperBlockByNonceCalled func(nonce uint64, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error)
	GetHyperBlockByHashCalled  func(hash string, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error)
}

// GetHyperBlockByNonce -
func (hbf *HyperBlockFacadeStub) GetHyperBlockByNonce(nonce uint64, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
	if hbf.GetHyperBlockByNonceCalled != nil {
		return hbf.GetHyperBlockByNonceCalled(nonce, options)
	}

	return nil, nil
}

// GetHyperBlockByHash -
func (hbf *HyperBlockFacadeStub) GetHyperBlockByHash(hash string, options api.HyperBlockQueryOptions) (*api.CovalentHyperBlockApiResponse, error) {
	if hbf.GetHyperBlockByHashCalled != nil {
		return hbf.GetHyperBlockByHashCalled(hash, options)
	}

	return nil, nil
}
