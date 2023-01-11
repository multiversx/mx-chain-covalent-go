package apiMocks

import "github.com/multiversx/mx-chain-covalent-go/api"

// MultiversxHyperBlockEndPointStub -
type MultiversxHyperBlockEndPointStub struct {
	GetHyperBlockCalled func(path string) (*api.MultiversxHyperBlockApiResponse, error)
}

// GetHyperBlock -
func (ehb *MultiversxHyperBlockEndPointStub) GetHyperBlock(path string) (*api.MultiversxHyperBlockApiResponse, error) {
	if ehb.GetHyperBlockCalled != nil {
		return ehb.GetHyperBlockCalled(path)
	}

	return &api.MultiversxHyperBlockApiResponse{}, nil
}
