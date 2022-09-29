package apiMocks

import "github.com/ElrondNetwork/covalent-indexer-go/api"

type ElrondHyperBlockEndPointStub struct {
	GetHyperBlockCalled func(path string) (*api.ElrondHyperBlockApiResponse, error)
}

func (ehb *ElrondHyperBlockEndPointStub) GetHyperBlock(path string) (*api.ElrondHyperBlockApiResponse, error) {
	if ehb.GetHyperBlockCalled != nil {
		return ehb.GetHyperBlockCalled(path)
	}

	return nil, nil
}
