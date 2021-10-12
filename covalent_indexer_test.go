package covalent_test

import (
	"context"
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
	defer func() {
		_ = ci.Close()
	}()

	called1 := atomic.Flag{}
	called2 := atomic.Flag{}

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
	defer func() {
		_ = ci.Close()
	}()

	called1 := atomic.Flag{}
	called2 := atomic.Flag{}

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
	defer func() {
		_ = ci.Close()
	}()

	require.Panics(t, func() { ci.SaveBlock(context.Background(), nil) })
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
	defer func() {
		_ = ci.Close()
	}()

	require.Panics(t, func() { ci.SaveBlock(context.Background(), nil) })
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
	defer func() {
		_ = ci.Close()
	}()

	wssCalled := atomic.Flag{}
	wss := &mock.WSConnStub{
		WriteMessageCalled: func(messageType int, data []byte) error {
			wssCalled.Set()
			return nil
		},
	}

	wsrCalled := atomic.Flag{}
	wsr := &mock.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wsrCalled.Set()
			return websocket.BinaryMessage, blockRes.Block.Hash, nil
		},
	}

	go func() {
		ci.SaveBlock(context.Background(), nil)

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
	defer func() {
		_ = ci.Close()
	}()

	wssCalledCt := atomic.Counter{}
	wss := &mock.WSConnStub{
		WriteMessageCalled: func(messageType int, data []byte) error {
			wssCalledCt.Increment()
			return nil
		},
	}

	wsrCalledCt := atomic.Counter{}
	wsr := &mock.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wsrCalledCt.Increment()
			// After 4 consecutive failed acknowledged messages, send the correct ack data
			if wsrCalledCt.Get() == 4 {
				return websocket.BinaryMessage, blockRes.Block.Hash, nil
			}
			return websocket.BinaryMessage, []byte{0x1}, nil
		},
	}

	go func() {
		ci.SaveBlock(context.Background(), nil)

		// Expect data is sent/received 4 times (until a correct ack msg is sent) after WSS & WSR are set
		require.Equal(t, wssCalledCt.Get(), int64(4))
		require.Equal(t, wsrCalledCt.Get(), int64(4))
	}()

	time.Sleep(time.Millisecond * 200)
	// Expect no data is sent/received until WSS & WSR are set
	require.Equal(t, wssCalledCt.Get(), int64(0))
	require.Equal(t, wsrCalledCt.Get(), int64(0))

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
	defer func() {
		_ = ci.Close()
	}()

	wssCalledCt := atomic.Counter{}
	wss := &mock.WSConnStub{
		WriteMessageCalled: func(messageType int, data []byte) error {
			wssCalledCt.Increment()
			return nil
		},
	}

	wsrCalledCt := atomic.Counter{}
	wsr := &mock.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wsrCalledCt.Increment()
			return 0, nil, errors.New("read message error")
		},
	}

	wsrReconnectedCalledCt := atomic.Counter{}
	go func() {
		ci.SaveBlock(context.Background(), nil)

		require.Equal(t, int64(2), wssCalledCt.Get())
		require.Equal(t, int64(1), wsrCalledCt.Get())
		require.Equal(t, int64(1), wsrReconnectedCalledCt.Get())
	}()

	time.Sleep(time.Millisecond * 200)
	require.Equal(t, int64(0), wssCalledCt.Get())
	require.Equal(t, int64(0), wsrCalledCt.Get())
	require.Equal(t, int64(0), wsrReconnectedCalledCt.Get())

	go ci.SetWSSender(wss)
	go ci.SetWSReceiver(wsr)
	time.Sleep(time.Millisecond * 200)

	wsrReconnected := &mock.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wsrReconnectedCalledCt.Increment()
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
	defer func() {
		_ = ci.Close()
	}()

	wssCalledCt1 := atomic.Counter{}
	wss1 := &mock.WSConnStub{
		WriteMessageCalled: func(messageType int, data []byte) error {
			wssCalledCt1.Increment()
			if wssCalledCt1.Get() == 2 {
				return errors.New("write message error")
			}
			return nil
		},
	}

	wsrCalledCt1 := atomic.Counter{}
	wsr1 := &mock.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wsrCalledCt1.Increment()
			if wsrCalledCt1.Get() == 3 {
				return websocket.BinaryMessage, blockRes.Block.Hash, nil
			}
			return websocket.BinaryMessage, []byte{0x1}, nil
		},
	}

	wss2Called := atomic.Flag{}

	go func() {
		ci.SaveBlock(context.Background(), nil)

		require.Equal(t, int64(2), wssCalledCt1.Get())
		require.Equal(t, int64(3), wsrCalledCt1.Get())
		require.True(t, wss2Called.IsSet())
	}()

	time.Sleep(time.Millisecond * 200)
	require.Equal(t, int64(0), wssCalledCt1.Get())
	require.Equal(t, int64(0), wsrCalledCt1.Get())
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
