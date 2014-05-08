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

// The PackIndexV2 type represents the improved version 2 pack index format. The
// v2 format starts out with a magic number that is interpreted by v1 format
// decoders as an invalid fanout value. Following the magic signature is an
// actual 4-byte version number and a fan-out table identical to that of the v1
// format.
//
// Entries now consist of three parts: the SHA-1 hash of an object, the CRC32
// checksum of the packed entry data, and the offset into the pack. What
// distinguishes the v2 format from the v1 format is how entries are organized:
// all the SHA-1 hashes are packed together, as are the CRC32 checksums and the
// offsets themselves. In other words, three parallel arrays are used to
// represent entries for better cache locality.
//
// Lastly, the offsets are encoded so that an offset larger than 0x7fffffff
// actually points to another table of 64-bit offsets, so the v2 format supports
// pack files up to 16 EiB.
//
// For more information on the pack and pack index format, see:
// https://www.kernel.org/pub/software/scm/git/docs/technical/pack-format.txt
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

// EntryForSha1 returns a PackIndexEntry whose Sha1() value matches that of the
// given object. If the given object is not found in the pack index, nil is
// returned.
func (idx *PackIndexV2) EntryForSha1(object core.Sha1) PackIndexEntry {
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
		return nil
	}

	return packIndexV2Entry{idx: idx, pos: pos + lower}
}

type packIndexV2Entry struct {
	idx *PackIndexV2
	pos int
}

var _ PackIndexEntry = packIndexV2Entry{}

func (e packIndexV2Entry) Offset() int64 {
	return int64(e.idx.offsets[e.pos])
}

func (e packIndexV2Entry) Sha1() core.Sha1 {
	return e.idx.objectNames[e.pos]
}

func (e packIndexV2Entry) Crc32() core.Crc32 {
	return e.idx.crc32Checksums[e.pos]
}

// Entries returns a slice that represents entries in this pack index.
func (idx *PackIndexV2) Entries() []PackIndexEntry {
	entries := make([]PackIndexEntry, len(idx.objectNames))
	for i := 0; i < len(entries); i++ {
		entries[i] = packIndexV2Entry{idx: idx, pos: i}
	}

	return entries
}
