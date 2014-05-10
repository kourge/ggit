package core

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"sort"
)

// A Tree is a Git object type that points to multiple Blobs and multiple Trees.
// It is conceptually a list of items, sorted ascending lexicographically by
// name. This sort order invariant is preserved even when this struct is
// decoded from a source that does not have its entries sorted.
type Tree struct {
	entries []TreeEntry
	buffer  []byte
}

var _ Object = &Tree{}

// NewTree returns a Tree containing the given TreeEntry structs. The entries
// are sorted to preserve the sort order invariant.
func NewTree(entries []TreeEntry) *Tree {
	tree := Tree{entries: entries}
	tree.loadIfNeeded()
	return &tree
}

func (tree *Tree) Type() string {
	return "tree"
}

func (tree *Tree) Size() int {
	tree.loadIfNeeded()
	return len(tree.buffer)
}

// Reader returns an io.Reader that yields each TreeEntry contained in this Tree
// in serialized format.
func (tree *Tree) Reader() io.Reader {
	tree.loadIfNeeded()
	return bytes.NewReader(tree.buffer)
}

// Decode reads from an io.Reader item by item and attempts to decode each as a
// TreeEntry. If any item is improperly formatted, an error is returned. After
// each item is properly decoded into a TreeEntry and stored, they are then
// sorted to satisfy the invariant.
func (tree *Tree) Decode(reader io.Reader) error {
	entries := make([]TreeEntry, 0)
	r := bufio.NewReader(reader)

	for {
		line, err := r.ReadBytes(byte(0))
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		var checksum Sha1
		n, err := r.Read(checksum[:])
		if n != len(checksum) {
			return errors.New("truncated checksum")
		} else if err != nil {
			return err
		}

		entryReader := io.MultiReader(
			bytes.NewBuffer(line),
			bytes.NewBuffer(checksum[:]),
		)
		treeEntry := &TreeEntry{}
		if parseErr := treeEntry.Decode(entryReader); parseErr != nil {
			return parseErr
		}

		entries = append(entries, *treeEntry)
	}

	tree.entries = entries
	tree.sort()
	tree.load()
	return nil
}

func (tree *Tree) loadIfNeeded() {
	if len(tree.buffer) == 0 {
		tree.load()
	}
}

func (tree *Tree) load() {
	tree.sort()

	buffer := new(bytes.Buffer)

	for _, entry := range tree.entries {
		if _, err := buffer.ReadFrom(entry.Reader()); err != nil {
			Die(err)
		}
	}
	tree.buffer = buffer.Bytes()
}

func (tree *Tree) sort() {
	sort.Sort(TreeEntrySlice(tree.entries))
}

// Entries returns a slice of the sorted tree entries in this Tree.
func (tree *Tree) Entries() []TreeEntry {
	entries := make([]TreeEntry, len(tree.entries))
	for i, entry := range tree.entries {
		entries[i] = entry
	}

	return entries
}
