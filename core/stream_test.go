package core

import (
	"bytes"
	"reflect"
	"testing"
)

type streamFixture struct {
	Object
	Body string
	Hash Sha1
}

var (
	_fixture1 = streamFixture{
		Object: &Blob{Content: []byte("what is up, doc?")},
		Body:   "blob 16\x00what is up, doc?",
		Hash:   Sha1("bd9dbf5aae1a3862dd1526723246b20206e5fc37"),
	}
	_fixture2 = streamFixture{
		Object: &Blob{Content: []byte("my hovercraft is full of eels")},
		Body:   "blob 29\x00my hovercraft is full of eels",
		Hash:   Sha1("7400f1589a11d1b912d6a90574d4f836087599b1"),
	}
	_fixture3 = streamFixture{
		Object: &Tree{Entries: []TreeEntry{
			{_frw_r__r__, "blob", _fixture1.Hash, "blob1"},
			{_frw_r__r__, "blob", _fixture2.Hash, "blob2"},
		}},
		Body: "tree 117\x00100644 blob bd9dbf5aae1a3862dd1526723246b20206e5fc37\tblob1\n100644 blob 7400f1589a11d1b912d6a90574d4f836087599b1\tblob2",
		Hash: Sha1("75242a2234bfad3ddc24e6c352a17cfbcef308b5"),
	}

	_fixture1Stream *Stream = &Stream{Object: _fixture1.Object}
	_fixture3Stream *Stream = &Stream{Object: _fixture3.Object}
)

func TestStream_Reader_Blob(t *testing.T) {
	var actual []byte
	var expected []byte = []byte(_fixture1.Body)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(_fixture1Stream.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Error("stream.Reader() did not generate same byte sequence")
	}
}

func TestStream_Reader_Tree(t *testing.T) {
	var actual []byte
	var expected []byte = []byte(_fixture3.Body)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(_fixture3Stream.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Error("stream.Reader() did not generate same byte sequence")
	}
}

func TestStream_Bytes_Blob(t *testing.T) {
	var actual []byte = _fixture1Stream.Bytes()
	var expected []byte = []byte(_fixture1.Body)

	if !bytes.Equal(actual, expected) {
		t.Error("stream.Bytes() did not generate same byte sequence")
	}
}

func TestStream_Bytes_Tree(t *testing.T) {
	var actual []byte = _fixture3Stream.Bytes()
	var expected []byte = []byte(_fixture3.Body)

	if !bytes.Equal(actual, expected) {
		t.Error("stream.Bytes() did not generate same byte sequence")
	}
}

func TestStream_Hash_Blob(t *testing.T) {
	var actual Sha1 = _fixture1Stream.Hash()
	var expected Sha1 = _fixture1.Hash

	if actual != expected {
		t.Errorf("stream.Hash() = %v, want %v", actual, expected)
	}
}

func TestStream_Hash_Tree(t *testing.T) {
	var actual Sha1 = _fixture3Stream.Hash()
	var expected Sha1 = _fixture3.Hash

	if actual != expected {
		t.Errorf("stream.Hash() = %v, want %v", actual, expected)
	}
}

func TestStream_Rehash(t *testing.T) {
	stream := &Stream{Object: _fixture1.Object}

	if actual, expected := stream.Hash(), _fixture1.Hash; actual != expected {
		t.Errorf("stream.Hash() = %v, want %v", actual, expected)
	}

	stream.Object = _fixture2.Object

	// Call stream.Hash() again. The hash should stay the same.
	if actual, expected := stream.Hash(), _fixture1.Hash; actual != expected {
		t.Errorf("stream.Hash() = %v, want %v", actual, expected)
	}

	stream.Rehash()

	// This time the hash should be updated.
	if actual, expected := stream.Hash(), _fixture2.Hash; actual != expected {
		t.Errorf("stream.Hash() = %v, want %v", actual, expected)
	}
}

func TestStream_Decode_Blob(t *testing.T) {
	var actual *Blob
	var expected *Blob = _fixture1.Object.(*Blob)
	stream := &Stream{}

	err := stream.Decode(bytes.NewBufferString(_fixture1.Body))
	if err != nil {
		t.Errorf("stream.Decode() returned error %v", err)
	}

	actual = stream.Object.(*Blob)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("stream.Object = %v, want %v", actual, expected)
	}
}

func TestStream_Decode_Tree(t *testing.T) {
	var actual []TreeEntry
	var expected []TreeEntry = _fixture3.Object.(*Tree).Entries
	stream := &Stream{}

	err := stream.Decode(bytes.NewBufferString(_fixture3.Body))
	if err != nil {
		t.Errorf("stream.Decode() returned error %v", err)
	}

	actual = stream.Object.(*Tree).Entries
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("stream.Object = %v, want %v", actual, expected)
	}
}
