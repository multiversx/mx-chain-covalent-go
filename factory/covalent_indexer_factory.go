package factory

import (
	"log"
	"net/http"

	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/factory"
	ws2 "github.com/ElrondNetwork/covalent-indexer-go/process/ws"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// ArgsCovalentIndexerFactory holds all input dependencies required by covalent data indexer factory
// in order to create new instances
type ArgsCovalentIndexerFactory struct {
	Enabled          bool
	URL              string
	PubKeyConverter  core.PubkeyConverter
	Accounts         covalent.AccountsAdapter
	Hasher           hashing.Hasher
	Marshaller       marshal.Marshalizer
	ShardCoordinator process.ShardCoordinator
}

// CreateCovalentIndexer creates a new Driver instance of type covalent data indexer
func CreateCovalentIndexer(args *ArgsCovalentIndexerFactory) (covalent.Driver, error) {
	if check.IfNil(args.PubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}
	if check.IfNil(args.Accounts) {
		return nil, covalent.ErrNilAccountsAdapter
	}
	if check.IfNil(args.Hasher) {
		return nil, covalent.ErrNilHasher
	}
	if check.IfNil(args.Marshaller) {
		return nil, covalent.ErrNilMarshaller
	}

	argsDataProcessor := &factory.ArgsDataProcessor{
		PubKeyConvertor:  args.PubKeyConverter,
		Accounts:         args.Accounts,
		Hasher:           args.Hasher,
		Marshaller:       args.Marshaller,
		ShardCoordinator: args.ShardCoordinator,
	}

	dataProcessor, err := factory.CreateDataProcessor(argsDataProcessor)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	ci, err := covalent.NewCovalentDataIndexer(dataProcessor, server)
	if err != nil {
		return nil, err
	}

	router.HandleFunc("block", func(w http.ResponseWriter, r *http.Request) {
		// We'll need to define an Upgrader
		// this will require a Read and Write buffer size
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		// upgrade this connection to a WebSocket connection
		ws, errUpgrade := upgrader.Upgrade(w, r, nil)
		if errUpgrade != nil {
			log.Println(errUpgrade)
		}

		wss := &ws2.WSSender{
			Conn: ws,
		}
		ci.SetWSSender(wss)
	})

	return ci, nil
}
