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

	"github.com/multiversx/mx-chain-covalent-go/hyperBlock"
	"github.com/multiversx/mx-chain-covalent-go/testscommon/mock"
	"github.com/stretchr/testify/require"
)

func TestNewMultiversxHyperBlockEndPoint(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		endPoint, err := NewMultiversxHyperBlockEndPoint(&mock.HTTPClientStub{})
		require.NotNil(t, endPoint)
		require.Nil(t, err)
	})

	t.Run("nil http client, should return error", func(t *testing.T) {
		t.Parallel()

		endPoint, err := NewMultiversxHyperBlockEndPoint(nil)
		require.Nil(t, endPoint)
		require.Equal(t, errNilHttpServer, err)
	})
}

func TestMultiversxHyperBlockEndPoint_GetHyperBlock(t *testing.T) {
	t.Parallel()

	path := "path"
	expectedMultiversxApiResponse := &MultiversxHyperBlockApiResponse{
		Data: MultiversxHyperBlockApiResponsePayload{
			HyperBlock: hyperBlock.HyperBlock{
				Hash: "hash",
			},
		},
		Error: "",
		Code:  "success",
	}
	bodyResponse, errMarshal := json.Marshal(expectedMultiversxApiResponse)
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

		multiversxEndPoint, _ := NewMultiversxHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := multiversxEndPoint.GetHyperBlock(path)
		require.Nil(t, err)
		require.Equal(t, expectedMultiversxApiResponse, hyperBlockApiResponse)
	})

	t.Run("close response body failed, should work anyway and return no error", func(t *testing.T) {
		t.Parallel()

		reader := ioutil.NopCloser(bytes.NewBuffer(bodyResponse))
		errClosingBody := errors.New("error closing body")
		body := &mock.ReadCloserStub{
			ReadCalled: reader.Read,
			CloseCalled: func() error {
				return errClosingBody
			},
		}
		client := &mock.HTTPClientStub{
			GetCalled: func(url string) (resp *http.Response, err error) {
				require.Equal(t, path, url)

				return &http.Response{
					Body:       body,
					StatusCode: http.StatusOK,
				}, nil
			},
		}

		multiversxEndPoint, _ := NewMultiversxHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := multiversxEndPoint.GetHyperBlock(path)
		require.Nil(t, err)
		require.Equal(t, expectedMultiversxApiResponse, hyperBlockApiResponse)
	})

	t.Run("could not get response from http client, should return error", func(t *testing.T) {
		t.Parallel()

		errHttpClient := errors.New("http client local err")
		client := &mock.HTTPClientStub{
			GetCalled: func(url string) (resp *http.Response, err error) {
				return nil, errHttpClient
			},
		}

		multiversxEndPoint, _ := NewMultiversxHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := multiversxEndPoint.GetHyperBlock(path)
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

		multiversxEndPoint, _ := NewMultiversxHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := multiversxEndPoint.GetHyperBlock(path)
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

		multiversxEndPoint, _ := NewMultiversxHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := multiversxEndPoint.GetHyperBlock(path)
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

		multiversxEndPoint, _ := NewMultiversxHyperBlockEndPoint(client)
		hyperBlockApiResponse, err := multiversxEndPoint.GetHyperBlock(path)
		require.Nil(t, hyperBlockApiResponse)
		require.NotNil(t, err)
		require.True(t, strings.Contains(err.Error(), strconv.Itoa(http.StatusBadRequest)))
	})
}
