package format

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kourge/ggit/core"
)

var (
	ErrPackFileAlreadyOpen  = errors.New("pack file already open")
	ErrObjectNotFoundInPack = errors.New("object not found in pack")
)

// A Pack is a set of objects that have been compressed into one file.
// Accessing any object stored in that file (called the "pack file") is sped
// up by a companion file called the "pack index".
type Pack struct {
	packPath string
	idxPath  string
	file     *os.File
	idx      PackIndex
}

// NewPack returns a Pack at the given path. The path can be a path to the
// pack file itself or the pack index file.
func NewPack(path string) *Pack {
	path = filepath.Clean(path)
	if dot := strings.LastIndex(path, "."); dot != -1 {
		path = path[:dot]
		return &Pack{packPath: path + ".pack", idxPath: path + ".idx"}
	}

	return &Pack{packPath: path}
}

// Open attempts to open this pack's pack file and its pack index. The pack file
// header and the pack index are verified by header checks and SHA-1 checksum
// matching.  The entire pack index is also parsed and loaded on Open.
func (p *Pack) Open() error {
	if p.file != nil {
		return ErrPackFileAlreadyOpen
	}

	file, err := os.Open(p.packPath)
	if err != nil {
		return err
	}

	err = verifyPack(file)
	if err != nil {
		return err
	}
	p.file = file

	// Reading a pack without an index is unsupported.
	if p.idxPath == "" {
		return Errorf("missing pack index file: %s", p.packPath)
	}

	// Load the corresponding index.
	if idxFile, err := os.Open(p.idxPath); err != nil {
		return err
	} else if idx, err := PackIndexFromReader(idxFile); err != nil {
		return err
	} else {
		p.idx = idx
		idxFile.Close()
	}

	return nil
}

var packHeaderSignature = [4]byte{'P', 'A', 'C', 'K'}

type packHeader struct {
	Signature   [4]byte // == packHeaderSignature
	Version     uint32
	ObjectCount uint32
}

func verifyPack(reader io.Reader) error {
	header := &packHeader{}
	if err := binary.Read(reader, binary.BigEndian, header); err != nil {
		return err
	}

	if header.Signature != packHeaderSignature {
		return Errorf("%v is not a valid pack header", header.Signature)
	}

	return nil
}

// Close closes the pack file associated with this Pack. If the pack file has
// never been opened in the first place, nothing happens.
func (p *Pack) Close() error {
	if p.file == nil {
		return nil
	}
	return p.file.Close()
}

// Objects returns a slice of sorted SHA-1 checksums of the objects in this
// pack.
func (p *Pack) Objects() []core.Sha1 {
	return p.idx.Objects()
}

// ObjectBySha1 returns the object corresponding to the given sha. If the object
// does not exist in this pack, nil is returned for the object and the error
// ErrObjectNotFoundInPack is returned. If there was an error in seeking in the
// pack file or decoding the pack entry within, nil is returned for the object
// along with the error that occurred.
func (p *Pack) ObjectBySha1(sha core.Sha1) (core.Object, error) {
	pos := p.idx.PosForSha1(sha)
	if pos == PackIndexPosNotFound {
		return nil, ErrObjectNotFoundInPack
	}

	if offset, err := p.idx.OffsetForPos(pos); err != nil {
		return nil, err
	} else if _, err := p.file.Seek(offset, 0); err != nil {
		return nil, err
	}

	entry := &packEntry{}
	err := entry.Decode(p.file)
	if err != nil {
		return nil, err
	}

	return entry, nil
}
