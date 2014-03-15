package plumbing

import (
	"os"
	"path/filepath"

	"github.com/kourge/ggit/core"
)

const (
	DefaultZlibCompressionLevel int         = 6
	DefaultObjectFileMode       os.FileMode = 0444
)

var (
	EmptySha core.Sha1 = core.Sha1([20]byte{})
)

// Errorf is a wrapper around errors.New(fmt.Sprintf(format, rest...)).
var Errorf = core.Errorf

// IsRepo returns true if the directory at path is a valid Git repository.
func IsRepo(path string) bool {
	dir, err := os.Open(filepath.Clean(path))
	if err != nil {
		return false
	}
	defer dir.Close()

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
