package abi

import "io"

// SingleValue is the interface to be implemented by all "single value" types, with respect to:
// https://docs.multiversx.com/developers/data/serialization-overview
type SingleValue interface {
	EncodeNested(writer io.Writer) error
	EncodeTopLevel(writer io.Writer) error
	DecodeNested(reader io.Reader) error
	DecodeTopLevel(data []byte) error
}
