module github.com/ElrondNetwork/covalent-indexer-go

go 1.16

require (
	github.com/ElrondNetwork/elrond-go v1.3.49-0.20221101105926-2f349fe6063c
	github.com/ElrondNetwork/elrond-go-core v1.1.24-0.20221101091159-7dd46b891266
	github.com/ElrondNetwork/elrond-go-logger v1.0.9
	github.com/elodina/go-avro v0.0.0-20160406082632-0c8185d9a3ba
	github.com/gin-gonic/gin v1.8.1
	github.com/pelletier/go-toml v1.9.3
	github.com/stretchr/testify v1.7.1
	github.com/urfave/cli v1.22.10
)

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.41 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.41

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.41 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.41

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_4 v1.4.58 => github.com/ElrondNetwork/arwen-wasm-vm v1.4.58
