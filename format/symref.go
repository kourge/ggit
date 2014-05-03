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
	symrefHeaderMagic      = []byte("ref: ")
	ErrInvalidSymref       = errors.New("invalid symref header")
	ErrInvalidSymrefTarget = errors.New("invalid symref target")
)

// A Symref represents a symbolic ref. The most prominent example of a Git
// symbolic ref is HEAD, which points to the ref "refs/heads/{branch}".
// Historically, a symbolic ref was literally a symbolic link on the file system
// that pointed to a loose ref file, but this was deemed not portable enough, so
// now, a symbolic ref is a regular file whose content indicates the ref it
// points to.
type Symref struct {
	Target string
}

var _ core.EncodeDecoder = &Symref{}

func (r Symref) GoString() string {
	return fmt.Sprintf("Symref{%s}", strconv.Quote(r.Target))
}

// Reader returns an io.Reader that yields a version of this Symref that can be
// written to disk.
func (r Symref) Reader() io.Reader {
	return io.MultiReader(
		bytes.NewReader(symrefHeaderMagic),
		strings.NewReader(r.Target),
	)
}

// Decode reads from an io.Reader and attempts to parse the byte stream as
// a Symref that was previously written to disk. If the data is invalid and
// cannot be interpreted as a symbolic ref, ErrInvalidSymref is returned. If the
// data is valid but points to something that does not look like a ref (i.e.
// does not start with "refs/"), ErrInvalidSymrefTarget is returned.
func (r *Symref) Decode(reader io.Reader) error {
	header := make([]byte, len(symrefHeaderMagic))
	if n, err := reader.Read(header); err != nil || n != len(symrefHeaderMagic) {
		return ErrInvalidSymref
	} else if !bytes.Equal(header, symrefHeaderMagic) {
		return ErrInvalidSymref
	}

	if rest, err := ioutil.ReadAll(reader); err != nil {
		return err
	} else {
		if newline := bytes.IndexByte(rest, '\n'); newline != -1 {
			rest = rest[:newline]
		}
		r.Target = string(rest)
	}

	if !strings.HasPrefix(r.Target, "refs/") {
		return ErrInvalidSymrefTarget
	}

	return nil
}
