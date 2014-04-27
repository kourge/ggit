package format

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"hash/crc32"
	"io"
	"math"
	"math/big"

	"github.com/kourge/ggit/core"
	"github.com/kourge/ggit/util"
)

var BigMaxInt64 = big.NewInt(math.MaxInt64)

type packEntryHeader struct {
	packEntryFlag
}

type packEntry struct {
	packEntryHeader
	size  *util.VariableSize
	data  []byte
	sha1  core.Sha1
	crc32 core.Crc32
}

var _ core.Object = &packEntry{}

func (entry *packEntry) Type() string {
	switch entry.packEntryHeader.Type() {
	case PackedObjectCommit:
		return "commit"
	case PackedObjectTree:
		return "tree"
	case PackedObjectBlob:
		return "blob"
	case PackedObjectTag:
		return "tag"
	default:
		return "unknown"
	}
}

func (entry *packEntry) Size() int {
	return int(entry.size.Int64())
}

func (entry *packEntry) Reader() io.Reader {
	return bytes.NewReader(entry.data)
}

func (entry *packEntry) Decode(reader io.Reader) error {
	header := &(entry.packEntryHeader)
	if err := binary.Read(reader, binary.BigEndian, header); err != nil {
		return err
	}

	size0 := header.Size0()
	entry.size = util.NewVariableSize(0)
	if header.SizeExtension() {
		if err := entry.size.Decode(reader); err != nil {
			return err
		}
		entry.size.Lsh(entry.size.Int, 4)
		entry.size.Or(entry.size.Int, big.NewInt(int64(size0)))
	} else {
		entry.size.SetInt64(int64(size0))
	}

	switch entry.packEntryHeader.Type() {
	case PackedObjectCommit, PackedObjectTree, PackedObjectBlob, PackedObjectTag:
		if entry.size.Cmp(BigMaxInt64) > 0 {
			return Errorf("pack object size %d is too big", entry.size)
		}
		r, err := zlib.NewReader(reader)
		if err != nil {
			return err
		}

		buffer := new(bytes.Buffer)
		crc32 := crc32.NewIEEE()
		w := io.MultiWriter(buffer, crc32)
		if _, err := io.Copy(w, r); err != nil {
			return err
		}

		entry.crc32 = core.Crc32FromByteSlice(crc32.Sum(nil))
		entry.data = buffer.Bytes()
	case PackedObjectOfsDelta:
		return Errorf("ofs_delta objects not yet supported")
	case PackedObjectRefDelta:
		return Errorf("ref_delta objects not yet supported")
	default:
		return Errorf("%d is not a valid pack object type", entry.Type())
	}

	return nil
}

type packEntryFlag byte

func (f packEntryFlag) SizeExtension() bool {
	return byte(f)>>7 == 1
}

func (f packEntryFlag) Type() PackedObjectType {
	return PackedObjectType((byte(f) >> 4) & 0x7)
}

func (f packEntryFlag) Size0() byte {
	return byte(f) & 0xf
}

type PackedObjectType byte

const (
	PackedObjectNone     PackedObjectType = 0
	PackedObjectCommit   PackedObjectType = 1
	PackedObjectTree     PackedObjectType = 2
	PackedObjectBlob     PackedObjectType = 3
	PackedObjectTag      PackedObjectType = 4
	PackedObjectOfsDelta PackedObjectType = 6
	PackedObjectRefDelta PackedObjectType = 7
)
