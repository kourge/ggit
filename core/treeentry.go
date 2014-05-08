package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// A TreeEntry represents an item in a Git tree object. A TreeEntry itself is
// not a Git object, but points to one through its listed SHA-1 checksum.
type TreeEntry struct {
	Mode GitMode
	Name string
	Sha  Sha1
}

var _ EncodeDecoder = &TreeEntry{}

// Reader returns an io.Reader that yields a byte sequence in the format of
// "<mode> <name>\0<sha>", where mode is a 6-digit octal number representing
// a valid GitMode, name is a string that must not contain the NULL byte, and
// sha is a 20-byte-long raw representation of a SHA-1 checksum.
func (entry TreeEntry) Reader() io.Reader {
	return io.MultiReader(
		entry.Mode.Reader(),
		bytes.NewReader([]byte{' '}),
		strings.NewReader(entry.Name),
		bytes.NewReader([]byte{0}),
		entry.Sha.Reader(),
	)
}

// Decode parses a serialized tree item assumed to be in the format of "<mode>
// <name>\0<sha>". An error is returned if any of the following is true: the
// mode is not a 6-digit octal in ASCII form, the mode is not a valid Git mode,
// or the SHA-1 checksum is not exactly 20 bytes long.
func (entry *TreeEntry) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)

	if modeString, err := r.ReadString(byte(' ')); err != nil {
		return err
	} else {
		modeString = modeString[:len(modeString)-1]
		if mode, err := GitModeFromString(modeString); err != nil {
			return errors.New("invalid tree entry mode")
		} else {
			entry.Mode = mode
		}
	}

	if filename, err := r.ReadString(byte(0)); err != nil {
		return err
	} else {
		entry.Name = filename[:len(filename)-1]
	}

	var checksum Sha1
	if n, err := r.Read(checksum[:]); n != len(checksum) {
		return errors.New("failed to read fixed bytes for checksum")
	} else if err != nil {
		return err
	}
	entry.Sha = checksum

	return nil
}

// TreeEntryFromObject makes a new TreeEntry and attempts to populate it with
// a supported object. Supported object types are tree and blob. Since the only
// attribute that can be inferred from the given object is the SHA-1 checksum,
// you must supply the intended mode and name for the resulting TreeEntry.
func TreeEntryFromObject(object Object, mode GitMode, name string) (*TreeEntry, error) {
	entry := &TreeEntry{Mode: mode, Name: name}

	if t := object.Type(); t != "blob" && t != "tree" {
		message := fmt.Sprintf("%s is not a valid object type for a tree entry", t)
		return nil, errors.New(message)
	}

	entry.Sha = NewStream(object).Hash()
	return entry, nil
}

// A slice type that satisfies sort.Interface so that a slice of TreeEntries can
// be lexicographically sorted by their entry name.
type TreeEntrySlice []TreeEntry

func (entries TreeEntrySlice) Len() int {
	return len(entries)
}

func (entries TreeEntrySlice) Less(i, j int) bool {
	return entries[i].Name < entries[j].Name
}

func (entries TreeEntrySlice) Swap(i, j int) {
	temp := entries[i]
	entries[i] = entries[j]
	entries[j] = temp
}
