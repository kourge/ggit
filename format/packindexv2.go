package format

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"io"
	"sort"

	"github.com/kourge/ggit/core"
)

var packIndexV2HeaderMagic = [4]byte{0xff, 0x74, 0x4f, 0x63}

type packIndexV2Header struct {
	Magic   [4]byte // == packIndexV2HeaderMagic
	Version uint32  // == 2
	Fanout  [256]uint32
}

type PackIndexV2 struct {
	packIndexV2Header
	objectNames    []core.Sha1 // sorted
	crc32Checksums []core.Crc32
	offsets        []uint32
	higherOffsets  []uint64
	packfileSha1   core.Sha1
	packIndexSha1  core.Sha1
}

var _ PackIndex = &PackIndexV2{}

func (idx *PackIndexV2) Decode(reader io.Reader) error {
	hash := sha1.New()
	r := io.TeeReader(reader, hash)

	header := &(idx.packIndexV2Header)
	if err := binary.Read(r, binary.BigEndian, header); err != nil {
		return err
	}

	if header.Magic != packIndexV2HeaderMagic {
		return errors.New("invalid pack index header")
	}
	if header.Version != 2 {
		return Errorf("unexpected pack index header version %d", header.Version)
	}

	entryCount := int(header.Fanout[255])

	idx.objectNames = make([]core.Sha1, entryCount)
	if err := binary.Read(r, binary.BigEndian, idx.objectNames); err != nil {
		return err
	}

	idx.crc32Checksums = make([]core.Crc32, entryCount)
	if err := binary.Read(r, binary.BigEndian, idx.crc32Checksums); err != nil {
		return err
	}

	idx.offsets = make([]uint32, entryCount)
	if err := binary.Read(r, binary.BigEndian, idx.offsets); err != nil {
		return err
	}

	higherOffsetCount := 0
	for _, offset := range idx.offsets {
		if (offset >> 31) == 1 {
			higherOffsetCount += 1
		}
	}

	idx.higherOffsets = make([]uint64, higherOffsetCount)
	if err := binary.Read(r, binary.BigEndian, idx.higherOffsets); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &idx.packfileSha1); err != nil {
		return err
	}

	if err := binary.Read(reader, binary.BigEndian, &idx.packIndexSha1); err != nil {
		return err
	}

	actualSha1 := core.Sha1FromByteSlice(hash.Sum(nil))
	if idx.packIndexSha1 != actualSha1 {
		return Errorf("pack index SHA-1 is %s, expected %s", actualSha1, idx.packIndexSha1)
	}

	return nil
}

// Size returns the number of objects present in the pack index.
func (idx *PackIndexV2) Size() int {
	return len(idx.objectNames)
}

// Objects returns a sorted slice of objects present in the pack index.
func (idx *PackIndexV2) Objects() []core.Sha1 {
	return idx.objectNames
}

// PosForSha1 returns the abstract position of the given object within the pack
// index. If the given object is not found in the pack index, the value
// PackIndexPosNotFound is returned.
func (idx *PackIndexV2) PosForSha1(object core.Sha1) PackIndexPos {
	lower := 0
	if object[0] != 0x00 {
		lower = int(idx.Fanout[int(object[0])-1])
	}
	upper := int(idx.Fanout[int(object[0])])
	entries := idx.objectNames[lower:upper]

	pos := sort.Search(len(entries), func(i int) bool {
		return entries[i].Compare(object) >= 0
	})

	if pos == len(entries) {
		return PackIndexPosNotFound
	}

	return PackIndexPos(pos + lower)
}

// OffsetForPos returns the byte offset of an object within the pack index's
// corresponding pack file, given that object's abstract position. If the given
// position is invalid, -1 is returned as the offset and the value
// ErrInvalidPackIndexPos is returned as the error.
func (idx *PackIndexV2) OffsetForPos(pos PackIndexPos) (offset int64, err error) {
	if pos < PackIndexPos(0) || int(pos) >= len(idx.offsets) {
		return -1, ErrInvalidPackIndexPos
	}

	return int64(idx.offsets[pos]), nil
}

// Crc32ForPos returns the CRC-32 checksum of an object within the pack index's
// corresponding pack file, given that object's abstract position. If the given
// position is invalid, 0 is returned as the checksum and the value
// ErrInvalidPackIndexPos is returned as the error.
func (idx *PackIndexV2) Crc32ForPos(pos PackIndexPos) (checksum core.Crc32, err error) {
	if pos < PackIndexPos(0) || int(pos) >= len(idx.offsets) {
		return core.Crc32{}, ErrInvalidPackIndexPos
	}

	return idx.crc32Checksums[pos], nil
}
