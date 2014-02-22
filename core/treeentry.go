package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// A TreeEntry represents an item in a Git tree object. A TreeEntry itself is
// not a Git object, but points to one through its listed SHA-1 checksum.
type TreeEntry struct {
	Mode GitMode
	Type string
	Sha  Sha1
	Name string
}

var _ EncodeDecoder = &TreeEntry{}

// Reader returns an io.Reader that yields a byte sequence in the format of
// "<mode> <type> <sha>\t<name>", where mode is a 6-digit octal number
// representing a GitMode and the rest are strings.
func (entry TreeEntry) Reader() io.Reader {
	return io.MultiReader(
		entry.Mode.Reader(),
		bytes.NewReader([]byte{' '}),
		strings.NewReader(entry.Type),
		bytes.NewReader([]byte{' '}),
		strings.NewReader(string(entry.Sha)),
		bytes.NewReader([]byte{'\t'}),
		strings.NewReader(entry.Name),
	)
}

// Decode parses a serialized tree item assumed to be in the format of
// "<mode> <type> <sha>\t<name>". An error is returned if any of the following
// is true: the mode is not a 6-digit octal, the type is not blob or tree, or
// the SHA-1 checksum is not 40 characters long.
func (entry *TreeEntry) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)

	if modeString, err := r.ReadString(byte(' ')); err != nil {
		return err
	} else {
		modeString = modeString[:len(modeString)-1]
		if len(modeString) != 6 {
			return errors.New("tree entry mode is not 6 characters long")
		}
		if mode, err := GitModeFromString(modeString); err != nil {
			return errors.New("invalid tree entry mode")
		} else {
			entry.Mode = mode
		}
	}

	if typeString, err := r.ReadString(byte(' ')); err != nil {
		return err
	} else {
		typeString = typeString[:len(typeString)-1]
		if typeString != "blob" && typeString != "tree" {
			return errors.New(fmt.Sprintf("invalid tree item object type %#v", typeString))
		}
		entry.Type = typeString
	}

	if shaString, err := r.ReadString(byte('\t')); err != nil {
		return err
	} else {
		shaString = shaString[:len(shaString)-1]
		sha1 := Sha1(shaString)
		if !sha1.IsValid() {
			return errors.New(fmt.Sprintf("invalid tree item SHA-1 %v", shaString))
		}
		entry.Sha = sha1
	}

	if filenameBytes, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		entry.Name = string(filenameBytes)
	}

	return nil
}

// TreeEntryFromObject makes a new TreeEntry and attempts to populate it with a
// supported object. Supported object types are tree and blob. Since the only
// attributes that can be inferred from the given object are the type and the
// SHA-1 checksum, you must supply the intended mode and name for the resulting
// TreeEntry.
func TreeEntryFromObject(object Object, mode GitMode, name string) (*TreeEntry, error) {
	entry := &TreeEntry{Mode: mode, Name: name}

	if t := object.Type(); t != "blob" && t != "tree" {
		message := fmt.Sprintf("%s is not a valid object type for a tree entry", t)
		return nil, errors.New(message)
	} else {
		entry.Type = t
	}

	entry.Sha = (&Stream{Object: object}).Hash()
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
