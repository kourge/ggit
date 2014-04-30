package plumbing

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/kourge/ggit/core"
	"github.com/kourge/ggit/format"
)

var (
	ErrRefNotFound = errors.New("ref not found in repo")
	ErrInvalidRef  = format.ErrInvalidRef
)

// Sha1ByRef looks for a ref within this repository. It first searches the
// repository for the ref by calling Sha1FromLooseRef. If doing so yields an
// error that is not ErrRefNotFound, then the search is aborted and that error
// is returned. Otherwise, it falls back to calling Sha1FromPackedRefs.
func (repo *Repository) Sha1ByRef(ref string) (core.Sha1, error) {
	if sha1, err := repo.Sha1FromLooseRef(ref); err == nil {
		return sha1, err
	} else if err != ErrRefNotFound {
		return core.Sha1{}, err
	}

	return repo.Sha1FromPackedRefs(ref)
}

// Sha1FromLooseRef looks for a loose ref within this repository with the given
// name. If the ref file does not exist, the error ErrRefNotFound is returned.
// If the file exists but contains invalid data, then ErrInvalidRef is returned.
// If the file exists but failed to be opened, an appropriate error is returned.
// In any case, if an error occurred, the Sha1 returned is empty.
func (repo *Repository) Sha1FromLooseRef(ref string) (core.Sha1, error) {
	looseRef := &format.Ref{Name: ref}
	path := looseRef.Path(repo.Path())

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return core.Sha1{}, ErrRefNotFound
	} else if err != nil {
		return core.Sha1{}, err
	}
	defer file.Close()

	if err := looseRef.Decode(file); err != nil {
		return core.Sha1{}, err
	}

	return looseRef.Sha1, nil
}

// Sha1FromPackedRefs looks for a packed ref within this repository with the
// given name. If the packed refs file does not exist, the error ErrRefNotFound
// is returned. If the file exists but does not contain the ref with the given
// name, the error ErrRefNotFound is also returned. If the file exists but
// failed to be opened, an appropriate error is returned. If the file exists but
// contains invalid data, an appropriate error is returned.
func (repo *Repository) Sha1FromPackedRefs(ref string) (core.Sha1, error) {
	path := filepath.Join(repo.Path(), "packed-refs")

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return core.Sha1{}, ErrRefNotFound
	} else if err != nil {
		return core.Sha1{}, err
	}
	defer file.Close()

	packedRefs := &format.PackedRefs{}
	if err := packedRefs.Decode(file); err != nil {
		return core.Sha1{}, err
	}

	if sha1 := packedRefs.Sha1ForName(ref); sha1.IsEmpty() {
		return sha1, ErrRefNotFound
	} else {
		return sha1, nil
	}
}
