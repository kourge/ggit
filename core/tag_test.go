package core

import (
	"bytes"
	"reflect"
	"testing"
)

var (
	_fixtureTagObject Sha1 = Sha1{
		0x6b, 0x6f, 0x8b, 0x56, 0x6e, 0xf3, 0x24, 0x5f, 0x5b, 0x25,
		0xd0, 0x3c, 0x61, 0xb2, 0xaf, 0x0a, 0x1f, 0x55, 0x30, 0x1e,
	}
	_fixtureTagObjectType string = "commit"
	_fixtureTagName       string = "v4.1.0.rc2"
	_fixtureTagTagger     Person = NewPerson(
		"David Heinemeier Hansson", "david@loudthinking.com", 1395778247, 1*3600,
	)
	_fixtureTagMessage string = "v4.1.0.rc2 release"
	_fixtureTag        *Tag   = NewTag(
		_fixtureTagObject,
		_fixtureTagObjectType,
		_fixtureTagName,
		_fixtureTagTagger,
		_fixtureTagMessage,
	)
	_fixtureTagString string = `object 6b6f8b566ef3245f5b25d03c61b2af0a1f55301e
type commit
tag v4.1.0.rc2
tagger David Heinemeier Hansson <david@loudthinking.com> 1395778247 +0100

v4.1.0.rc2 release
`
)

func TestTag_Type(t *testing.T) {
	var actual string = _fixtureTag.Type()
	var expected string = "tag"

	if actual != expected {
		t.Errorf("tag.Type() = %v, want %v", actual, expected)
	}
}

func TestTag_Size(t *testing.T) {
	var actual int = _fixtureTag.Size()
	var expected int = len(_fixtureTagString)

	if actual != expected {
		t.Errorf("tag.Size() = %v, want %v", actual, expected)
	}
}

func TestTag_Reader(t *testing.T) {
	var actual []byte
	var expected []byte = []byte(_fixtureTagString)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(_fixtureTag.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Error("tag.Reader() did not generate same byte sequence")
	}
}

func TestTag_Decode(t *testing.T) {
	var actual *Tag = &Tag{}
	var expected *Tag = _fixtureTag

	reader := bytes.NewBuffer([]byte(_fixtureTagString))
	err := actual.Decode(reader)

	if err != nil {
		t.Errorf("tag.Decode() returned error %v", err)
	}

	if !reflect.DeepEqual(*actual, *expected) {
		t.Errorf("tag.Decode() produced %v, want %v", actual, expected)
	}
}

func TestTag_Object(t *testing.T) {
	var actual Sha1 = _fixtureTag.Object()
	var expected Sha1 = _fixtureTagObject

	if actual != expected {
		t.Errorf("tag.Object() = %v, want %v", actual, expected)
	}
}

func TestTag_ObjectType(t *testing.T) {
	var actual string = _fixtureTag.ObjectType()
	var expected string = _fixtureTagObjectType

	if actual != expected {
		t.Errorf("tag.ObjectType() = %v, want %v", actual, expected)
	}
}

func TestTag_Name(t *testing.T) {
	var actual string = _fixtureTag.Name()
	var expected string = _fixtureTagName

	if actual != expected {
		t.Errorf("tag.Name() = %v, want %v", actual, expected)
	}
}

func TestTag_Tagger(t *testing.T) {
	var actual Person = _fixtureTag.Tagger()
	var expected Person = _fixtureTagTagger

	if !actual.Equal(expected) {
		t.Errorf("tag.Tagger() = %v, want %v", actual, expected)
	}
}

func TestTag_Message(t *testing.T) {
	var actual string = _fixtureTag.Message()
	var expected string = _fixtureTagMessage

	if actual != expected {
		t.Errorf("tag.Message() = %v, want %v", actual, expected)
	}
}
