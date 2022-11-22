package facade

import "errors"

var errEmptyElrondProxyUrl = errors.New("empty proxy url provided")

var errNilAvroEncoder = errors.New("nil avro encoder provided")

var errNilHyperBlockProcessor = errors.New("nil hyper block processor provided")

var errNilHyperBlockEndpointHandler = errors.New("nil hyper block endpoint handler provided")

var errInvalidNoncesInterval = errors.New("invalid nonces interval")
