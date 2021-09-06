package ws

import (
	"github.com/ElrondNetwork/covalent-indexer-go/process"
)

// WSReceiver handles sending receiving data through websockets
type WSReceiver struct {
	Conn process.WSConn
}

// ReadData reads a data buffer in binary format through a websocket
func (wss *WSReceiver) ReadData() (int, []byte, error) {
	return wss.Conn.ReadMessage()
}
