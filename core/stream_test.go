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
	_fixture4 = streamFixture{
		Object: NewCommit(
			Sha1{
				0x93, 0x5e, 0x0a, 0x5c, 0x83, 0x61, 0xe5, 0x9f, 0x8b, 0xbc,
				0x01, 0xb2, 0xdb, 0xfb, 0xec, 0x3a, 0x44, 0xe2, 0x49, 0x04,
			},
			[]Sha1{{
				0x77, 0x5c, 0x72, 0x28, 0x62, 0x15, 0x59, 0x62, 0x34, 0x06,
				0x85, 0x7d, 0x18, 0x10, 0xa3, 0x15, 0x36, 0x16, 0x33, 0x6f,
			}},
			NewPerson(
				"Kosuke Asami", "tfortress58@gmail.com", 1395160458, 9*3600,
			),
			NewPerson(
				"Jack Nagel", "jacknagel@gmail.com", 1395293290, -5*3600,
			),
			`byobu 5.75

This release includes fixes about prefix problem that is discussed
in #27045.

Closes #27667.

Signed-off-by: Jack Nagel <jacknagel@gmail.com>`,
		),
		Body: "commit 371\x00tree 935e0a5c8361e59f8bbc01b2dbfbec3a44e24904\nparent 775c7228621559623406857d1810a3153616336f\nauthor Kosuke Asami <tfortress58@gmail.com> 1395160458 +0900\ncommitter Jack Nagel <jacknagel@gmail.com> 1395293290 -0500\n\nbyobu 5.75\n\nThis release includes fixes about prefix problem that is discussed\nin #27045.\n\nCloses #27667.\n\nSigned-off-by: Jack Nagel <jacknagel@gmail.com>\n",
		Hash: _sha("91465a197c01a5f022a224a592e769147db145a2"),
	}
	_fixture5 = streamFixture{
		Object: NewTag(
			Sha1{
				0x6b, 0x6f, 0x8b, 0x56, 0x6e, 0xf3, 0x24, 0x5f, 0x5b, 0x25,
				0xd0, 0x3c, 0x61, 0xb2, 0xaf, 0x0a, 0x1f, 0x55, 0x30, 0x1e,
			},
			"commit",
			"v4.1.0.rc2",
			NewPerson(
				"David Heinemeier Hansson", "david@loudthinking.com", 1395778247, 1*3600,
			),
			"v4.1.0.rc2 release",
		),
		Body: "tag 169\x00object 6b6f8b566ef3245f5b25d03c61b2af0a1f55301e\ntype commit\ntag v4.1.0.rc2\ntagger David Heinemeier Hansson <david@loudthinking.com> 1395778247 +0100\n\nv4.1.0.rc2 release\n",
		Hash: _sha("d82b255a0f16a06ebd2a3fbfe4893719d697c043"),
	}

	_fixture1Stream *Stream = NewStream(_fixture1.Object)
	_fixture3Stream *Stream = NewStream(_fixture3.Object)
	_fixture4Stream *Stream = NewStream(_fixture4.Object)
	_fixture5Stream *Stream = NewStream(_fixture5.Object)

	_streamTests = []struct {
		Stream  *Stream
		Fixture streamFixture
	}{
		{_fixture1Stream, _fixture1},
		{_fixture3Stream, _fixture3},
		{_fixture4Stream, _fixture4},
		{_fixture5Stream, _fixture5},
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
