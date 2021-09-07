package covalent_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/elrond-vm-common/atomic"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestNewCovalentDataIndexer(t *testing.T) {
	tests := []struct {
		args        func() (processor covalent.DataHandler, server *http.Server)
		expectedErr error
	}{
		{
			args: func() (processor covalent.DataHandler, server *http.Server) {
				return nil, &http.Server{Addr: "localhost:22111"}
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
				return &mock.DataHandlerStub{}, &http.Server{Addr: "localhost:22112"}
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
			Addr: "localhost:21119",
		},
	)
	defer ci.Close()
	called1 := atomic.Flag{}
	called1.Unset()

	called2 := atomic.Flag{}
	called2.Unset()

	wss1 := &mock.WSConnStub{
		CloseCalled: func() error {
			called1.Set()
			return nil
		},
	}

	wss2 := &mock.WSConnStub{
		CloseCalled: func() error {
			called2.Set()
			return nil
		},
	}

	go ci.SetWSSender(nil)
	time.Sleep(time.Millisecond * covalent.RetrialTimeoutMS)
	require.False(t, called1.IsSet())
	require.False(t, called2.IsSet())

	go ci.SetWSSender(wss1)
	time.Sleep(time.Millisecond * covalent.RetrialTimeoutMS)
	require.False(t, called1.IsSet())
	require.False(t, called2.IsSet())

	go ci.SetWSSender(wss2)
	time.Sleep(time.Millisecond * covalent.RetrialTimeoutMS)
	require.True(t, called1.IsSet())
	require.False(t, called2.IsSet())
}

func TestCovalentIndexer_SetWSReceiver_SetTwoConsecutiveWebSockets_ExpectFirstOneClosed(t *testing.T) {
	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{},
		&http.Server{
			Addr: "localhost:21119",
		},
	)
	defer ci.Close()

	called1 := atomic.Flag{}
	called1.Unset()

	called2 := atomic.Flag{}
	called2.Unset()

	wss1 := &mock.WSConnStub{
		CloseCalled: func() error {
			called1.Set()
			return nil
		},
	}

	wss2 := &mock.WSConnStub{
		CloseCalled: func() error {
			called2.Set()
			return nil
		},
	}

	go ci.SetWSReceiver(nil)
	time.Sleep(time.Millisecond * covalent.RetrialTimeoutMS)
	require.False(t, called1.IsSet())
	require.False(t, called2.IsSet())

	go ci.SetWSReceiver(wss1)
	time.Sleep(time.Millisecond * covalent.RetrialTimeoutMS)
	require.False(t, called1.IsSet())
	require.False(t, called2.IsSet())

	go ci.SetWSReceiver(wss2)
	time.Sleep(time.Millisecond * covalent.RetrialTimeoutMS)
	require.True(t, called1.IsSet())
	require.False(t, called2.IsSet())

	_ = ci.Close()
}

func TestCovalentIndexer_SaveBlock_ErrorProcessingData_ExpectPanic(t *testing.T) {
	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{
			ProcessDataCalled: func(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
				return nil, errors.New("local error")
			},
		},
		&http.Server{
			Addr: "localhost:3333",
		},
	)
	defer ci.Close()

	require.Panics(t, func() { ci.SaveBlock(nil) })
}

func TestCovalentIndexer_SaveBlock_ErrorEncodingBlockRes_ExpectPanic(t *testing.T) {
	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{
			ProcessDataCalled: func(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
				return nil, nil
			},
		},
		&http.Server{
			Addr: "localhost:21119",
		},
	)
	defer ci.Close()

	require.Panics(t, func() { ci.SaveBlock(nil) })
}

func TestCovalentIndexer_SaveBlock_ExpectSuccess(t *testing.T) {
	blockRes := generateRandomValidBlockResult()

	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{
			ProcessDataCalled: func(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
				return blockRes, nil
			},
		},
		&http.Server{
			Addr: "localhost:21119",
		})
	defer ci.Close()

	wssCalled := atomic.Flag{}
	wssCalled.Unset()
	wss := &mock.WSConnStub{
		WriteMessageCalled: func(messageType int, data []byte) error {
			wssCalled.Set()
			return nil
		},
	}

	wsrCalled := atomic.Flag{}
	wsrCalled.Unset()
	wsr := &mock.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wsrCalled.Set()
			return websocket.BinaryMessage, blockRes.Block.Hash, nil
		},
	}

	go func() {
		ci.SaveBlock(nil)

		// Expect data is sent/received only after WSS & WSR are set
		require.True(t, wssCalled.IsSet())
		require.True(t, wsrCalled.IsSet())
	}()

	time.Sleep(time.Millisecond * 200)
	// Expect no data is sent/received until WSS & WSR are set
	require.False(t, wssCalled.IsSet())
	require.False(t, wsrCalled.IsSet())

	go ci.SetWSSender(wss)
	go ci.SetWSReceiver(wsr)
	time.Sleep(time.Millisecond * 200)
}

func TestCovalentIndexer_SaveBlock_WrongAcknowledgedDataFourTimes_ExpectSuccessAfterFourRetrials(t *testing.T) {
	blockRes := generateRandomValidBlockResult()

	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{
			ProcessDataCalled: func(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
				return blockRes, nil
			},
		},
		&http.Server{
			Addr: "localhost:21119",
		})
	defer ci.Close()

	wssCalledCt := 0
	wss := &mock.WSConnStub{
		WriteMessageCalled: func(messageType int, data []byte) error {
			wssCalledCt++
			return nil
		},
	}

	wsrCalledCt := 0
	wsr := &mock.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wsrCalledCt++
			// After 4 consecutive failed acknowledged messages, send the correct ack data
			if wsrCalledCt == 4 {
				return websocket.BinaryMessage, blockRes.Block.Hash, nil
			}
			return websocket.BinaryMessage, []byte{0x1}, nil
		},
	}

	go func() {
		ci.SaveBlock(nil)

		// Expect data is sent/received 4 times (until a correct ack msg is sent) after WSS & WSR are set
		require.Equal(t, wssCalledCt, 4)
		require.Equal(t, wsrCalledCt, 4)
	}()

	time.Sleep(time.Millisecond * 200)
	// Expect no data is sent/received until WSS & WSR are set
	require.Equal(t, wssCalledCt, 0)
	require.Equal(t, wsrCalledCt, 0)

	go ci.SetWSSender(wss)
	go ci.SetWSReceiver(wsr)
	time.Sleep(time.Millisecond * 200)
}

func TestCovalentIndexer_SaveBlock_ErrorAcknowledgeData_ReconnectedWSR_ExpectMessageResent(t *testing.T) {
	blockRes := generateRandomValidBlockResult()

	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{
			ProcessDataCalled: func(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
				return blockRes, nil
			},
		},
		&http.Server{
			Addr: "localhost:21119",
		})
	defer ci.Close()

	wssCalledCt := 0
	wss := &mock.WSConnStub{
		WriteMessageCalled: func(messageType int, data []byte) error {
			wssCalledCt++
			return nil
		},
	}

	wsrCalledCt := 0
	wsr := &mock.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wsrCalledCt++
			return 0, nil, errors.New("read message error")
		},
	}

	wsrReconnectedCalledCt := 0
	go func() {
		ci.SaveBlock(nil)

		require.Equal(t, wssCalledCt, 2)
		require.Equal(t, wsrCalledCt, 1)
		require.Equal(t, wsrReconnectedCalledCt, 1)
	}()

	time.Sleep(time.Millisecond * 200)
	require.Equal(t, wssCalledCt, 0)
	require.Equal(t, wsrCalledCt, 0)
	require.Equal(t, wsrReconnectedCalledCt, 0)

	go ci.SetWSSender(wss)
	go ci.SetWSReceiver(wsr)
	time.Sleep(time.Millisecond * 200)

	wsrReconnected := &mock.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wsrReconnectedCalledCt++
			return websocket.BinaryMessage, blockRes.Block.Hash, nil
		},
	}

	go ci.SetWSReceiver(wsrReconnected)
	time.Sleep(time.Millisecond * 200)
}

func TestCovalentIndexer_SaveBlock_WrongAcknowledgeThreeTimes_ErrorSendingBlockTwoTimes_ExpectSuccessAfterNewWSSConnection(t *testing.T) {
	blockRes := generateRandomValidBlockResult()

	ci, _ := covalent.NewCovalentDataIndexer(
		&mock.DataHandlerStub{
			ProcessDataCalled: func(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
				return blockRes, nil
			},
		},
		&http.Server{
			Addr: "localhost:21119",
		})
	defer ci.Close()

	wssCalledCt1 := 0
	wss1 := &mock.WSConnStub{
		WriteMessageCalled: func(messageType int, data []byte) error {
			wssCalledCt1++
			if wssCalledCt1 == 2 {
				return errors.New("write message error")
			}
			return nil
		},
	}

	wsrCalledCt1 := 0
	wsr1 := &mock.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wsrCalledCt1++
			if wsrCalledCt1 == 3 {
				return websocket.BinaryMessage, blockRes.Block.Hash, nil
			}
			return websocket.BinaryMessage, []byte{0x1}, nil
		},
	}

	wss2Called := atomic.Flag{}
	wss2Called.Unset()

	go func() {
		ci.SaveBlock(nil)

		require.Equal(t, wssCalledCt1, 2)
		require.Equal(t, wsrCalledCt1, 3)
		require.True(t, wss2Called.IsSet(), true)
	}()

	time.Sleep(time.Millisecond * 200)
	require.Equal(t, wssCalledCt1, 0)
	require.Equal(t, wsrCalledCt1, 0)
	require.False(t, wss2Called.IsSet())

	go ci.SetWSSender(wss1)
	go ci.SetWSReceiver(wsr1)
	time.Sleep(time.Millisecond * 500)

	wss2 := &mock.WSConnStub{
		WriteMessageCalled: func(messageType int, data []byte) error {
			wss2Called.Set()
			return nil
		},
	}

	go ci.SetWSSender(wss2)
	time.Sleep(time.Millisecond * 500)
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
