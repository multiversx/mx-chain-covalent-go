package ws

import (
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/gorilla/websocket"
)

var log = logger.GetOrCreate("covalent/websocket")

// WSSender handles sending binary data through websockets
type WSSender struct {
	Conn process.WSConn
}

// SendMessage sends a data buffer in binary format through a websocket
func (wss *WSSender) SendMessage(data []byte) error {
	return wss.Conn.WriteMessage(websocket.BinaryMessage, data)
}
