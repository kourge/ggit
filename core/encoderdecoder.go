package core

import (
	"io"
)

// Encoder is the interface that wraps the Reader method.
//
// Reader returns an io.Reader, which, when read, lazily yields a stream of
// bytes that represents the on-disk, serialized format of the underlying type.
type Encoder interface {
	Reader() io.Reader
}

// Decoder is the interface that wraps the Decode method.
//
// Decode mutates its pointer type and populates its underlying struct fields by
// interpreting the byte sequence yielded from reading the supplied io.Reader.
// Any existing field values may be clobbered, so it is customary to call Decode
// on the zero value of a type.
type Decoder interface {
	Decode(reader io.Reader) error
}

// EncodeDecoder is the interface that groups the Reader and Decode methods.
type EncodeDecoder interface {
	Encoder
	Decoder
}
