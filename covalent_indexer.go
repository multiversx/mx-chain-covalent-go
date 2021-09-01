package covalent

import (
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
		log.Error("could not initialize webserver", "error", err.Error())
	}
}

func (c *covalentIndexer) SendBlockToCovalent(result *schema.BlockResult) {
	data, err := utility.Encode(result)
	if err != nil {
		log.Error("could not encode block result", "error", err)
	}

	c.RLock()
	wss := c.wss
	c.RUnlock()

	if wss != nil {
		wss.SendMessage(data)
	}
}

// SaveBlock saves the block info and converts it in order to be sent to covalent
func (c *covalentIndexer) SaveBlock(args *indexer.ArgsSaveBlockData) {
	// TODO this function in future PRs
	// 1. Process data from args, format it according to avro schema
	blockResult, err := c.processor.ProcessData(args)
	if err != nil {

	}
	_ = blockResult

	c.SendBlockToCovalent(blockResult)
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
	return nil
}

// IsInterfaceNil DUMMY
func (c covalentIndexer) IsInterfaceNil() bool {
	return false
}
