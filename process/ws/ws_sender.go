package ws

import (
	"io"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/gorilla/websocket"
)

var log = logger.GetOrCreate("covalent")

type WsConn interface {
	io.Closer
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type WsSender struct {
	Conn WsConn
}

func (wss *WsSender) SendMessage(data []byte) {
	err := wss.Conn.WriteMessage(websocket.TextMessage, data) //TODO: CHANGE TO BINARY DATA
	if err != nil {
		log.Error("could not send message", "message", data, "error", err)
	}
}
