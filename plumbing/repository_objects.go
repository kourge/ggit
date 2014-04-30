package plumbing

import (
	"compress/zlib"
	"errors"
	"os"
	"path/filepath"

	"github.com/kourge/ggit/core"
	"github.com/kourge/ggit/format"
)

var (
	ErrObjectNotFoundInRepo = errors.New("object not found in repository")
)

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
