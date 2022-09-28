package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/stretchr/testify/require"
)

func TestElrondHyperBlockEndPoint_GetHyperBlock(t *testing.T) {
	path := "path"
	expectedElrondApiResponse := &ElrondHyperBlockApiResponse{
		Data: ElrondHyperBlockApiResponsePayload{
			HyperBlock: hyperBlock.HyperBlock{
				Hash: "hash",
			},
		},
		Error: "",
		Code:  "success",
	}
	bodyResponse, err := json.Marshal(expectedElrondApiResponse)
	require.Nil(t, err)

	t.Run("should work", func(t *testing.T) {
		client := &mock.HTTPClientStub{
			GetCalled: func(url string) (resp *http.Response, err error) {
				require.Equal(t, path, url)

				return &http.Response{
					Body:       ioutil.NopCloser(bytes.NewBufferString(string(bodyResponse))),
					StatusCode: http.StatusOK,
				}, nil
			},
		}

		elrondEndPoint, _ := NewElrondHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := elrondEndPoint.GetHyperBlock(path)
		require.Nil(t, err)
		require.Equal(t, expectedElrondApiResponse, hyperBlockApiResponse)
	})

	t.Run("could not get response from http client", func(t *testing.T) {
		expectedErr := errors.New("http client local err")
		client := &mock.HTTPClientStub{
			GetCalled: func(url string) (resp *http.Response, err error) {
				return nil, expectedErr
			},
		}

		elrondEndPoint, _ := NewElrondHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := elrondEndPoint.GetHyperBlock(path)
		require.Nil(t, hyperBlockApiResponse)
		require.Equal(t, expectedErr, err)

	})

	t.Run("status code not ok", func(t *testing.T) {
		client := &mock.HTTPClientStub{
			GetCalled: func(url string) (resp *http.Response, err error) {
				require.Equal(t, path, url)

				return &http.Response{
					Body:       ioutil.NopCloser(bytes.NewBufferString(string(bodyResponse))),
					StatusCode: http.StatusBadRequest,
				}, nil
			},
		}

		elrondEndPoint, _ := NewElrondHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := elrondEndPoint.GetHyperBlock(path)
		require.Nil(t, hyperBlockApiResponse)
		require.NotNil(t, err)
		require.True(t, strings.Contains(err.Error(), fmt.Sprintf("status code: %d", http.StatusBadRequest)))
	})
}
