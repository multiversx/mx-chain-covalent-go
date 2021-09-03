package covalent_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/ws"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/elrond-vm-common/atomic"
	"github.com/stretchr/testify/require"
)

func TestNewCovalentDataIndexer(t *testing.T) {
	tests := []struct {
		args        func() (processor covalent.DataHandler, server *http.Server)
		expectedErr error
	}{
		{
			args: func() (processor covalent.DataHandler, server *http.Server) {
				return nil, &http.Server{Addr: "localhost:8080"}
			},
			expectedErr: covalent.ErrNilDataHandler,
		},
		{
			args: func() (processor covalent.DataHandler, server *http.Server) {
				return &mock.DataHandlerStub{}, nil
			},
			expectedErr: covalent.ErrNilHTTPServer,
		},
		{
			args: func() (processor covalent.DataHandler, server *http.Server) {
				return &mock.DataHandlerStub{}, &http.Server{Addr: "localhost:8080"}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := covalent.NewCovalentDataIndexer(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestCovalentIndexer_SetWSSender_SetTwoConsecutiveWebSockets_ExpectFirstOneClosed(t *testing.T) {
	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{},
		&http.Server{
			Addr: "localhost:8080",
		})
	called := atomic.Flag{}
	called.Unset()
	wss := &ws.WSSender{
		Conn: &mock.WSConnStub{
			CloseCalled: func() error {
				called.Set()
				return nil
			},
		},
	}

	ci.SetWSSender(nil)
	require.False(t, called.IsSet())
	ci.SetWSSender(wss)
	require.False(t, called.IsSet())
	ci.SetWSSender(wss)
	require.True(t, called.IsSet())
}

func TestCovalentIndexer_sendBlockResultToCovalent_InvalidBlock_ExpectNoSentMessage(t *testing.T) {
	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{},
		&http.Server{
			Addr: "localhost:8080",
		})
	called := atomic.Flag{}
	called.Unset()
	wss := &ws.WSSender{
		Conn: &mock.WSConnStub{
			WriteMessageCalled: func(messageType int, data []byte) error {
				called.Set()
				return nil
			},
		},
	}
	ci.SetWSSender(wss)

	ci.SendBlockResultToCovalent(nil)
	require.False(t, called.IsSet())
}

func TestCovalentIndexer_sendBlockResultToCovalent_ExpectSentMessage(t *testing.T) {
	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{},
		&http.Server{
			Addr: "localhost:8080",
		})
	called := atomic.Flag{}
	called.Unset()
	wss := &ws.WSSender{
		Conn: &mock.WSConnStub{
			WriteMessageCalled: func(messageType int, data []byte) error {
				called.Set()
				return nil
			},
		},
	}

	ci.SetWSSender(wss)

	blockRes := generateRandomValidBlockResult()
	ci.SendBlockResultToCovalent(blockRes)
	require.True(t, called.IsSet())
}

func TestCovalentIndexer_SaveBlock_ExpectDataSent(t *testing.T) {
	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{
			ProcessDataCalled: func(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
				return generateRandomValidBlockResult(), nil
			},
		}, &http.Server{
			Addr: "localhost:3333",
		})

	called := atomic.Flag{}
	called.Unset()
	wss := &ws.WSSender{
		Conn: &mock.WSConnStub{
			WriteMessageCalled: func(messageType int, data []byte) error {
				called.Set()
				return nil
			},
		},
	}

	ci.SetWSSender(wss)
	ci.SaveBlock(nil)
	require.True(t, called.IsSet())
}

func TestCovalentIndexer_SaveBlock_ErrorProcessingData_ExpectDataNotSent(t *testing.T) {
	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{
			ProcessDataCalled: func(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
				return nil, errors.New("local error")
			},
		}, &http.Server{
			Addr: "localhost:3333",
		})

	called := atomic.Flag{}
	called.Unset()
	wss := &ws.WSSender{
		Conn: &mock.WSConnStub{
			WriteMessageCalled: func(messageType int, data []byte) error {
				called.Set()
				return nil
			},
		},
	}

	ci.SetWSSender(wss)
	ci.SaveBlock(nil)
	require.False(t, called.IsSet())
}

func generateRandomValidBlockResult() *schema.BlockResult {
	block := &schema.Block{
		Hash:          testscommon.GenerateRandomFixedBytes(32),
		StateRootHash: testscommon.GenerateRandomFixedBytes(32),
	}

	return &schema.BlockResult{
		Block: block,
	}
}
