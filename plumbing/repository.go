package plumbing

import (
	"compress/zlib"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/kourge/ggit/core"
	"github.com/kourge/ggit/format"
)

var (
	ErrObjectNotFoundInRepo = errors.New("object not found in repository")
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

// PackedObjectBySha1 returns an Object with the given Sha1. If the object in
// question is loose or does not exist, the error ErrObjectNotFoundInRepo is
// returned.
func (repo *Repository) PackedObjectBySha1(hash core.Sha1) (core.Object, error) {
	for _, pack := range repo.Packs() {
		if err := pack.Open(); err != nil {
			return nil, err
		}
		defer pack.Close()

		if object, err := pack.ObjectBySha1(hash); err == format.ErrObjectNotFoundInPack {
			continue
		} else if err != nil {
			return nil, err
		} else {
			return object, nil
		}
	}

	return nil, ErrObjectNotFoundInRepo
}

// LooseObjectBySha1 returns an Object with the given Sha1. If the object in
// question is packed or does not exist, the error ErrObjectNotFoundInRepo is
// returned.
func (repo *Repository) LooseObjectBySha1(hash core.Sha1) (core.Object, error) {
	prefix, rest := hash.Split(2)
	path := filepath.Join(repo.path, "objects", prefix, rest)

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, ErrObjectNotFoundInRepo
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	r, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}

	stream := &core.Stream{}
	if err := stream.Decode(r); err != nil {
		return nil, err
	}

	return stream.Object(), nil
}

// ObjectBySha1 returns an Object with the given Sha1. First, the object is
// searched amongst loose objects. If it is not found there, then all pack files
// are searched in order. By now if it is stil not found, the error
// ErrObjectNotFoundInRepo is returned.
func (repo *Repository) ObjectBySha1(hash core.Sha1) (core.Object, error) {
	object, err := repo.LooseObjectBySha1(hash)
	if err == nil {
		return object, nil
	} else if err != ErrObjectNotFoundInRepo {
		return nil, err
	}

	return repo.PackedObjectBySha1(hash)
}
