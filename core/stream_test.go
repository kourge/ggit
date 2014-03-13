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

func _sha(s string) Sha1 {
	sha, _ := Sha1FromString(s)
	return sha
}

var (
	_fixture1 = streamFixture{
		Object: &Blob{Content: []byte("what is up, doc?")},
		Body:   "blob 16\x00what is up, doc?",
		Hash:   _sha("bd9dbf5aae1a3862dd1526723246b20206e5fc37"),
	}
	_fixture2 = streamFixture{
		Object: &Blob{Content: []byte("my hovercraft is full of eels")},
		Body:   "blob 29\x00my hovercraft is full of eels",
		Hash:   _sha("7400f1589a11d1b912d6a90574d4f836087599b1"),
	}
	_fixture3 = streamFixture{
		Object: NewTree([]TreeEntry{
			{_frw_r__r__, "blob1", _fixture1.Hash},
			{_frw_r__r__, "blob2", _fixture2.Hash},
		}),
		Body: "tree 66\x00100644 blob1\x00\xbd\x9d\xbf\x5a\xae\x1a\x38\x62\xdd\x15\x26\x72\x32\x46\xb2\x02\x06\xe5\xfc\x37100644 blob2\x00\x74\x00\xf1\x58\x9a\x11\xd1\xb9\x12\xd6\xa9\x05\x74\xd4\xf8\x36\x08\x75\x99\xb1",
		Hash: _sha("dd08687e90cca5ce563867c40346781e3b115d36"),
	}

	_fixture1Stream *Stream = NewStream(_fixture1.Object)
	_fixture3Stream *Stream = NewStream(_fixture3.Object)
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

func TestStream_Decode_Blob(t *testing.T) {
	var actual *Blob
	var expected *Blob = _fixture1.Object.(*Blob)
	stream := &Stream{}

	err := stream.Decode(bytes.NewBufferString(_fixture1.Body))
	if err != nil {
		t.Errorf("stream.Decode() returned error %v", err)
	}

	actual = stream.Object().(*Blob)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("stream.Object = %v, want %v", actual, expected)
	}
}

func TestStream_Decode_Tree(t *testing.T) {
	var actual []TreeEntry
	var expected []TreeEntry = _fixture3.Object.(*Tree).entries
	stream := &Stream{}

	err := stream.Decode(bytes.NewBufferString(_fixture3.Body))
	if err != nil {
		t.Errorf("stream.Decode() returned error %v", err)
	}

	actual = stream.Object().(*Tree).entries
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("stream.Object = %v, want %v", actual, expected)
	}
}
