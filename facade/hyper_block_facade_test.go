package facade

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/elodina/go-avro"
	"github.com/multiversx/mx-chain-covalent-go/api"
	"github.com/multiversx/mx-chain-covalent-go/cmd/proxy/config"
	"github.com/multiversx/mx-chain-covalent-go/hyperBlock"
	"github.com/multiversx/mx-chain-covalent-go/schema"
	"github.com/multiversx/mx-chain-covalent-go/testscommon/mock"
	"github.com/multiversx/mx-chain-covalent-go/testscommon/mock/apiMocks"
	"github.com/stretchr/testify/require"
)

func TestNewHyperBlockFacade(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		facade, err := NewHyperBlockFacade("url", &mock.AvroEncoderStub{}, &apiMocks.MultiversxHyperBlockEndPointStub{}, &mock.HyperBlockProcessorStub{})
		require.NotNil(t, facade)
		require.Nil(t, err)
	})

	t.Run("empty url, should return error", func(t *testing.T) {
		t.Parallel()

		facade, err := NewHyperBlockFacade("", &mock.AvroEncoderStub{}, &apiMocks.MultiversxHyperBlockEndPointStub{}, &mock.HyperBlockProcessorStub{})
		require.Nil(t, facade)
		require.Equal(t, errEmptyMultiversxProxyUrl, err)
	})

	t.Run("nil encoder, should return error", func(t *testing.T) {
		t.Parallel()

		facade, err := NewHyperBlockFacade("url", nil, &apiMocks.MultiversxHyperBlockEndPointStub{}, &mock.HyperBlockProcessorStub{})
		require.Nil(t, facade)
		require.Equal(t, errNilAvroEncoder, err)
	})

	t.Run("nil multiversx endpoint, should return error", func(t *testing.T) {
		t.Parallel()

		facade, err := NewHyperBlockFacade("url", &mock.AvroEncoderStub{}, nil, &mock.HyperBlockProcessorStub{})
		require.Nil(t, facade)
		require.Equal(t, errNilHyperBlockEndpointHandler, err)
	})

	t.Run("nil processor, should return error", func(t *testing.T) {
		t.Parallel()

		facade, err := NewHyperBlockFacade("url", &mock.AvroEncoderStub{}, &apiMocks.MultiversxHyperBlockEndPointStub{}, nil)
		require.Nil(t, facade)
		require.Equal(t, errNilHyperBlockProcessor, err)
	})
}

func TestHyperBlockFacade_GetHyperBlockByNonce(t *testing.T) {
	t.Parallel()

	multiversxProxyUrl := "url"
	requestedNonce := uint64(4)

	multiversxApiResponse := &api.MultiversxHyperBlockApiResponse{
		Data: api.MultiversxHyperBlockApiResponsePayload{HyperBlock: hyperBlock.HyperBlock{
			Hash: "hash",
		}},
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}
	multiversxEndPoint := &apiMocks.MultiversxHyperBlockEndPointStub{
		GetHyperBlockCalled: func(path string) (*api.MultiversxHyperBlockApiResponse, error) {
			require.Equal(t, fmt.Sprintf("%s%s/%d", multiversxProxyUrl, hyperBlockPathByNonce, requestedNonce), path)
			return multiversxApiResponse, nil
		},
	}

	blockResult := &schema.HyperBlock{
		Hash: []byte(multiversxApiResponse.Data.HyperBlock.Hash),
	}
	processor := &mock.HyperBlockProcessorStub{
		ProcessCalled: func(hyperBlock *hyperBlock.HyperBlock) (*schema.HyperBlock, error) {
			require.Equal(t, &multiversxApiResponse.Data.HyperBlock, hyperBlock)
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

	facade, _ := NewHyperBlockFacade(multiversxProxyUrl, encoder, multiversxEndPoint, processor)

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

	multiversxProxyUrl := "url"
	requestedHash := "hash"

	multiversxApiResponse := &api.MultiversxHyperBlockApiResponse{
		Data: api.MultiversxHyperBlockApiResponsePayload{HyperBlock: hyperBlock.HyperBlock{
			Hash: "hash",
		}},
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}
	multiversxEndPoint := &apiMocks.MultiversxHyperBlockEndPointStub{
		GetHyperBlockCalled: func(path string) (*api.MultiversxHyperBlockApiResponse, error) {
			require.Equal(t, fmt.Sprintf("%s%s/%s", multiversxProxyUrl, hyperBlockPathByHash, requestedHash), path)
			return multiversxApiResponse, nil
		},
	}

	blockResult := &schema.HyperBlock{
		Hash: []byte(multiversxApiResponse.Data.HyperBlock.Hash),
	}
	processor := &mock.HyperBlockProcessorStub{
		ProcessCalled: func(hyperBlock *hyperBlock.HyperBlock) (*schema.HyperBlock, error) {
			require.Equal(t, &multiversxApiResponse.Data.HyperBlock, hyperBlock)
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

	facade, _ := NewHyperBlockFacade(multiversxProxyUrl, encoder, multiversxEndPoint, processor)

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

	multiversxProxyUrl := "url"

	t.Run("cannot get hyper block from endpoint, expect error", func(t *testing.T) {
		t.Parallel()

		errGetHyperBlock := errors.New("error getting hyper block")
		multiversxEndPoint := &apiMocks.MultiversxHyperBlockEndPointStub{
			GetHyperBlockCalled: func(path string) (*api.MultiversxHyperBlockApiResponse, error) {
				return nil, errGetHyperBlock
			},
		}

		facade, _ := NewHyperBlockFacade(
			multiversxProxyUrl,
			&mock.AvroEncoderStub{},
			multiversxEndPoint,
			&mock.HyperBlockProcessorStub{},
		)

		block, err := facade.getHyperBlock(multiversxProxyUrl)
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
			multiversxProxyUrl,
			&mock.AvroEncoderStub{},
			&apiMocks.MultiversxHyperBlockEndPointStub{},
			processor,
		)

		block, err := facade.getHyperBlock(multiversxProxyUrl)
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
			multiversxProxyUrl,
			encoder,
			&apiMocks.MultiversxHyperBlockEndPointStub{},
			&mock.HyperBlockProcessorStub{},
		)

		block, err := facade.getHyperBlock(multiversxProxyUrl)
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

	multiversxProxyUrl := "url"

	interval := &api.Interval{
		Start: 4,
		End:   45,
	}
	numNonces := interval.End - interval.Start + 1

	multiversxEndPointCallsCt := uint64(0)
	processHyperBlocksCt := uint64(0)
	encodeCt := uint64(0)

	multiversxEndPoint := &apiMocks.MultiversxHyperBlockEndPointStub{
		GetHyperBlockCalled: func(path string) (*api.MultiversxHyperBlockApiResponse, error) {
			atomic.AddUint64(&multiversxEndPointCallsCt, 1)

			nonceFromRequest := getNonceFromRequest(t, path)
			requireNonceInInterval(t, nonceFromRequest, interval)
			require.Equal(t, fmt.Sprintf("%s%s/%d", multiversxProxyUrl, hyperBlockPathByNonce, nonceFromRequest), path)

			return &api.MultiversxHyperBlockApiResponse{
				Data: api.MultiversxHyperBlockApiResponsePayload{
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

	facade, _ := NewHyperBlockFacade(multiversxProxyUrl, encoder, multiversxEndPoint, processor)

	expectedEncodedHyperBlocks := make([][]byte, 0)
	for nonce := interval.Start; nonce <= interval.End; nonce++ {
		encodedBlock := []byte(fmt.Sprintf("encodedBlock%d", nonce))
		expectedEncodedHyperBlocks = append(expectedEncodedHyperBlocks, encodedBlock)
	}
	options := config.HyperBlocksQueryOptions{
		BatchSize: 10,
	}
	blocks, err := facade.GetHyperBlocksByInterval(interval, options)
	require.Nil(t, err)
	require.Equal(t, &api.CovalentHyperBlocksApiResponse{
		Data:  expectedEncodedHyperBlocks,
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}, blocks)

	require.Equal(t, numNonces, multiversxEndPointCallsCt)
	require.Equal(t, numNonces, processHyperBlocksCt)
	require.Equal(t, numNonces, encodeCt)
}

func TestHyperBlockFacade_GetHyperBlocksByInterval_CouldNotFetchAllHyperBlocks_ExpectError(t *testing.T) {
	t.Parallel()

	multiversxProxyUrl := "url"

	interval := &api.Interval{
		Start: 4,
		End:   45,
	}

	multiversxEndPointCallsCt := uint64(0)
	processHyperBlocksCt := uint64(0)
	encodeCt := uint64(0)

	expectedErr := errors.New("could not get hyper block")
	invalidNonce := uint64(40)
	multiversxEndPoint := &apiMocks.MultiversxHyperBlockEndPointStub{
		GetHyperBlockCalled: func(path string) (*api.MultiversxHyperBlockApiResponse, error) {
			atomic.AddUint64(&multiversxEndPointCallsCt, 1)

			nonceFromRequest := getNonceFromRequest(t, path)
			requireNonceInInterval(t, nonceFromRequest, interval)
			require.Equal(t, fmt.Sprintf("%s%s/%d", multiversxProxyUrl, hyperBlockPathByNonce, nonceFromRequest), path)

			if nonceFromRequest == invalidNonce {
				return nil, expectedErr
			}

			return &api.MultiversxHyperBlockApiResponse{
				Data: api.MultiversxHyperBlockApiResponsePayload{
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

	facade, _ := NewHyperBlockFacade(multiversxProxyUrl, encoder, multiversxEndPoint, processor)

	expectedEncodedHyperBlocks := make([][]byte, 0)
	for nonce := interval.Start; nonce <= interval.End; nonce++ {
		encodedBlock := []byte(fmt.Sprintf("encodedBlock%d", nonce))
		expectedEncodedHyperBlocks = append(expectedEncodedHyperBlocks, encodedBlock)
	}

	options := config.HyperBlocksQueryOptions{
		BatchSize: 10,
	}
	blocks, err := facade.GetHyperBlocksByInterval(interval, options)
	require.Nil(t, blocks)
	require.True(t, strings.Contains(err.Error(), errCouldNotGetHyperBlock.Error()))
	require.True(t, strings.Contains(err.Error(), expectedErr.Error()))
	require.True(t, strings.Contains(err.Error(), fmt.Sprintf("%s%s/%d", multiversxProxyUrl, hyperBlockPathByNonce, invalidNonce)))
	require.True(t, strings.Contains(err.Error(), fmt.Sprintf("%d", maxRequestsRetrial)))

	require.Equal(t, uint64(41)+maxRequestsRetrial, multiversxEndPointCallsCt) // 41 calls in [4,40] + maxRequestsRetrial
	require.Equal(t, uint64(41), processHyperBlocksCt)                         // 41 calls in [4,40]
	require.Equal(t, uint64(41), encodeCt)                                     // 41 calls in [4,40]
}

func TestHyperBlockFacade_GetHyperBlocksByInterval_GetHyperBlockAfterNumRetrials(t *testing.T) {
	t.Parallel()

	multiversxProxyUrl := "url"

	interval := &api.Interval{
		Start: 4,
		End:   45,
	}
	numNonces := interval.End - interval.Start + 1

	multiversxEndPointCallsCt := uint64(0)
	processHyperBlocksCt := uint64(0)
	encodeCt := uint64(0)

	expectedErr := errors.New("could not get hyper block")
	invalidNonce := uint64(40)
	numRetrials := 0
	maxNumRetrials := maxRequestsRetrial / 2
	multiversxEndPoint := &apiMocks.MultiversxHyperBlockEndPointStub{
		GetHyperBlockCalled: func(path string) (*api.MultiversxHyperBlockApiResponse, error) {
			atomic.AddUint64(&multiversxEndPointCallsCt, 1)

			nonceFromRequest := getNonceFromRequest(t, path)
			requireNonceInInterval(t, nonceFromRequest, interval)
			require.Equal(t, fmt.Sprintf("%s%s/%d", multiversxProxyUrl, hyperBlockPathByNonce, nonceFromRequest), path)

			if nonceFromRequest == invalidNonce && numRetrials < maxNumRetrials {
				numRetrials++
				return nil, expectedErr
			}

			return &api.MultiversxHyperBlockApiResponse{
				Data: api.MultiversxHyperBlockApiResponsePayload{
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

	facade, _ := NewHyperBlockFacade(multiversxProxyUrl, encoder, multiversxEndPoint, processor)

	expectedEncodedHyperBlocks := make([][]byte, 0)
	for nonce := interval.Start; nonce <= interval.End; nonce++ {
		encodedBlock := []byte(fmt.Sprintf("encodedBlock%d", nonce))
		expectedEncodedHyperBlocks = append(expectedEncodedHyperBlocks, encodedBlock)
	}
	options := config.HyperBlocksQueryOptions{
		BatchSize: 10,
	}
	blocks, err := facade.GetHyperBlocksByInterval(interval, options)
	require.Nil(t, err)
	require.Equal(t, &api.CovalentHyperBlocksApiResponse{
		Data:  expectedEncodedHyperBlocks,
		Error: "",
		Code:  api.ReturnCodeSuccess,
	}, blocks)

	require.Equal(t, numNonces+uint64(numRetrials), multiversxEndPointCallsCt)
	require.Equal(t, numNonces, processHyperBlocksCt)
	require.Equal(t, numNonces, encodeCt)
}

func TestHyperBlockFacade_GetHyperBlocksByInterval_IntervalEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("invalid nonces interval, should return error", func(t *testing.T) {
		t.Parallel()

		facade, _ := NewHyperBlockFacade("url",
			&mock.AvroEncoderStub{},
			&apiMocks.MultiversxHyperBlockEndPointStub{},
			&mock.HyperBlockProcessorStub{},
		)

		interval := &api.Interval{
			Start: 10,
			End:   9,
		}
		options := config.HyperBlocksQueryOptions{
			BatchSize: 10,
		}
		blocks, err := facade.GetHyperBlocksByInterval(interval, options)
		require.Nil(t, blocks)
		require.Equal(t, errInvalidNoncesInterval, err)
	})

	t.Run("invalid batch size, should return error", func(t *testing.T) {
		t.Parallel()

		facade, _ := NewHyperBlockFacade("url",
			&mock.AvroEncoderStub{},
			&apiMocks.MultiversxHyperBlockEndPointStub{},
			&mock.HyperBlockProcessorStub{},
		)

		interval := &api.Interval{
			Start: 10,
			End:   12,
		}
		options := config.HyperBlocksQueryOptions{
			BatchSize: 0,
		}
		blocks, err := facade.GetHyperBlocksByInterval(interval, options)
		require.Nil(t, blocks)
		require.Equal(t, errInvalidBatchSize, err)
	})

	t.Run("start interval = end interval, should only return one encoded hyper block", func(t *testing.T) {
		t.Parallel()

		encodedHyperBlock := []byte("encodedHyperBlock")
		encoder := &mock.AvroEncoderStub{
			EncodeCalled: func(record avro.AvroRecord) ([]byte, error) {
				return encodedHyperBlock, nil
			},
		}
		facade, _ := NewHyperBlockFacade("url",
			encoder,
			&apiMocks.MultiversxHyperBlockEndPointStub{},
			&mock.HyperBlockProcessorStub{},
		)

		interval := &api.Interval{
			Start: 10,
			End:   10,
		}
		options := config.HyperBlocksQueryOptions{
			BatchSize: 10,
		}
		blocks, err := facade.GetHyperBlocksByInterval(interval, options)
		require.Nil(t, err)
		require.Equal(t, &api.CovalentHyperBlocksApiResponse{
			Data:  [][]byte{encodedHyperBlock},
			Error: "",
			Code:  api.ReturnCodeSuccess,
		}, blocks)
	})
}
