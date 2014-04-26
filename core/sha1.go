package core

import (
	"bytes"
	"fmt"
	"io"
)

// A Sha1 represents a 40-character hexadecimal SHA-1 hash checksum.
type Sha1 [20]byte

var _ Encoder = Sha1{}

// An empty, all-zero SHA-1 checksum.
var EmptySha1 = Sha1{}

// Reader returns an io.Reader that yields the bytes that comprise the SHA-1
// checksum.
func (sha Sha1) Reader() io.Reader {
	return bytes.NewBuffer(sha[:])
}

// String returns this checksum in form of a 40-digit, left zero-padded
// hexadecimal number.
func (sha Sha1) String() string {
	return fmt.Sprintf("%040x", [20]byte(sha))
}

func (sha Sha1) GoString() string {
	return fmt.Sprintf("Sha1{%s}", sha)
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

// Compare returns an integer comparing the two checksums lexicographically.
// The result will be 0 if sha==other, -1 if sha < other, and +1 if sha < other.
func (sha Sha1) Compare(other Sha1) int {
	return bytes.Compare(sha[:], other[:])
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

// Sha1FromByteSlice converts an arbitrarily-sized byte slice into a Sha1
// checksum by copying into a 20-byte array. If the slice's size is lesser than
// 20 bytes, then the remainder of the array is padded with zero bytes. If the
// slice's size is greater than 20 bytes, anything past the 20th byte is
// ignored.
func Sha1FromByteSlice(slice []byte) Sha1 {
	sha := Sha1{}
	copy(sha[:], slice[:])
	return sha
}
