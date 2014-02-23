package core

import (
	"bytes"
	"io"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var (
	_frw_r__r__ = GitModeRegular | GitModeReadWritable
	_d_________ = GitModeDir | GitModeNullPerm

	_fixtureTreeEntry TreeEntry = TreeEntry{
		_frw_r__r__, "blob", "3618cb8c4131839885ac273d74ee2eb8a7dd6970", "README.md",
	}
	_fixtureTreeEntryString string = "100644 blob 3618cb8c4131839885ac273d74ee2eb8a7dd6970\tREADME.md"

	_fixtureGitignoreBlob *Blob = &Blob{Content: []byte(`/*
!/.gitignore
!/Library/
!/CONTRIBUTING.md
!/README.md
!/SUPPORTERS.md
!/bin
/bin/*
!/bin/brew
!/share/man/man1/brew.1
.DS_Store
/Library/LinkedKegs
/Library/PinnedKegs
/Library/Taps
/Library/Formula/.gitignore
`)}
	_fixtureGitIgnoreTreeEntry TreeEntry = TreeEntry{
		_frw_r__r__, "blob", "458a3c1f135f68e9650d344525cd12a46d7048f5", ".gitignore",
	}
	_fixtureGitIgnoreHash   Sha1   = "458a3c1f135f68e9650d344525cd12a46d7048f5"
	_fixtureGitIgnoreString string = "100644 blob 458a3c1f135f68e9650d344525cd12a46d7048f5\t.gitignore"
)

func TestTreeEntry_Reader(t *testing.T) {
	treeEntry := _fixtureTreeEntry
	var actual []byte
	var expected []byte = []byte(_fixtureTreeEntryString)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(treeEntry.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Error("treeEntry.Reader() did not generate same byte sequence")
	}
}

func TestTreeEntry_Decode(t *testing.T) {
	var actual *TreeEntry = &TreeEntry{}
	var expected *TreeEntry = &_fixtureGitIgnoreTreeEntry
	err := actual.Decode(strings.NewReader(_fixtureGitIgnoreString))

	if err != nil {
		t.Errorf("treeEntry.Decode() returned error %v", err)
	}

	if *actual != *expected {
		t.Errorf("treeEntry.Decode() produced %v, want %v", actual, expected)
	}
}

func TestTreeEntryFromObject_Blob(t *testing.T) {
	blob := _fixtureGitignoreBlob
	name := ".gitignore"
	treeEntry, err := TreeEntryFromObject(blob, _frw_r__r__, name)

	if err != nil {
		t.Errorf("TreeEntryFromObject returned error %v", err)
	}

	if treeEntry.Mode != _frw_r__r__ {
		t.Errorf("treeEntry.Mode = %06o, want %06o", treeEntry.Mode, _frw_r__r__)
	}

	if treeEntry.Type != blob.Type() {
		t.Errorf("treeEntry.Type = %v, want %v", treeEntry.Type, blob.Type())
	}

	if treeEntry.Sha != _fixtureGitIgnoreHash {
		t.Errorf("treeEntry.Sha = %v, want %v", treeEntry.Sha, _fixtureGitIgnoreHash)
	}

	if treeEntry.Name != name {
		t.Errorf("treeEntry.Name = %v, want %v", treeEntry.Name, name)
	}
}

func TestTreeEntryFromObject_Tree(t *testing.T) {
}

type unknownObject struct{}

func (o *unknownObject) Type() string           { return "unknown" }
func (o *unknownObject) Size() int              { return 0 }
func (o *unknownObject) Reader() io.Reader      { return new(bytes.Buffer) }
func (o *unknownObject) Decode(io.Reader) error { return nil }

func TestTreeEntryFromObject_Unknown(t *testing.T) {
	var unknown Object = &unknownObject{}
	var _, err = TreeEntryFromObject(unknown, _frw_r__r__, "foobar")

	if err == nil {
		t.Errorf("TreeEntryFromObject did not return error when it should have", err)
	}
}

func TestTreeEntrySlice_Sort(t *testing.T) {
	_1 := TreeEntry{
		_frw_r__r__, "blob", "3618cb8c4131839885ac273d74ee2eb8a7dd6970", "README.md",
	}
	_2 := TreeEntry{
		_frw_r__r__, "blob", "00268614f04567605359c96e714e834db9cebab6", ".gitignore",
	}
	_3 := TreeEntry{
		_frw_r__r__, "blob", "bf4b7bee80cf3f910fce252f73b189f1f3c2042a", "LICENSE",
	}
	treeEntries := []TreeEntry{_1, _2, _3}

	var expected TreeEntrySlice = TreeEntrySlice([]TreeEntry{_2, _3, _1})
	var actual TreeEntrySlice = TreeEntrySlice(treeEntries)
	sort.Sort(actual)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Sorted tree entries = %v, want %v", actual, expected)
	}
}
