package covalent

import "errors"

// ErrNilPubKeyConverter signals that a pub key converter input parameter is nil
var ErrNilPubKeyConverter = errors.New("received nil input value: pub key converter")

// ErrNilAccountsAdapter signals that an accounts adapter input parameter is nil
var ErrNilAccountsAdapter = errors.New("received nil input value: accounts adapter")

// ErrBlockBodyAssertion signals that an error occurred when trying to assert BodyHandler interface of type block body
var ErrBlockBodyAssertion = errors.New("error asserting BodyHandler interface of type block body")
