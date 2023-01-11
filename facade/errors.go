package facade

import "errors"

var errEmptyMultiversxProxyUrl = errors.New("empty proxy url provided")

var errNilAvroEncoder = errors.New("nil avro encoder provided")

var errNilHyperBlockProcessor = errors.New("nil hyper block processor provided")

var errNilHyperBlockEndpointHandler = errors.New("nil hyper block endpoint handler provided")

var errInvalidNoncesInterval = errors.New("invalid nonces interval")

var errInvalidBatchSize = errors.New("received zero batch size")

var errCouldNotGetHyperBlock = errors.New("could not get hyper block")

var errCouldNotGetAllHyperBlocks = errors.New("could not get all hyper blocks")
