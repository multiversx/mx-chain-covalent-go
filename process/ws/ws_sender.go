package ws

import (
	"io"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/gorilla/websocket"
)

var log = logger.GetOrCreate("covalent")

type WSConn interface {
	io.Closer
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type WSSender struct {
	Conn WSConn
}

func (wss *WSSender) SendMessage(data []byte) {
	err := wss.Conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Error("could not send message", "message", data, "error", err)
	}
}
