package core

import (
	"bytes"
	"fmt"
	"io"
)

// A Crc32 represents a 8-character hexadecimal CRC-32 hash checksum.
type Crc32 [4]byte

// Reader returns an io.Reader that yields the bytes that comprise the CRC-32
// checksum.
func (crc Crc32) Reader() io.Reader {
	return bytes.NewBuffer(crc[:])
}

// String returns this checksum in form of an 8-digit, left zero-padded
// hexadecimal number.
func (crc Crc32) String() string {
	return fmt.Sprintf("%08x", [4]byte(crc))
}

func (crc Crc32) GoString() string {
	return fmt.Sprintf("Crc32{%s}", crc)
}

// IsEmpty returns true if this checksum is all zero.
func (crc Crc32) IsEmpty() bool {
	return bytes.Count(crc[:], []byte{0}) == 4
}

// Compare returns an integer comparing the two checksums lexicographically.
// The result will be 0 if crc==other, -1 if crc < other, and +1 if crc < other.
func (crc Crc32) Compare(other Crc32) int {
	return bytes.Compare(crc[:], other[:])
}

// Crc32FromString attempts to parse a 8-character hexadecimal string as a Crc32
// checksum. An error is returned if the string is not well-formed.
func Crc32FromString(s string) (crc Crc32, err error) {
	bytes := make([]byte, 40)
	if _, err := fmt.Sscanf(s, "%08x", &bytes); err != nil {
		return crc, err
	}
	copy(crc[:], bytes[:])
	return
}

// Crc32FromByteSlice converts an arbitrarily-sized byte slice into a Crc32
// checksum by copying into a 4-byte array. If the slice's size is lesser than
// 4 bytes, then the remainder of the array is padded with zero bytes. If the
// slice's size is greater than 4 bytes, anything past the 4th byte is ignored.
func Crc32FromByteSlice(slice []byte) Crc32 {
	crc := Crc32{}
	copy(crc[:], slice[:])
	return crc
}
