package api

import "errors"

var errNilHttpServer = errors.New("nil http server provided")

var errNilHyperBlockFacade = errors.New("nil hyper block facade provided")

var errInvalidBlockNonce = errors.New("invalid block nonce")

var errInvalidBlockHash = errors.New("invalid block hash")

var errInvalidHyperBlocksBatchSize = errors.New("invalid hyper blocks batch size")

var errMissingQueryParameter = errors.New("missing query parameter")
