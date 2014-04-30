package plumbing

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kourge/ggit/core"
	"github.com/kourge/ggit/format"
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

// Path returns the path used to initialize the Repository.
func (repo *Repository) Path() string {
	return repo.path
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

// Packs returns a slice of all the packs in this repository.
func (repo *Repository) Packs() (packs []*format.Pack) {
	packPath := filepath.Join(repo.path, "objects", "pack")
	packDir, err := os.Open(packPath)
	if err != nil {
		core.Die(err)
	} else {
		defer packDir.Close()
	}

	if filenames, err := packDir.Readdirnames(0); err != nil {
		core.Die(err)
	} else {
		for _, filename := range filenames {
			if strings.HasSuffix(filename, ".pack") {
				packs = append(packs, format.NewPack(filepath.Join(packPath, filename)))
			}
		}
	}

	return packs
}
