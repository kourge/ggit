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
	_fixture4Stream *Stream = NewStream(_fixture4.Object)

	_streamTests = []struct {
		Stream *Stream
		Fixture streamFixture
	}{
		{_fixture1Stream, _fixture1},
		{_fixture3Stream, _fixture3},
	}
)

func TestStream_Reader(t *testing.T) {
	for _, s := range _streamTests {
		var actual []byte
		var expected []byte = []byte(s.Fixture.Body)

		buffer := new(bytes.Buffer)
		buffer.ReadFrom(s.Stream.Reader())
		actual = buffer.Bytes()

		if !bytes.Equal(actual, expected) {
			t.Error("stream.Reader() did not generate same byte sequence")
		}
	}
}

func TestStream_Bytes(t *testing.T) {
	for _, s := range _streamTests {
		var actual []byte = s.Stream.Bytes()
		var expected []byte = []byte(s.Fixture.Body)

		if !bytes.Equal(actual, expected) {
			t.Error("stream.Bytes() did not generate same byte sequence")
		}
	}
}

func TestStream_Hash(t *testing.T) {
	for _, s := range _streamTests {
		var actual Sha1 = s.Stream.Hash()
		var expected Sha1 = s.Fixture.Hash

		if actual != expected {
			t.Errorf("stream.Hash() = %v, want %v", actual, expected)
		}
	}
}

func TestStream_Decode(t *testing.T) {
	for _, s := range _streamTests {
		var actual Object
		var expected Object = s.Fixture.Object
		stream := &Stream{}

		err := stream.Decode(bytes.NewBufferString(s.Fixture.Body))
		if err != nil {
			t.Errorf("stream.Decode() returned error %v", err)
		}

		actual = stream.Object()
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("stream.Object = %v, want %v", actual, expected)
		}
	}
}

