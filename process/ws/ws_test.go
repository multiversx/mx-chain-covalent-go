package ws_test

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/process/ws"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
)

func TestWSSender_SendMessage(t *testing.T) {
	wss := &ws.WSSender{
		Conn: &mock.WSConnStub{
			WriteMessageCalled: func(messageType int, data []byte) error {
				return errors.New("local err")
			},
		},
	}
	wss.SendMessage([]byte{})
}
