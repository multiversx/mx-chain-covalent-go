package covalent

import "errors"

// ErrNilPubKeyConverter signals that a pub key converter input parameter is nil
var ErrNilPubKeyConverter = errors.New("received nil input value: pub key converter")

// ErrNilAccountsAdapter signals that am accounts adapter input parameter is nil
var ErrNilAccountsAdapter = errors.New("received nil input value: accounts adapter")
