package core

import (
	"bytes"
	"fmt"
	"io"
)

// A Sha1 represents a 40-character hexadecimal SHA-1 hash checksum.
type Sha1 [20]byte

// Reader returns an io.Reader that yields the bytes that comprise the SHA-1
// checksum.
func (sha Sha1) Reader() io.Reader {
	return bytes.NewBuffer(sha[:])
}

// String returns this checksum in hexadecimal form.
func (sha Sha1) String() string {
	return fmt.Sprintf("%040x", [20]byte(sha))
}

// Split returns two string values, obtained by splitting the checksum at index
// n.
func (sha Sha1) Split(n int) (string, string) {
	s := sha.String()
	return s[:n], s[n:]
}

// IsEmpty returns true if this checksum is all zero.
func (sha Sha1) IsEmpty() bool {
	return bytes.Count(sha[:], []byte{0}) == 20
}

// Sha1FromString attempts to parse a 40-character hexadecimal string as a Sha1
// checksum. An error is returned if the string is not well-formed.
func Sha1FromString(s string) (sha Sha1, err error) {
	bytes := make([]byte, 40)
	if _, err := fmt.Sscanf(s, "%040x", &bytes); err != nil {
		return sha, err
	}
	copy(sha[:], bytes[:])
	return
}
