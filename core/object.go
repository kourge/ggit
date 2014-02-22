package core

// Object is the interface that groups the three defining attributes of a Git
// object: a type string, a size int, and the ability to encode and decode
// itself into and from a stream of bytes by implementing EncoderDecoder.
//
// Type should return a string that should be alphanumeric and must not contain
// a space.
//
// The stream of bytes must not include the header fields, which consist of the
// type string, a space, the content length, and the NULL character. In other
// words, this stream is everything after the NULL character.
type Object interface {
	Type() string
	Size() int
	EncodeDecoder
}
