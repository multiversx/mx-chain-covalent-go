package process

import "errors"

var errNilTransactionHandler = errors.New("nil transaction handler provided")

var errNilShardBlocksHandler = errors.New("nil shard blocks handler provided")
