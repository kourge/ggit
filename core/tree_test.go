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
		Sha1: _sha("3618cb8c4131839885ac273d74ee2eb8a7dd6970"),
		Name: "README.md",
	}
	_fixtureReadmeTreeEntryString string    = "100644 README.md\x00\x36\x18\xcb\x8c\x41\x31\x83\x98\x85\xac\x27\x3d\x74\xee\x2e\xb8\xa7\xdd\x69\x70"
	_fixtureLicenseTreeEntry      TreeEntry = TreeEntry{
		Mode: GitModeRegular | GitModeReadWritable,
		Sha1: _sha("bf4b7bee80cf3f910fce252f73b189f1f3c2042a"),
		Name: "LICENSE",
	}
	_fixtureLicenseTreeEntryString string = "100644 LICENSE\x00\xbf\x4b\x7b\xee\x80\xcf\x3f\x91\x0f\xce\x25\x2f\x73\xb1\x89\xf1\xf3\xc2\x04\x2a"
	_fixtureTree                   *Tree  = NewTree([]TreeEntry{
		_fixtureLicenseTreeEntry,
		_fixtureReadmeTreeEntry,
	})
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
	var expected = len([]byte(_fixtureLicenseTreeEntryString + _fixtureReadmeTreeEntryString))

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(_fixtureTree.Reader())
	actual = buffer.Len()

	if actual != expected {
		t.Errorf("tree.Size() = %v, want %v", actual, expected)
	}
}

func TestTree_Reader(t *testing.T) {
	var actual []byte
	var expected = []byte(_fixtureLicenseTreeEntryString + _fixtureReadmeTreeEntryString)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(_fixtureTree.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Error("tree.Reader() did not generate same byte sequence")
	}
}

func TestTree_Decode(t *testing.T) {
	var actual *Tree = &Tree{}
	var expected *Tree = _fixtureTree

	r := io.MultiReader(
		strings.NewReader(_fixtureLicenseTreeEntryString),
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

func TestTree_Entries(t *testing.T) {
	var actual []TreeEntry = _fixtureTree.Entries()
	var expected []TreeEntry = []TreeEntry{
		_fixtureLicenseTreeEntry,
		_fixtureReadmeTreeEntry,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("tree.Entries() produced %v, want %v", actual, expected)
	}
}
