package core

import (
	"io"
	"io/ioutil"
	"strings"
)

// sha1field wraps a Sha1 type to enable encoding and decoding SHA-1 checksums
// to and from their hexadecimal representations.
type sha1field struct {
	Sha1
}

var _ EncodeDecoder = &sha1field{}

// Reader returns an io.Reader that yields a string that represents the wrapped
// checksum in a 40-digit, left zero-padded hexadecimal number.
func (s *sha1field) Reader() io.Reader {
	return strings.NewReader(s.Sha1.String())
}

// Decode treats the given reader as a byte sequence representing a human-
// readable hexadecimal representation of a SHA-1 checksum and attempts to
// decode it.
func (s *sha1field) Decode(reader io.Reader) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	sha1, err := Sha1FromString(string(bytes))
	if err != nil {
		return err
	}

	s.Sha1 = sha1
	return nil
}
