package covalent

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/process/ws"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/gorilla/websocket"
)

var log = logger.GetOrCreate("covalent")

type covalentIndexer struct {
	processor        DataHandler
	server           *http.Server
	wss              *ws.WSSender
	mutWss           sync.RWMutex
	wsr              *ws.WSReceiver
	mutWsr           sync.RWMutex
	newConnectionWSR chan struct{}
	newConnectionWSS chan struct{}
}

// NewCovalentDataIndexer creates a new instance of covalent data indexer, which implements Driver interface and
// converts protocol input data to covalent required data
func NewCovalentDataIndexer(processor DataHandler, server *http.Server) (*covalentIndexer, error) {
	if processor == nil {
		return nil, ErrNilDataHandler
	}
	if server == nil {
		return nil, ErrNilHTTPServer
	}
	ci := &covalentIndexer{
		processor: processor,
		server:    server,
	}
	ci.newConnectionWSR = make(chan struct{})
	ci.newConnectionWSS = make(chan struct{})

	go ci.start()

	return ci, nil
}

func (c *covalentIndexer) SetWSSender(wss *ws.WSSender) {
	c.mutWss.Lock()
	if c.wss != nil {
		_ = c.wss.Conn.Close()
	}
	c.wss = wss
	c.mutWss.Unlock()

	c.newConnectionWSS <- struct{}{}
	log.Error("AM SETAT WSS")
}

func (c *covalentIndexer) SetWSReceiver(wsr *ws.WSReceiver) {
	c.mutWsr.Lock()
	if c.wsr != nil {
		_ = c.wsr.Conn.Close()
	}
	c.wsr = wsr
	c.mutWsr.Unlock()

	c.newConnectionWSR <- struct{}{}
	log.Error("AM SETAT WSR")
}

func (c *covalentIndexer) start() {
	err := c.server.ListenAndServe()
	if err != nil {
		log.Error("could not initialize webserver", "error", err)
	}
}

func (c *covalentIndexer) sendBlockResultToCovalent(result *schema.BlockResult) {
	binaryData, err := utility.Encode(result)
	if err != nil {
		log.Error("could not encode block result to binary data", "error", err)
		return
	}

	c.sendWithRetrial(binaryData, result.Block.Hash)
}

func (c *covalentIndexer) waitForWSSConnection() {
	for {
		select {
		case <-c.newConnectionWSS:
			log.Warn("got a new WSS connection")
			return
		}
	}
}

func (c *covalentIndexer) waitForWSRConnection() {
	for {
		select {
		case <-c.newConnectionWSR:
			log.Warn("got a new WSS connection")
			return
		}
	}
}

func (c *covalentIndexer) sendWithRetrial(binaryData []byte, ackData []byte) {
	//resendTimeout := time.Duration(0)
	c.mutWss.RLock()
	wss := c.wss
	c.mutWss.RUnlock()

	c.mutWsr.RLock()
	wsr := c.wsr
	c.mutWsr.RUnlock()

	if wss == nil {
		c.waitForWSSConnection()
	}
	if wsr == nil {
		c.waitForWSRConnection()
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: //<-time.After(resendTimeout):
			c.mutWss.RLock()
			wss = c.wss
			c.mutWss.RUnlock()

			c.mutWsr.RLock()
			wsr = c.wsr
			c.mutWsr.RUnlock()

			log.Error("CEASUL INCA TICAIE")
			if wss == nil || wsr == nil {
				log.Error("WSR SI WSS INCA SUNT NILE")
			}
			if wss != nil && wsr != nil {
				log.Error("AM INTRAT UNDE WSS SI WSR NU SUNT NIL")
				msgTypeReceive, errSend, errReceive := c.sendMessageWithAck(binaryData, ackData, wss, wsr)
				if errSend == nil && errReceive == nil && msgTypeReceive == 0 {
					return
				}
			}
		}
	}
}

func (c *covalentIndexer) sendMessageWithAck(
	msgToSend []byte,
	ackData []byte,
	wss *ws.WSSender,
	wsr *ws.WSReceiver,
) (int, error, error) {
	errSend := wss.SendMessage(msgToSend)
	if errSend != nil {
		log.Error("GOT A  errSend, waiting for SENDER TO RECONNECT", "errSend", errSend)
		c.waitForWSSConnection()
	}

	msgType, receivedData, errReadData := wsr.ReadData()
	if errReadData != nil {
		log.Error("GOT A errReadData, waiting for RECEIVER TO RECONNECT", "errReadData", errReadData)
		c.waitForWSRConnection()
	}
	log.Error("PE sendMessageWithAck", "errSend", errSend, "msgType", msgType,
		"receivedData", hex.EncodeToString(receivedData), "errReadData", errReadData)

	if errSend == nil && errReadData == nil && msgType == websocket.BinaryMessage {
		log.Error("AM TRIMIS CU SUCCES, am primit", "hash", hex.EncodeToString(receivedData))
		if bytes.Compare(receivedData, ackData) == 0 {
			log.Error("HASH-URILE SUNT EGALE")
			return 0, nil, nil
		}
	}

	return msgType, errSend, errReadData
}

// SaveBlock saves the block info and converts it in order to be sent to covalent
func (c *covalentIndexer) SaveBlock(args *indexer.ArgsSaveBlockData) {
	blockResult, err := c.processor.ProcessData(args)
	if err != nil {
		log.Error("SaveBlock failed. Could not process block",
			"error", err, "headerHash", hex.EncodeToString(args.HeaderHash))
		panic("could not process block, please check log")
	}

	c.sendBlockResultToCovalent(blockResult)
}

// RevertIndexedBlock DUMMY
func (c covalentIndexer) RevertIndexedBlock(header data.HeaderHandler, body data.BodyHandler) {}

// SaveRoundsInfo DUMMY
func (c covalentIndexer) SaveRoundsInfo(roundsInfos []*indexer.RoundInfo) {}

// SaveValidatorsPubKeys DUMMY
func (c covalentIndexer) SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) {}

// SaveValidatorsRating DUMMY
func (c covalentIndexer) SaveValidatorsRating(indexID string, infoRating []*indexer.ValidatorRatingInfo) {
}

// SaveAccounts DUMMY
func (c covalentIndexer) SaveAccounts(blockTimestamp uint64, acc []data.UserAccountHandler) {}

// Close DUMMY
func (c covalentIndexer) Close() error {
	if c.wss != nil && c.wss.Conn != nil {
		err := c.wss.Conn.Close()
		log.LogIfError(err)
	}
	if c.server != nil {
		return c.server.Close()
	}
	return nil
}

// IsInterfaceNil DUMMY
func (c covalentIndexer) IsInterfaceNil() bool {
	return false
}
