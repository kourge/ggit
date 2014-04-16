package format

import (
	"encoding/binary"

	"github.com/kourge/ggit/core"
)

type indexEntryHeaderV2 struct {
	CtimeSecs     uint32
	CtimeNanosecs uint32
	MtimeSecs     uint32
	MtimeNanosecs uint32
	Dev           uint32
	Ino           uint32
	Mode          core.GitMode
	Uid           uint32
	Gid           uint32
	FileSize      uint32
	Sha1          core.Sha1
	Flags         indexEntryHeaderV2Flag
}

func (hv2 indexEntryHeaderV2) IndexEntryHeaderSize() int {
	return binary.Size(hv2)
}

type indexEntryHeaderV2Flag uint16

func (f indexEntryHeaderV2Flag) AssumeValid() bool {
	return (f >> 15) == 1
}

func (f indexEntryHeaderV2Flag) Extended() bool {
	return (f >> 14) == 1
}

type indexEntryStage uint8

func (f indexEntryHeaderV2Flag) Stage() indexEntryStage {
	return indexEntryStage((f >> 12) & 2)
}

// 12-bit
func (f indexEntryHeaderV2Flag) NameLength() uint16 {
	return uint16(f) & 0xfff
}

////////////////////////////////////////////////////////////////////////////////

type indexEntryHeaderV3 struct {
	indexEntryHeaderV2
	V3Flags indexEntryHeaderV3Flag
}

func (hv3 indexEntryHeaderV3) IndexEntryHeaderSize() int {
	return binary.Size(hv3)
}

type indexEntryHeaderV3Flag uint16

func (f indexEntryHeaderV3Flag) SkipWorktree() bool {
	return (f >> 14) == 1
}

func (f indexEntryHeaderV3Flag) IntentToAdd() bool {
	return (f >> 13) == 1
}

////////////////////////////////////////////////////////////////////////////////

type indexEntryHeader interface {
	IndexEntryHeaderSize() int
}

type indexEntry struct {
	indexEntryHeader
	pathName string
}
