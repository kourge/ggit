package format

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kourge/ggit/core"
)

// A GlobTable holds multiple glob patterns.
type GlobTable struct {
	Globs []string
}

var _ core.EncodeDecoder = &GlobTable{}

// GlobTableAtPath attempts to open the file at the given path and decode it
// into a GlobTable. If this succeeded, a GlobTable is returned with a nil
// error. If this failed, the error is returned with a nil GlobTable.
func GlobTableAtPath(path string) (*GlobTable, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	globs := &GlobTable{}
	if err := globs.Decode(file); err != nil {
		return nil, err
	}

	return globs, nil
}

// Match takes a name and returns whether the name matches any of the glob
// patterns in this GlobTable. If there is an invalid glob pattern in this
// GlobTable, then an error may be returned and the matched boolean is
// meaningless. Note that this is not guaranteed, since it is possible that
// the name matches an earlier valid pattern before an attempt to match the name
// against an invalid pattern even occurs.
func (t GlobTable) Match(name string) (matched bool, err error) {
	for _, pattern := range t.Globs {
		matched, err = filepath.Match(pattern, name)
		if err != nil {
			return
		} else if matched {
			return
		}
	}
	return false, nil
}

// WalkFunc wraps a filepath.WalkFunc so that only file names that are matched
// by this GlobTable will call the underlying WalkFunc.
func (t GlobTable) WalkFunc(f filepath.WalkFunc) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if match, err := t.Match(info.Name()); err != nil {
			return err
		} else if match {
			return f(path, info, err)
		}
		return nil
	}
}

// Reader returns an io.Reader that yields this GlobTable in a format that can
// be stored on disk in a human-readable format.
func (t GlobTable) Reader() io.Reader {
	globCount := len(t.Globs)
	readers := make([]io.Reader, globCount*2-1)

	for i, glob := range t.Globs {
		readers[i*2] = strings.NewReader(glob)
		if i != globCount-1 {
			readers[i*2+1] = bytes.NewReader([]byte{'\n'})
		}
	}

	return io.MultiReader(readers...)
}

// Decode takes an io.Reader and attempts to parse the resulting stream of bytes
// as a GlobTable, ignoring comments and whitespace.
func (t *GlobTable) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)

	atEof := false
	for !atEof {
		if line, err := r.ReadString('\n'); err == io.EOF {
			atEof = true
		} else if err != nil {
			return err
		} else {
			if pound := strings.IndexRune(line, '#'); pound != -1 {
				line = line[:pound]
			}
			line = strings.TrimSpace(line)

			if line == "" {
				continue
			}

			t.Globs = append(t.Globs, line)
		}
	}

	return nil
}
