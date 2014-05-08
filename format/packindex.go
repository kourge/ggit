package format

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/kourge/ggit/core"
)

var ErrInvalidPackIndexPos = errors.New("invalid pack index position")

// PackIndex represents the common set of operations expected of a pack index,
// regardless of its on-disk format. A pack index allows quick seeking of an
// object that has been packed into a pack file. It is still possible to find an
// object in a pack file without a pack index, although doing so is a costly
// task.
type PackIndex interface {
	core.Decoder

	// Size returns the number of objects present in the pack index.
	Size() int

	// Objects returns a sorted slice of objects present in the pack index.
	Objects() []core.Sha1

	// EntryForSha1 returns a PackIndexEntry whose Sha1() value matches that of
	// the given object. If the given object is not found in the pack index, nil
	// is returned.
	EntryForSha1(object core.Sha1) PackIndexEntry

	// Entries returns a slice that represents entries in the pack index.
	Entries() []PackIndexEntry
}

// PackIndexEntry represents an entry within a pack index. An entry is consisted
// of three things: an offset into the corresponding pack, the SHA-1 of the
// object that resides at that offset, and the CRC-32 of said object in its raw
// form in the pack. Note that only Offset and Sha1 are guaranteed to return
// meaningful values; early versions of pack indices did not include the CRC-32
// of objects.
type PackIndexEntry interface {
	Offset() int64
	Sha1() core.Sha1
	Crc32() core.Crc32
}

// PackIndexFromReader examines the given io.Reader, detects the right version
// of the pack index format, and decodes the stream.
func PackIndexFromReader(reader io.Reader) (idx PackIndex, err error) {
	r := bufio.NewReader(reader)

	if magic, err := r.Peek(4); err != nil {
		return nil, err
	} else if bytes.Compare(magic, packIndexV2HeaderMagic[:]) == 0 {
		idx = &PackIndexV2{}
	} else {
		idx = &PackIndexV1{}
	}

	err = idx.Decode(r)
	return idx, err
}
