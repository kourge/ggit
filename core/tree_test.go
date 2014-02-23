package core

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

var (
	_fixtureReadmeTreeEntry TreeEntry = TreeEntry{
		Mode: GitModeRegular | GitModeReadWritable,
		Type: "blob",
		Sha:  "3618cb8c4131839885ac273d74ee2eb8a7dd6970",
		Name: "README.md",
	}
	_fixtureReadmeTreeEntryString string    = "100644 blob 3618cb8c4131839885ac273d74ee2eb8a7dd6970\tREADME.md"
	_fixtureLicenseTreeEntry      TreeEntry = TreeEntry{
		Mode: GitModeRegular | GitModeReadWritable,
		Type: "blob",
		Sha:  "bf4b7bee80cf3f910fce252f73b189f1f3c2042a",
		Name: "LICENSE",
	}
	_fixtureLicenseTreeEntryString string = "100644 blob bf4b7bee80cf3f910fce252f73b189f1f3c2042a\tLICENSE"
	_fixtureTree                   Tree   = Tree{Entries: []TreeEntry{
		_fixtureLicenseTreeEntry,
		_fixtureReadmeTreeEntry,
	}}
)

func TestTree_Type(t *testing.T) {
	var actual string = _fixtureTree.Type()
	var expected string = "tree"

	if actual != expected {
		t.Errorf("tree.Type() = %v, want %v", actual, expected)
	}
}

func TestTree_Size(t *testing.T) {
	var actual int
	var expected = len([]byte(_fixtureLicenseTreeEntryString + "\n" + _fixtureReadmeTreeEntryString))

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(_fixtureTree.Reader())
	actual = buffer.Len()

	if actual != expected {
		t.Errorf("tree.Size() = %v, want %v", actual, expected)
	}
}

func TestTree_Reader(t *testing.T) {
	var actual []byte
	var expected = []byte(_fixtureLicenseTreeEntryString + "\n" + _fixtureReadmeTreeEntryString)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(_fixtureTree.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Error("tree.Reader() did not generate same byte sequence")
	}
}

func TestTree_Decode(t *testing.T) {
	var actual *Tree = &Tree{}
	var expected *Tree = &_fixtureTree

	r := io.MultiReader(
		strings.NewReader(_fixtureLicenseTreeEntryString),
		bytes.NewReader([]byte{'\n'}),
		strings.NewReader(_fixtureReadmeTreeEntryString),
	)
	err := actual.Decode(r)

	if err != nil {
		t.Errorf("tree.Decode() returned error %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("tree.Decode() produced %v, want %v", actual, expected)
	}
}
