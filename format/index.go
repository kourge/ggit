package format

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"

	"github.com/kourge/ggit/core"
)

var _ = bytes.Count

var indexHeaderSignature = [4]byte{'D', 'I', 'R', 'C'}

type indexHeader struct {
	Signature    [4]byte // == indexHeaderSignature
	Version      uint32
	EntriesCount uint32
}

// An Index is the in-memory representation of a Git index file, which is
// a stored version of a repository's working tree.
type Index struct {
	version    uint32
	entries    []indexEntry
	extensions []indexExtension
	sha1       core.Sha1
}

var _ core.Decoder = &Index{}

// Decode takes a reader and treats the stream it yields as an index file and
// parses it. An error is returned if the stream forms an invalid index file.
// Otherwise nil is returned.
//
// See
// https://www.kernel.org/pub/software/scm/git/docs/technical/index-format.txt
// for more information on the file format.
func (idx *Index) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)
	// TODO: Verify SHA-1

	header, err := decodeIndexHeader(idx, r)
	if err != nil {
		return err
	}

	if err = decodeIndexEntries(idx, r, header.EntriesCount); err != nil {
		return err
	}

	if err = decodeIndexExtensionsAndSha1(idx, r); err != nil {
		return err
	}

	return nil
}

func decodeIndexHeader(idx *Index, r *bufio.Reader) (*indexHeader, error) {
	header := &indexHeader{}
	if err := binary.Read(r, binary.BigEndian, header); err != nil {
		return header, err
	}

	if header.Signature != indexHeaderSignature {
		return header, Errorf("%v is not a valid index header", header.Signature)
	}
	if header.Version != 2 {
		return header, Errorf("%d is not a valid index version", header.Version)
	} else {
		idx.version = header.Version
	}

	return header, nil
}

func decodeIndexEntries(idx *Index, r *bufio.Reader, n uint32) error {
	for i := uint32(0); i < n; i++ {
		var entryHeader indexEntryHeader

		v2EntryHeader := &indexEntryHeaderV2{}
		if err := binary.Read(r, binary.BigEndian, v2EntryHeader); err != nil {
			return err
		}
		entryHeader = v2EntryHeader

		if v2EntryHeader.Flags.Extended() {
			v3EntryHeader := &indexEntryHeaderV3{*v2EntryHeader, 0}
			if err := binary.Read(r, binary.BigEndian, &(v3EntryHeader.V3Flags)); err != nil {
				return err
			}
			entryHeader = v3EntryHeader
		}

		entry := &indexEntry{}
		totalRead := entryHeader.IndexEntryHeaderSize()
		if pathName, err := r.ReadString(0); err != nil {
			return err
		} else {
			entry.pathName = pathName[:len(pathName)-1]
			entry.indexEntryHeader = entryHeader

			totalRead += len(pathName)
			nearestMultiple := roundUpToNearestMultipleOfEight(totalRead)
			paddingSize := nearestMultiple - totalRead
			for i := 1; i <= paddingSize; i++ {
				if c, err := r.ReadByte(); err != nil {
					return err
				} else if c != 0 {
					return Errorf("path name padding is not null byte")
				}
			}
		}

		idx.entries = append(idx.entries, *entry)
	}

	return nil
}

func roundUpToNearestMultipleOfEight(i int) int {
	return (i + 8 - 1) & ^(8 - 1)
}

// All extensions, including cached tree and resolve undo extensions are
// currently simply parsed as a byte array. Their contents are neither
// interpreted nor validated.
func decodeIndexExtensionsAndSha1(idx *Index, r *bufio.Reader) error {
	for {
		if tentativeSha1, err := tryReadSha1(r); err == nil {
			idx.sha1 = tentativeSha1
			break
		}

		extensionHeader := &indexExtensionHeader{}
		if err := binary.Read(r, binary.BigEndian, extensionHeader); err != nil {
			return err
		}

		extensionData := make([]byte, extensionHeader.Size)
		if n, err := io.ReadFull(r, extensionData); err != nil {
			return err
		} else if n != int(extensionHeader.Size) {
			return Errorf(
				"failed to read %d bytes for extension %v, read %d bytes instead",
				extensionHeader.Size,
				extensionHeader.Signature,
				n,
			)
		}

		extension := indexExtension{*extensionHeader, extensionData}
		if !extension.Optional() {
			return Errorf(
				"cannot handle non-optional index extension %v",
				extension.Signature,
			)
		}
		idx.extensions = append(idx.extensions, extension)
	}

	return nil
}

type errSha1NotYetReached struct{}

func (e errSha1NotYetReached) Error() string {
	return "SHA-1 checksum not yet reached in index"
}

func tryReadSha1(r *bufio.Reader) (sha1 core.Sha1, err error) {
	sha1Size := 20
	if tentativeSha1, _ := r.Peek(sha1Size + 1); len(tentativeSha1) == sha1Size {
		return core.Sha1FromByteSlice(tentativeSha1), nil
	} else {
		return core.Sha1{}, errSha1NotYetReached{}
	}
}

// Pathnames returns an array of strings that are the path names of each entry
// within this index.
func (idx *Index) Pathnames() []string {
	pathnames := make([]string, len(idx.entries))
	for i, entry := range idx.entries {
		pathnames[i] = entry.pathName
	}
	return pathnames
}
