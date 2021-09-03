package ws_test

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/process/ws"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/stretchr/testify/require"
)

func TestWSSender_SendMessage_ExpectNoError(t *testing.T) {
	wss := &ws.WSSender{Conn: &mock.WSConnStub{}}
	require.Nil(t, wss.SendMessage([]byte{}))
}

func TestWSSender_SendMessage_ExpectError(t *testing.T) {
	wss := &ws.WSSender{
		Conn: &mock.WSConnStub{
			WriteMessageCalled: func(messageType int, data []byte) error {
				return errors.New("local err")
			},
		},
	}
	require.NotNil(t, wss.SendMessage([]byte{}))
}
