package covalent

import "errors"

// ErrNilPubKeyConverter signals that a pub key converter input parameter is nil
var ErrNilPubKeyConverter = errors.New("received nil input value: pub key converter")

// ErrNilAccountsAdapter signals that an accounts adapter input parameter is nil
var ErrNilAccountsAdapter = errors.New("received nil input value: accounts adapter")

// ErrBlockBodyAssertion signals that an error occurred when trying to assert BodyHandler interface of type block body
var ErrBlockBodyAssertion = errors.New("error asserting BodyHandler interface of type block body")

// ErrNilHasher signals that a nil hasher has been provided
var ErrNilHasher = errors.New("nil hasher provided")

// ErrNilMarshaller signals that a nil marshaller has been provided
var ErrNilMarshaller = errors.New("nil marshaller provided")

// ErrNilMiniBlockHandler signals that a nil mini block handler has been provided
var ErrNilMiniBlockHandler = errors.New("nil mini block handler provided")
