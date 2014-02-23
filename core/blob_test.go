package core

import (
	"bytes"
	"testing"
)

var (
	_fixtureContent []byte = []byte("what is up, doc?")
	_fixtureBlob    Blob   = Blob{Content: _fixtureContent}
)

func TestBlob_Type(t *testing.T) {
	var actual string = _fixtureBlob.Type()
	var expected string = "blob"

	if actual != expected {
		t.Errorf("blob.Type() = %v, want %v", actual, expected)
	}
}

func TestBlob_Size(t *testing.T) {
	var actual int = _fixtureBlob.Size()
	var expected int = len(_fixtureContent)

	if actual != expected {
		t.Errorf("blob.Size() = %v, want %v", actual, expected)
	}
}

func TestBlob_Reader(t *testing.T) {
	var actual []byte
	var expected []byte = _fixtureContent

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(_fixtureBlob.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Error("blob.Reader() did not generate same byte sequence")
	}
}

func TestBlob_Decode(t *testing.T) {
	blob := &Blob{}
	var actual []byte
	var expected []byte = _fixtureContent

	reader := bytes.NewBuffer(_fixtureContent)
	err := blob.Decode(reader)
	actual = blob.Content

	if err != nil {
		t.Errorf("blob.Decode() returned error %v", err)
	}

	if !bytes.Equal(actual, expected) {
		t.Error("blob.Decode() produced %v, want %v", actual, expected)
	}
}
