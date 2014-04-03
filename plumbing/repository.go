package plumbing

import (
	"os"
	"path/filepath"

	"github.com/kourge/ggit/core"
)

// A Repository represents a potential Git repository.
type Repository struct {
	path string
}

// NewRepository returns a Repository given a path. The path is first cleaned
// on initialization.
func NewRepository(path string) *Repository {
	return &Repository{path: filepath.Clean(path)}
}

// IsRepo returns true if the path used to initialize the Repository is in fact
// a valid one.
func (repo *Repository) IsValid() bool {
	dir, err := os.Open(repo.path)
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

// UnpackedObjectBySha1 returns an Object with the given Sha1. If the object in
// question is packed or does not exist, an error is returned.
func (repo *Repository) UnpackedObjectBySha1(hash core.Sha1) (core.Object, error) {
	prefix, rest := hash.Split(2)
	path := filepath.Join(repo.path, prefix, rest)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	stream := &core.Stream{}
	if err := stream.Decode(file); err != nil {
		return nil, err
	}

	return stream.Object(), nil
}
