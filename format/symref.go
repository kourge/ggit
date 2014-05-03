package format

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/kourge/ggit/core"
)

var (
	symrefHeaderMagic   = []byte("ref: ")
	ErrRefInvalidHeader = errors.New("invalid symref header")
	ErrRefInvalidPath   = errors.New("invalid symref path")
)

// A Symref represents a symbolic ref. The most prominent example of a Git
// symbolic ref is HEAD, which points to the ref "refs/heads/{branch}".
// Historically, a symbolic ref was literally a symbolic link on the file system
// that pointed to a loose ref file, but this was deemed not portable enough, so
// now, a symbolic ref is a regular file whose content indicates the ref it
// points to.
type Symref struct {
	Path string
}

var _ core.EncodeDecoder = &Symref{}

// String returns the path of the ref to which this Symref points.
func (r Symref) String() string {
	return r.Path
}

func (r Symref) GoString() string {
	return fmt.Sprintf("Symref{%s}", strconv.Quote(r.Path))
}

// Reader returns an io.Reader that yields a version of this Symref that can be
// written to disk.
func (r Symref) Reader() io.Reader {
	return io.MultiReader(
		bytes.NewReader(symrefHeaderMagic),
		strings.NewReader(r.Path),
	)
}

// Decode reads from an io.Reader and attempts to parse the byte stream as a
// Symref that was previously written to disk.
func (r *Symref) Decode(reader io.Reader) error {
	header := make([]byte, len(symrefHeaderMagic))
	if n, err := reader.Read(header); err != nil || n != len(symrefHeaderMagic) {
		return ErrRefInvalidHeader
	} else if !bytes.Equal(header, symrefHeaderMagic) {
		return ErrRefInvalidHeader
	}

	if rest, err := ioutil.ReadAll(reader); err != nil {
		return err
	} else {
		r.Path = string(rest)
	}

	if !strings.HasPrefix(r.Path, "refs/") {
		return ErrRefInvalidPath
	}

	return nil
}
