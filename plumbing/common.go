package plumbing

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kourge/ggit/core"
)

const (
	EmptySha                    core.Sha1   = 0
	DefaultZlibCompressionLevel int         = 6
	DefaultObjectFileMode       os.FileMode = 0444
)

type fileFunc func(f *os.File) error

var fileNop fileFunc = func(f *os.File) error {
	return nil
}

type fileCreation struct {
	Name string
	Do   fileFunc
}

// Errorf is a wrapper around errors.New(fmt.Sprintf(format, rest...)).
func Errorf(format string, rest ...interface{}) error {
	return errors.New(fmt.Sprintf(format, rest...))
}

// IsRepo returns true if the directory at path is a valid Git repository.
func IsRepo(path string) bool {
	dir, err := os.Open(filepath.Clean(path))
	if err != nil {
		return false
	}

	if filenames, err := dir.Readdirnames(0); os.IsNotExist(err) {
		return false
	} else {
		for _, mustHave := range []string{"hooks", "info", "objects", "refs"} {
			found := false
		LOOKUP:
			for _, filename := range filenames {
				if filename == mustHave {
					found = true
					break LOOKUP
				}
			}

			if !found {
				return false
			}
		}
	}

	return true
}

