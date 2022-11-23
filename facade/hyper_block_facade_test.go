package facade

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/cmd/proxy/config"
	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock/apiMocks"
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
		Code:  api.ReturnCodeSuccess,
	}
	elrondEndPoint := &apiMocks.ElrondHyperBlockEndPointStub{
		GetHyperBlockCalled: func(path string) (*api.ElrondHyperBlockApiResponse, error) {
			require.Equal(t, fmt.Sprintf("%s%s/%d", elrondProxyUrl, hyperBlockPathByNonce, requestedNonce), path)
			return elrondApiResponse, nil
		},
	}

	blockResult := &schema.HyperBlock{
		Hash: []byte(elrondApiResponse.Data.HyperBlock.Hash),
	}
	processor := &mock.HyperBlockProcessorStub{
		ProcessCalled: func(hyperBlock *hyperBlock.HyperBlock) (*schema.HyperBlock, error) {
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

	block, err := facade.GetHyperBlockByNonce(4, config.HyperBlockQueryOptions{})
	require.Nil(t, err)
	require.Equal(t, &api.CovalentHyperBlockApiResponse{
		Data:  encodedBlock,
		Error: "",
		Code:  api.ReturnCodeSuccess,
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
		Code:  api.ReturnCodeSuccess,
	}
	elrondEndPoint := &apiMocks.ElrondHyperBlockEndPointStub{
		GetHyperBlockCalled: func(path string) (*api.ElrondHyperBlockApiResponse, error) {
			require.Equal(t, fmt.Sprintf("%s%s/%s", elrondProxyUrl, hyperBlockPathByHash, requestedHash), path)
			return elrondApiResponse, nil
		},
	}

	blockResult := &schema.HyperBlock{
		Hash: []byte(elrondApiResponse.Data.HyperBlock.Hash),
	}
	processor := &mock.HyperBlockProcessorStub{
		ProcessCalled: func(hyperBlock *hyperBlock.HyperBlock) (*schema.HyperBlock, error) {
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

	block, err := facade.GetHyperBlockByHash(requestedHash, config.HyperBlockQueryOptions{})
	require.Nil(t, err)
	require.Equal(t, &api.CovalentHyperBlockApiResponse{
		Data:  encodedBlock,
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}, block)
}

func TestHyperBlockFacade_buildUrlWithBlockQueryOptions(t *testing.T) {
	t.Parallel()

	path := "path"

	fullPath := buildUrlWithBlockQueryOptions(path, config.HyperBlockQueryOptions{
		WithLogs:            false,
		WithAlteredAccounts: false,
	})
	require.Equal(t, path, fullPath)

	fullPath = buildUrlWithBlockQueryOptions(path, config.HyperBlockQueryOptions{
		WithLogs:            true,
		WithAlteredAccounts: false,
	})
	require.Equal(t, fmt.Sprintf("%s?%s=true", path, api.UrlParameterWithLogs), fullPath)

	fullPath = buildUrlWithBlockQueryOptions(path, config.HyperBlockQueryOptions{
		WithLogs:            false,
		WithAlteredAccounts: true,
	})
	require.Equal(t, fmt.Sprintf("%s?%s=true", path, api.UrlParameterWithAlteredAccounts), fullPath)

	fullPath = buildUrlWithBlockQueryOptions(path, config.HyperBlockQueryOptions{
		WithLogs:            true,
		WithAlteredAccounts: true,
	})
	require.Equal(t, fmt.Sprintf("%s?%s=true&%s=true", path, api.UrlParameterWithAlteredAccounts, api.UrlParameterWithLogs), fullPath)

	fullPath = buildUrlWithBlockQueryOptions(path, config.HyperBlockQueryOptions{
		WithLogs:            true,
		WithAlteredAccounts: true,
		Tokens:              "all",
		NotarizedAtSource:   true,
	})
	require.Equal(t,
		fmt.Sprintf("%s?%s=true&%s=all&%s=true&%s=true",
			path,
			api.UrlParameterNotarizedAtSource,
			api.UrlParameterTokens,
			api.UrlParameterWithAlteredAccounts,
			api.UrlParameterWithLogs,
		),
		fullPath)
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
			ProcessCalled: func(hyperBlock *hyperBlock.HyperBlock) (*schema.HyperBlock, error) {
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

func getNonceFromRequest(t *testing.T, request string) uint64 {
	splits := strings.Split(request, "/")
	numSplits := len(splits)
	require.True(t, numSplits >= 1)

	nonceFromRequestStr := splits[numSplits-1]
	nonceFromRequestInt, err := strconv.Atoi(nonceFromRequestStr)
	require.Nil(t, err)

	return uint64(nonceFromRequestInt)
}

func requireNonceInInterval(t *testing.T, nonce uint64, interval *api.Interval) {
	require.True(t, nonce >= interval.Start && nonce <= interval.End)
}

func TestHyperBlockFacade_GetHyperBlocksByInterval(t *testing.T) {
	t.Parallel()

	elrondProxyUrl := "url"

	interval := &api.Interval{
		Start: 4,
		End:   45,
	}
	numNonces := interval.End - interval.Start + 1

	elrondEndPointCallsCt := uint64(0)
	processHyperBlocksCt := uint64(0)
	encodeCt := uint64(0)

	elrondEndPoint := &apiMocks.ElrondHyperBlockEndPointStub{
		GetHyperBlockCalled: func(path string) (*api.ElrondHyperBlockApiResponse, error) {
			atomic.AddUint64(&elrondEndPointCallsCt, 1)

			nonceFromRequest := getNonceFromRequest(t, path)
			requireNonceInInterval(t, nonceFromRequest, interval)
			require.Equal(t, fmt.Sprintf("%s%s/%d", elrondProxyUrl, hyperBlockPathByNonce, nonceFromRequest), path)

			return &api.ElrondHyperBlockApiResponse{
				Data: api.ElrondHyperBlockApiResponsePayload{
					HyperBlock: hyperBlock.HyperBlock{
						Nonce: nonceFromRequest,
					}},
				Error: "",
				Code:  api.ReturnCodeSuccess,
			}, nil
		},
	}

	processor := &mock.HyperBlockProcessorStub{
		ProcessCalled: func(hyperBlock *hyperBlock.HyperBlock) (*schema.HyperBlock, error) {
			atomic.AddUint64(&processHyperBlocksCt, 1)

			requireNonceInInterval(t, hyperBlock.Nonce, interval)
			return &schema.HyperBlock{
				Nonce: int64(hyperBlock.Nonce),
			}, nil
		},
	}

	encoder := &mock.AvroEncoderStub{
		EncodeCalled: func(record avro.AvroRecord) ([]byte, error) {
			atomic.AddUint64(&encodeCt, 1)

			hyperBlockRecord, castOk := record.(*schema.HyperBlock)
			require.True(t, castOk)

			requireNonceInInterval(t, uint64(hyperBlockRecord.Nonce), interval)
			encodedBlock := []byte(fmt.Sprintf("encodedBlock%d", hyperBlockRecord.Nonce))
			return encodedBlock, nil
		},
	}

	facade, _ := NewHyperBlockFacade(elrondProxyUrl, encoder, elrondEndPoint, processor)

	expectedEncodedHyperBlocks := make([][]byte, 0)
	for nonce := interval.Start; nonce <= interval.End; nonce++ {
		encodedBlock := []byte(fmt.Sprintf("encodedBlock%d", nonce))
		expectedEncodedHyperBlocks = append(expectedEncodedHyperBlocks, encodedBlock)
	}
	blocks, err := facade.GetHyperBlocksByInterval(&api.Interval{
		Start: 4,
		End:   45,
	}, config.HyperBlocksQueryOptions{
		BatchSize: 10,
	})
	require.Nil(t, err)
	require.Equal(t, &api.CovalentHyperBlocksApiResponse{
		Data:  expectedEncodedHyperBlocks,
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}, blocks)

	require.Equal(t, numNonces, elrondEndPointCallsCt)
	require.Equal(t, numNonces, processHyperBlocksCt)
	require.Equal(t, numNonces, encodeCt)
}
