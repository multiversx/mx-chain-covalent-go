package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/hyperBlock"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/stretchr/testify/require"
)

func TestElrondHyperBlockEndPoint_GetHyperBlock(t *testing.T) {
	t.Parallel()

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
	bodyResponse, errMarshal := json.Marshal(expectedElrondApiResponse)
	require.Nil(t, errMarshal)

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		client := &mock.HTTPClientStub{
			GetCalled: func(url string) (resp *http.Response, err error) {
				require.Equal(t, path, url)

				return &http.Response{
					Body:       ioutil.NopCloser(bytes.NewBuffer(bodyResponse)),
					StatusCode: http.StatusOK,
				}, nil
			},
		}

		elrondEndPoint, _ := NewElrondHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := elrondEndPoint.GetHyperBlock(path)
		require.Nil(t, err)
		require.Equal(t, expectedElrondApiResponse, hyperBlockApiResponse)
	})

	t.Run("could not get response from http client, should return error", func(t *testing.T) {
		t.Parallel()

		errHttpClient := errors.New("http client local err")
		client := &mock.HTTPClientStub{
			GetCalled: func(url string) (resp *http.Response, err error) {
				return nil, errHttpClient
			},
		}

		elrondEndPoint, _ := NewElrondHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := elrondEndPoint.GetHyperBlock(path)
		require.Nil(t, hyperBlockApiResponse)
		require.Equal(t, errHttpClient, err)
	})

	t.Run("could not read body, should return error and close body", func(t *testing.T) {
		t.Parallel()

		errReadBytes := errors.New("error reading bytes")
		wasReaderClosed := false
		body := &mock.ReadCloserStub{
			ReadCalled: func(p []byte) (n int, err error) {
				return 0, errReadBytes
			},
			CloseCalled: func() error {
				wasReaderClosed = true
				return nil
			},
		}
		client := &mock.HTTPClientStub{
			GetCalled: func(url string) (resp *http.Response, err error) {
				require.Equal(t, path, url)

				return &http.Response{
					Body:       body,
					StatusCode: http.StatusBadRequest,
				}, nil
			},
		}

		elrondEndPoint, _ := NewElrondHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := elrondEndPoint.GetHyperBlock(path)
		require.Nil(t, hyperBlockApiResponse)
		require.Equal(t, errReadBytes, err)
		require.True(t, wasReaderClosed)
	})

	t.Run("could not unmarshall response, should return error", func(t *testing.T) {
		t.Parallel()

		body := &mock.ReadCloserStub{
			ReadCalled: func(p []byte) (n int, err error) {
				return 0, io.EOF
			},
		}
		client := &mock.HTTPClientStub{
			GetCalled: func(url string) (resp *http.Response, err error) {
				require.Equal(t, path, url)

				return &http.Response{
					Body:       body,
					StatusCode: http.StatusBadRequest,
				}, nil
			},
		}

		elrondEndPoint, _ := NewElrondHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := elrondEndPoint.GetHyperBlock(path)
		require.Nil(t, hyperBlockApiResponse)
		require.NotNil(t, err)
	})

	t.Run("status code not ok, should return error", func(t *testing.T) {
		t.Parallel()

		client := &mock.HTTPClientStub{
			GetCalled: func(url string) (resp *http.Response, err error) {
				require.Equal(t, path, url)

				return &http.Response{
					Body:       ioutil.NopCloser(bytes.NewBuffer(bodyResponse)),
					StatusCode: http.StatusBadRequest,
				}, nil
			},
		}

		elrondEndPoint, _ := NewElrondHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := elrondEndPoint.GetHyperBlock(path)
		require.Nil(t, hyperBlockApiResponse)
		require.NotNil(t, err)
		require.True(t, strings.Contains(err.Error(), strconv.Itoa(http.StatusBadRequest)))
	})
}
