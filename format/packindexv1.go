package format

import (
	"crypto/sha1"
	"encoding/binary"
	"io"
	"sort"

	"github.com/kourge/ggit/core"
)

type packIndexV1Header struct {
	Fanout [256]uint32
}

type packIndexV1Entry struct {
	Offset     uint32
	ObjectName core.Sha1
}

type packIndexV1EntryWrapper struct {
	entry *packIndexV1Entry
}

var _ PackIndexEntry = packIndexV1EntryWrapper{}

func (e packIndexV1EntryWrapper) Offset() int64 {
	return int64(e.entry.Offset)
}

func (e packIndexV1EntryWrapper) Sha1() core.Sha1 {
	return e.entry.ObjectName
}

func (e packIndexV1EntryWrapper) Crc32() core.Crc32 {
	return core.Crc32{}
}

// The PackIndexV1 type represents the original (version 1) pack index format.
// The v1 format is very simple: it consists of a first-level fan-out table
// followed by entries that each consist of a 32-bit pack offset and the 20-byte
// SHA-1 hash of the object that said offset points to.
//
// This format has two major weaknesses: due to the entry offset size, it cannot
// support pack files over 4 GiB in size, and the structure of an entry tends to
// poor cache locality because consecutive entries form a stream of
// zebra-striped offsets and hashes. In addition, this format can detect if the
// pack file has any data corruption, but not which individual pack entry is
// corrupted. All of these problems are addressed in the v2 format.
//
// For more information on the pack and pack index format, see:
// https://www.kernel.org/pub/software/scm/git/docs/technical/pack-format.txt
type PackIndexV1 struct {
	packIndexV1Header
	entries       []packIndexV1Entry
	packfileSha1  core.Sha1
	packIndexSha1 core.Sha1
}

var _ PackIndex = &PackIndexV1{}

func (idx *PackIndexV1) Decode(reader io.Reader) error {
	hash := sha1.New()
	r := io.TeeReader(reader, hash)

	header := &(idx.packIndexV1Header)
	if err := binary.Read(r, binary.BigEndian, header); err != nil {
		return err
	}

	entryCount := int(header.Fanout[255])

	idx.entries = make([]packIndexV1Entry, entryCount)
	if err := binary.Read(r, binary.BigEndian, idx.entries); err != nil {
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
func (idx *PackIndexV1) Size() int {
	return len(idx.entries)
}

// Objects returns a sorted slice of objects present in the pack index.
func (idx *PackIndexV1) Objects() []core.Sha1 {
	objects := make([]core.Sha1, len(idx.entries))
	for i, entry := range idx.entries {
		objects[i] = entry.ObjectName
	}
	return objects
}

// EntryForSha1 returns a PackIndexEntry whose Sha1() value matches that of the
// given object. If the given object is not found in the pack index, nil is
// returned.
func (idx *PackIndexV1) EntryForSha1(object core.Sha1) PackIndexEntry {
	lower := 0
	if object[0] != 0x00 {
		lower = int(idx.Fanout[int(object[0])-1])
	}
	upper := int(idx.Fanout[int(object[0])])
	entries := idx.entries[lower:upper]

	pos := sort.Search(len(entries), func(i int) bool {
		return entries[i].ObjectName.Compare(object) >= 0
	})

	if pos == len(entries) {
		return nil
	}

	return packIndexV1EntryWrapper{&idx.entries[pos + lower]}
}

// PosForSha1 returns the abstract position of the given object within the pack
// index. If the given object is not found in the pack index, the value
// PackIndexPosNotFound is returned.
func (idx *PackIndexV1) PosForSha1(object core.Sha1) PackIndexPos {
	lower := 0
	if object[0] != 0x00 {
		lower = int(idx.Fanout[int(object[0])-1])
	}
	upper := int(idx.Fanout[int(object[0])])
	entries := idx.entries[lower:upper]

	pos := sort.Search(len(entries), func(i int) bool {
		return entries[i].ObjectName.Compare(object) >= 0
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
func (idx *PackIndexV1) OffsetForPos(pos PackIndexPos) (offset int64, err error) {
	if pos < PackIndexPos(0) || int(pos) >= len(idx.entries) {
		return -1, ErrInvalidPackIndexPos
	}

	return int64(idx.entries[pos].Offset), nil
}

// Entries returns a slice that represents entries in this pack index.
func (idx *PackIndexV1) Entries() []PackIndexEntry {
	entries := make([]PackIndexEntry, len(idx.entries))
	for i := 0; i < len(entries); i++ {
		entries[i] = packIndexV1EntryWrapper{&idx.entries[i]}
	}

	return entries
}
