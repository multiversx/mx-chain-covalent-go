package api

import "github.com/ElrondNetwork/covalent-indexer-go/schema"

type HyperBlockFacade struct {
}

func (hpf *HyperBlockFacade) GetHyperBlockByNonce(nonce uint64, options HyperBlockQueryOptions) (*schema.BlockResult, error) {
	return &schema.BlockResult{
		Block: &schema.Block{
			Nonce: int64(nonce),
		},
	}, nil
}
func (hpf *HyperBlockFacade) GetHyperBlockByHash(hash string, options HyperBlockQueryOptions) (*schema.BlockResult, error) {
	return nil, nil
}
