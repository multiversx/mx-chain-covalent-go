package api

import "errors"

var errNilHttpServer = errors.New("nil http server provided")

var errEmptyElrondProxyUrl = errors.New("empty proxy url provided")

var errNilHyperBlockFacade = errors.New("nil hyper block facade provided")

var errNilAvroEncoder = errors.New("nil avro encoder provided")

var errNilHyperBlockProcessor = errors.New("nil hyper block processor provided")

var errInvalidBlockNonce = errors.New("invalid block nonce")

var errInvalidBlockHash = errors.New("invalid block hash")
