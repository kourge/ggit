package util

import (
	"errors"
	"io"
	"math/big"

	"github.com/kourge/ggit/core"
)

// VariableSize is a big Int wrapper type that can decode the variable length
// format used by Git pack files.
//
// A variable size is represented by a series of bytes. In each byte, the most
// significant bit serves as a boolean that signals whether another byte
// follows. The remaining lower seven bits constitute the actual value. Let the
// first byte read be called "size0", and each subsequent nth byte be called
// "sizeN". Suppose a total of T bytes were read. Then the resulting value is
// a (T*7)-bit integer comprised of size0 as its least significant part, and
// sizeN as its most significant part.
type VariableSize struct {
	*big.Int
}

var _ core.Decoder = &VariableSize{}

func NewVariableSize(i int64) *VariableSize {
	return &VariableSize{big.NewInt(i)}
}

func (v *VariableSize) Decode(reader io.Reader) error {
	nextByte := make([]byte, 1)
	sizeByte := sevenBitsZeroMore
	nth := 0
	for sizeByte.More() {
		if n, err := reader.Read(nextByte); err != nil {
			return err
		} else if n != 1 {
			return errors.New("varsize: could not read next size byte")
		}
		sizeByte = sevenBits(nextByte[0])

		sizeChunk := sizeByte.ToBig()
		for i := 0; i < nth; i++ {
			sizeChunk.Lsh(sizeChunk, 7)
		}

		v.Or(v.Int, sizeChunk)
		nth += 1
	}

	return nil
}

type sevenBits byte

const (
	sevenBitsZeroMore sevenBits = 128
)

func (s sevenBits) More() bool {
	return byte(s)>>7 == 1
}

func (s sevenBits) SizeN() byte {
	return byte(s) & 0x7f
}

func (s sevenBits) ToBig() *big.Int {
	return big.NewInt(int64(s.SizeN()))
}
