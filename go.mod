module github.com/ElrondNetwork/covalent-indexer-go

go 1.16

require (
	github.com/ElrondNetwork/elrond-go v1.3.40
	github.com/ElrondNetwork/elrond-go-core v1.1.19
	github.com/ElrondNetwork/elrond-go-logger v1.0.7
	github.com/ElrondNetwork/elrond-vm-common v1.3.15
	github.com/elodina/go-avro v0.0.0-20160406082632-0c8185d9a3ba
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/stretchr/testify v1.7.1
	github.com/urfave/cli v1.22.9
)

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.41 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.41

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.41 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.41

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_4 v1.4.58 => github.com/ElrondNetwork/arwen-wasm-vm v1.4.58
