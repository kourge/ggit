package core

import (
	"bytes"
	"io"
)

// A Blob is a Git object type that stores the contents of files. Blobs are
// conceptually tuples of (Len, Bytes), where Bytes is a byte array and Len is
// the length of the byte array.
//
// Blobs carry no concept of a name. Trees hold that responsibility instead.
type Blob struct {
	Content []byte
}

var _ Object = &Blob{}

func (blob *Blob) Type() string {
	return "blob"
}

// Size returns the length of the blob content.
func (blob *Blob) Size() int {
	return len(blob.Content)
}

// Reader returns an io.Reader that yields the blob's internal byte array
// when read.
func (blob *Blob) Reader() io.Reader {
	return bytes.NewReader(blob.Content)
}

// Decode reads from an io.Reader into a byte array and replaces the blob's
// internal byte array with the one that was just read.
func (blob *Blob) Decode(reader io.Reader) error {
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(reader); err != nil {
		return err
	}
	blob.Content = buffer.Bytes()
	return nil
}
