package facade

import (
	"fmt"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock/apiMocks"
	"github.com/ElrondNetwork/elrond-go/api/shared"
	"github.com/elodina/go-avro"
	"github.com/stretchr/testify/require"
)

func TestHyperBlockFacade_GetHyperBlockByNonce(t *testing.T) {
	elrondProxyUrl := "url"
	requestedNonce := uint64(4)

	elrondApiResponse := &api.ElrondHyperBlockApiResponse{
		Data: api.ElrondHyperBlockApiResponsePayload{HyperBlock: hyperBlock.HyperBlock{
			Hash: "hash",
		}},
		Error: "",
		Code:  shared.ReturnCodeSuccess,
	}
	elrondEndPoint := &apiMocks.ElrondHyperBlockEndPointStub{
		GetHyperBlockCalled: func(path string) (*api.ElrondHyperBlockApiResponse, error) {
			require.Equal(t, fmt.Sprintf("%s%s/%d", elrondProxyUrl, hyperBlockPathByNonce, requestedNonce), path)
			return elrondApiResponse, nil
		},
	}

	blockResult := &schema.BlockResult{
		Block: &schema.Block{
			Hash: []byte(elrondApiResponse.Data.HyperBlock.Hash),
		},
	}
	processor := &mock.HyperBlockProcessorStub{
		ProcessCalled: func(hyperBlock *hyperBlock.HyperBlock) (*schema.BlockResult, error) {
			require.Equal(t, &elrondApiResponse.Data.HyperBlock, hyperBlock)
			return blockResult, nil
		},
	}

	encodedBlock := []byte("encodedBlock")
	encoder := &mock.AvroEncoderStub{
		EncodeCalled: func(record avro.AvroRecord) ([]byte, error) {
			require.Equal(t, blockResult, record)
			return encodedBlock, nil
		},
	}

	facade, _ := NewHyperBlockFacade(elrondProxyUrl, encoder, elrondEndPoint, processor)

	block, err := facade.GetHyperBlockByNonce(4, api.HyperBlockQueryOptions{})
	require.Nil(t, err)
	require.Equal(t, &api.CovalentHyperBlockApiResponse{
		Data:  encodedBlock,
		Error: "",
		Code:  shared.ReturnCodeSuccess,
	}, block)
}
