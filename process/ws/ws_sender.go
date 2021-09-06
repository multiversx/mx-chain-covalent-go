package ws

import (
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/gorilla/websocket"
)

// WSSender handles sending binary data through websockets
type WSSender struct {
	Conn process.WSConn
}

// SendMessage sends a data buffer in binary format through a websocket
func (wss *WSSender) SendMessage(data []byte) error {
	return wss.Conn.WriteMessage(websocket.BinaryMessage, data)
}
