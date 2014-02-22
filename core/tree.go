package core

import (
	"bufio"
	"bytes"
	"io"
	"sort"
)

// A Tree is a Git object type that points to multiple Blobs and multiple Trees.
// It is conceptually a list of items, sorted ascending lexicographically by
// name. This sort order invariant is preserved even when this struct is
// initialized with an unsorted slice of TreeEntry structs.
//
// If you replace tree.Entries manually after initializing the Tree struct or
// after calling Decode, you must call Reload to ensure that this invariant is
// maintained and that the cache is regenerated.
type Tree struct {
	Entries []TreeEntry
	buffer  []byte
}

var _ Object = &Tree{}

func (tree *Tree) Type() string {
	return "tree"
}

func (tree *Tree) Size() int {
	tree.loadIfNeeded()
	return len(tree.buffer)
}

// Reader returns an io.Reader that yields each TreeEntry contained in this Tree
// in serialized format, each separated by the new line rune '\n'. Before this
// is done, the entries are sorted if needed. This bulk serialization is only
// done once and then internally cached.
func (tree *Tree) Reader() io.Reader {
	tree.loadIfNeeded()
	return bytes.NewReader(tree.buffer)
}

// Decode reads from an io.Reader line by line and decodes each line as a
// TreeEntry. If any line is improperly formatted, an error is returned. After
// each line is properly decoded into a TreeEntry and stored, they are then
// sorted to satisfy the invariant.
func (tree *Tree) Decode(reader io.Reader) error {
	var err error
	buffer := new(bytes.Buffer)
	entries := make([]TreeEntry, 0)
	r := bufio.NewReader(reader)

	reachedEof := false
	for !reachedEof {
		line, err := r.ReadBytes(byte('\n'))
		if err == io.EOF {
			reachedEof = true
		} else if err != nil {
			break
		}

		if _, err := buffer.Write(line); err != nil {
			return err
		}

		i := len(line)
		if !reachedEof {
			i -= 1
		}
		lineReader := bytes.NewBuffer(line[:i])

		treeEntry := &TreeEntry{}
		if parseErr := treeEntry.Decode(lineReader); parseErr != nil {
			return parseErr
		}

		entries = append(entries, *treeEntry)
	}

	if err != nil && err != io.EOF {
		return err
	}

	tree.Entries = entries
	tree.buffer = buffer.Bytes()
	tree.sort()
	return nil
}

func (tree *Tree) loadIfNeeded() {
	if len(tree.buffer) == 0 {
		tree.Reload()
	}
}

// Reload ensures that the internal state is synchronized correctly with
// tree.Entries. Normally you do not need to call this method; however, if you
// replace tree.Entries after tree's initiallization or after calling Decode,
// you should call this method.
//
// Two things are done on reload: first, tree.Entries is sorted to maintain
// the sort order invariant; second, the internal serialization cache is
// immediately regenerated. For more information on the internal cache, see the
// description of the Reader method.
func (tree *Tree) Reload() {
	tree.sort()

	buffer := new(bytes.Buffer)

	for i, entry := range tree.Entries {
		if i != 0 {
			if err := buffer.WriteByte('\n'); err != nil {
				die(err)
			}
		}

		if _, err := buffer.ReadFrom(entry.Reader()); err != nil {
			die(err)
		}
	}
	tree.buffer = buffer.Bytes()
}

func (tree *Tree) sort() {
	sort.Sort(TreeEntrySlice(tree.Entries))
}
