package plumbing

import (
	"compress/zlib"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/kourge/ggit/core"
)

// HashObjectOptions contains all the possible options for HashObject.
//
// Type is a string that is a Git object type, such as blob, tree, commit, or
// tag. If left unspecified as an empty string, it defaults to blob.
//
// Reader is an io.Reader that represents the byte stream that is to be treated
// as the hashed object's content. An error is returned if Reader is nil.
//
// Write is a bool that, when set to true, also causes the hashed object to be
// written to the specified Repo.
//
// Repo is a string that is a path to a repository. It is required when Write is
// true. An error is returned if Write is true and Repo is either unspecified or
// not a valid repository.
type HashObjectOptions struct {
	Type   string
	Reader io.Reader
	Write  bool
	Repo   string
}

// HashObject calculates the hash for a potential Git object. Equivalent to
// `git hash-object`. See the documentation on HashObjectOptions for more
// details.
func HashObject(o HashObjectOptions) (hash core.Sha1, err error) {
	if o.Reader == nil {
		return EmptySha, errors.New("Reader must not be nil")
	}

	if o.Type == "" {
		o.Type = "blob"
	}

	stream := &core.Stream{}
	switch o.Type {
	case "blob":
		stream.Object = &core.Blob{}
	case "tree":
		stream.Object = &core.Tree{}
	default:
		return EmptySha, Errorf("%v is not a valid Type", o.Type)
	}

	if err := stream.Object.Decode(o.Reader); err != nil {
		return EmptySha, err
	}
	hash = stream.Hash()

	if !o.Write {
		return
	}

	if o.Repo == "" {
		return hash, errors.New("must specify Repo")
	}
	if !IsRepo(o.Repo) {
		return hash, Errorf("not a repo: %s", o.Repo)
	}

	first, rest := hash.Split(2)
	slot := filepath.Join(o.Repo, "objects", first)
	if err := os.MkdirAll(slot, os.FileMode(0755)); err != nil {
		return hash, err
	}

	filepath := filepath.Join(slot, rest)
	if _, err := os.Stat(filepath); os.IsExist(err) {
		// If an object with this SHA-1 already exists, there is no need to
		// write it again.
		return hash, nil
	}

	file, err := os.Create(filepath)
	defer file.Close()

	writer, err := zlib.NewWriterLevel(file, DefaultZlibCompressionLevel)
	if err != nil {
		return hash, nil
	}
	defer writer.Close()

	io.Copy(writer, stream.Reader())
	file.Chmod(os.FileMode(0444))

	return
}
