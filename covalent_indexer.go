package covalent

import (
	"encoding/hex"
	"net/http"
	"sync"

	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/process/ws"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("covalent")

type covalentIndexer struct {
	processor DataHandler
	server    *http.Server
	wss       *ws.WSSender
	sync.RWMutex
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

	go ci.start()

	return ci, nil
}

func (c *covalentIndexer) SetWSSender(wss *ws.WSSender) {
	c.Lock()
	defer c.Unlock()
	if c.wss != nil {
		_ = c.wss.Conn.Close()
	}
	c.wss = wss
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

	c.RLock()
	wss := c.wss
	c.RUnlock()

	if wss != nil {
		//TODO: Handle error in next PR
		_ = wss.SendMessage(binaryData)
	}
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
func (c *covalentIndexer) RevertIndexedBlock(header data.HeaderHandler, body data.BodyHandler) {}

// SaveRoundsInfo DUMMY
func (c *covalentIndexer) SaveRoundsInfo(roundsInfos []*indexer.RoundInfo) {}

// SaveValidatorsPubKeys DUMMY
func (c *covalentIndexer) SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) {
}

// SaveValidatorsRating DUMMY
func (c *covalentIndexer) SaveValidatorsRating(indexID string, infoRating []*indexer.ValidatorRatingInfo) {
}

// SaveAccounts DUMMY
func (c *covalentIndexer) SaveAccounts(blockTimestamp uint64, acc []data.UserAccountHandler) {}

// Close DUMMY
func (c *covalentIndexer) Close() error {
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
func (c *covalentIndexer) IsInterfaceNil() bool {
	return false
}
