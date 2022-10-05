package transactions

import "errors"

var errNilLogProcessor = errors.New("nil log processor provided")

var errNilReceiptProcessor = errors.New("nil receipt processor provided")
