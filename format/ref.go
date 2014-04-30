package format

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"

	"github.com/kourge/ggit/core"
)

var (
	ErrInvalidRef = errors.New("ref not well-formed")
)

// A Git ref is simply a name in a repository that points to an object's SHA-1
// hash. The name usually starts with "refs/" and has multiple components
// separated by slashes.
//
// In a repository, a ref traditionally starts out as simply a file residing at
// a particular path whose content contains the SHA-1 hash of the object to
// which it points. This is a loose ref, but if it is not updated often, it ends
// up wasting space. Git may take multiple refs and store them all in the same
// file, producing packed refs.
type Ref struct {
	Name string
	Sha1 core.Sha1
}

var _ core.Decoder = &Ref{}

func (ref Ref) GoString() string {
	return fmt.Sprintf("Ref{%s}", strconv.Quote(ref.Name))
}

// Path returns the path that this ref would be located at if it were a loose
// ref at the given a root path.
func (ref Ref) Path(root string) string {
	return filepath.Join(root, ref.Name)
}

// Decode reads from a an io.Reader, presumably a loose ref file, and extracts
// the SHA-1 that is supposed to be in the file. If the file does not contain a
// properly formatted SHA-1 checksum, the error ErrInvalidRef is returned.
func (ref *Ref) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)

	if line, err := r.ReadString('\n'); err != nil {
		return err
	} else if len(line) != 41 {
		return ErrInvalidRef
	} else {
		if sha1, err := core.Sha1FromString(line[:]); err != nil {
			return ErrInvalidRef
		} else {
			ref.Sha1 = sha1
		}
	}

	return nil
}
