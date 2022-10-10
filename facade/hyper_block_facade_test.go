package facade

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock/apiMocks"
	"github.com/ElrondNetwork/elrond-go/api/shared"
	"github.com/elodina/go-avro"
	"github.com/stretchr/testify/require"
)

func TestNewHyperBlockFacade(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		facade, err := NewHyperBlockFacade("url", &mock.AvroEncoderStub{}, &apiMocks.ElrondHyperBlockEndPointStub{}, &mock.HyperBlockProcessorStub{})
		require.NotNil(t, facade)
		require.Nil(t, err)
	})

	t.Run("empty url, should return error", func(t *testing.T) {
		t.Parallel()

		facade, err := NewHyperBlockFacade("", &mock.AvroEncoderStub{}, &apiMocks.ElrondHyperBlockEndPointStub{}, &mock.HyperBlockProcessorStub{})
		require.Nil(t, facade)
		require.Equal(t, errEmptyElrondProxyUrl, err)
	})

	t.Run("nil encoder, should return error", func(t *testing.T) {
		t.Parallel()

		facade, err := NewHyperBlockFacade("url", nil, &apiMocks.ElrondHyperBlockEndPointStub{}, &mock.HyperBlockProcessorStub{})
		require.Nil(t, facade)
		require.Equal(t, errNilAvroEncoder, err)
	})

	t.Run("nil elrond endpoint, should return error", func(t *testing.T) {
		t.Parallel()

		facade, err := NewHyperBlockFacade("url", &mock.AvroEncoderStub{}, nil, &mock.HyperBlockProcessorStub{})
		require.Nil(t, facade)
		require.Equal(t, errNilHyperBlockEndpointHandler, err)
	})

	t.Run("nil processor, should return error", func(t *testing.T) {
		t.Parallel()

		facade, err := NewHyperBlockFacade("url", &mock.AvroEncoderStub{}, &apiMocks.ElrondHyperBlockEndPointStub{}, nil)
		require.Nil(t, facade)
		require.Equal(t, errNilHyperBlockProcessor, err)
	})
}

func TestHyperBlockFacade_GetHyperBlockByNonce(t *testing.T) {
	t.Parallel()

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

	blockResult := &schemaV2.HyperBlock{
		Hash: []byte(elrondApiResponse.Data.HyperBlock.Hash),
	}
	processor := &mock.HyperBlockProcessorStub{
		ProcessCalled: func(hyperBlock *hyperBlock.HyperBlock) (*schemaV2.HyperBlock, error) {
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

func TestHyperBlockFacade_GetHyperBlockByHash(t *testing.T) {
	t.Parallel()

	elrondProxyUrl := "url"
	requestedHash := "hash"

	elrondApiResponse := &api.ElrondHyperBlockApiResponse{
		Data: api.ElrondHyperBlockApiResponsePayload{HyperBlock: hyperBlock.HyperBlock{
			Hash: "hash",
		}},
		Error: "",
		Code:  shared.ReturnCodeSuccess,
	}
	elrondEndPoint := &apiMocks.ElrondHyperBlockEndPointStub{
		GetHyperBlockCalled: func(path string) (*api.ElrondHyperBlockApiResponse, error) {
			require.Equal(t, fmt.Sprintf("%s%s/%s", elrondProxyUrl, hyperBlockPathByHash, requestedHash), path)
			return elrondApiResponse, nil
		},
	}

	blockResult := &schemaV2.HyperBlock{
		Hash: []byte(elrondApiResponse.Data.HyperBlock.Hash),
	}
	processor := &mock.HyperBlockProcessorStub{
		ProcessCalled: func(hyperBlock *hyperBlock.HyperBlock) (*schemaV2.HyperBlock, error) {
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

	block, err := facade.GetHyperBlockByHash(requestedHash, api.HyperBlockQueryOptions{})
	require.Nil(t, err)
	require.Equal(t, &api.CovalentHyperBlockApiResponse{
		Data:  encodedBlock,
		Error: "",
		Code:  shared.ReturnCodeSuccess,
	}, block)
}

func TestHyperBlockFacade_buildUrlWithBlockQueryOptions(t *testing.T) {
	t.Parallel()

	path := "path"

	fullPath := buildUrlWithBlockQueryOptions(path, api.HyperBlockQueryOptions{
		WithLogs:     false,
		WithBalances: false,
	})
	require.Equal(t, path, fullPath)

	fullPath = buildUrlWithBlockQueryOptions(path, api.HyperBlockQueryOptions{
		WithLogs:     true,
		WithBalances: false,
	})
	require.Equal(t, fmt.Sprintf("%s?%s=true", path, api.UrlParameterWithLogs), fullPath)

	fullPath = buildUrlWithBlockQueryOptions(path, api.HyperBlockQueryOptions{
		WithLogs:     false,
		WithBalances: true,
	})
	require.Equal(t, fmt.Sprintf("%s?%s=true", path, api.UrlParameterWithBalances), fullPath)

	fullPath = buildUrlWithBlockQueryOptions(path, api.HyperBlockQueryOptions{
		WithLogs:     true,
		WithBalances: true,
	})
	require.Equal(t, fmt.Sprintf("%s?%s=true&%s=true", path, api.UrlParameterWithBalances, api.UrlParameterWithLogs), fullPath)
}

func TestHyperBlockFacade_GetHyperBlock_ErrorCases(t *testing.T) {
	t.Parallel()

	elrondProxyUrl := "url"

	t.Run("cannot get hyper block from endpoint, expect error", func(t *testing.T) {
		t.Parallel()

		errGetHyperBlock := errors.New("error getting hyper block")
		elrondEndPoint := &apiMocks.ElrondHyperBlockEndPointStub{
			GetHyperBlockCalled: func(path string) (*api.ElrondHyperBlockApiResponse, error) {
				return nil, errGetHyperBlock
			},
		}

		facade, _ := NewHyperBlockFacade(
			elrondProxyUrl,
			&mock.AvroEncoderStub{},
			elrondEndPoint,
			&mock.HyperBlockProcessorStub{},
		)

		block, err := facade.getHyperBlock(elrondProxyUrl)
		require.Nil(t, block)
		require.Equal(t, errGetHyperBlock, err)
	})

	t.Run("cannot process hyper block, expect error", func(t *testing.T) {
		t.Parallel()

		errProcessor := errors.New("error processing hyper block")
		processor := &mock.HyperBlockProcessorStub{
			ProcessCalled: func(hyperBlock *hyperBlock.HyperBlock) (*schemaV2.HyperBlock, error) {
				return nil, errProcessor
			},
		}

		facade, _ := NewHyperBlockFacade(
			elrondProxyUrl,
			&mock.AvroEncoderStub{},
			&apiMocks.ElrondHyperBlockEndPointStub{},
			processor,
		)

		block, err := facade.getHyperBlock(elrondProxyUrl)
		require.Nil(t, block)
		require.Equal(t, errProcessor, err)
	})

	t.Run("cannot encode hyper block, expect error", func(t *testing.T) {
		t.Parallel()

		errEncoder := errors.New("error encoding hyper block")
		encoder := &mock.AvroEncoderStub{
			EncodeCalled: func(record avro.AvroRecord) ([]byte, error) {
				return nil, errEncoder
			},
		}

		facade, _ := NewHyperBlockFacade(
			elrondProxyUrl,
			encoder,
			&apiMocks.ElrondHyperBlockEndPointStub{},
			&mock.HyperBlockProcessorStub{},
		)

		block, err := facade.getHyperBlock(elrondProxyUrl)
		require.Nil(t, block)
		require.Equal(t, errEncoder, err)
	})
}
