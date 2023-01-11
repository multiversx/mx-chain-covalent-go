package apiMocks

import "github.com/ElrondNetwork/covalent-indexer-go/api"

// ElrondHyperBlockEndPointStub -
type ElrondHyperBlockEndPointStub struct {
	GetHyperBlockCalled func(path string) (*api.ElrondHyperBlockApiResponse, error)
}

// GetHyperBlock -
func (ehb *ElrondHyperBlockEndPointStub) GetHyperBlock(path string) (*api.ElrondHyperBlockApiResponse, error) {
	if ehb.GetHyperBlockCalled != nil {
		return ehb.GetHyperBlockCalled(path)
	}

	return &api.ElrondHyperBlockApiResponse{}, nil
}
