package core

import (
	"bytes"
	"strings"
	"testing"
)

var (
	_fixtureAuthor       Author = Author{Name: "John Doe", Email: "john@example.com"}
	_fixtureAuthorString string = "John Doe <john@example.com>"
)

func TestAuthorReader(t *testing.T) {
	var actual []byte
	var expected []byte = []byte(_fixtureAuthorString)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(_fixtureAuthor.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Error("author.Reader() did not generate same byte sequence")
	}
}

func TestAuthorString(t *testing.T) {
	var actual string = _fixtureAuthor.String()
	var expected string = _fixtureAuthorString

	if actual != expected {
		t.Errorf("author.String() = %v, want %v", actual, expected)
	}
}

func TestAuthorDecode(t *testing.T) {
	var actual *Author = &Author{}
	var expected *Author = &_fixtureAuthor
	err := actual.Decode(strings.NewReader(_fixtureAuthorString))

	if err != nil {
		t.Errorf("author.Decode() returned error %v", err)
	}

	if *actual != *expected {
		t.Errorf("author.Decode() produced %#v, want %#v", actual, expected)
	}
}
