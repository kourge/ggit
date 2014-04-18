package core

import (
	"io"
	"io/ioutil"
)

// A Tag is a Git object type that points to a Commit. There are two kinds of
// tags in Git: a lightweight tag that is simply a ref with a name that points
// to the SHA-1 checksum of a commit, and an annotated tag that is a Git object
// containing a tagger, a timestamp, and a tag message. This type represents
// the latter kind.
type Tag struct {
	object     sha1field
	objectType string
	name       string
	tagger     *Person
	message    string

	buffer []byte
}

var _ Object = &Tag{}

// NewTag returns a new Tag with the given attributes. The objectType is almost
// always "commit".
func NewTag(
	object Sha1,
	objectType string,
	name string,
	tagger Person,
	message string,
) *Tag {
	return &Tag{sha1field{object}, objectType, name, &tagger, message, nil}
}

func (tag *Tag) Object() Sha1 {
	return tag.object.Sha1
}

func (tag *Tag) ObjectType() string {
	return tag.objectType
}

func (tag *Tag) Name() string {
	return tag.name
}

func (tag *Tag) Tagger() Person {
	return *tag.tagger
}

func (tag *Tag) Message() string {
	return tag.message
}

func (tag *Tag) Type() string {
	return "tag"
}

func (tag *Tag) Size() int {
	if len(tag.buffer) == 0 {
		tag.load()
	}
	return len(tag.buffer)
}

// Reader returns an io.Reader that yields this annotated tag in a serialized
// format.
func (tag *Tag) Reader() io.Reader {
	return io.MultiReader(fieldslice{
		{"object", &tag.object},
		{"type", &StringCoder{tag.objectType}},
		{"tag", &StringCoder{tag.name}},
		{"tagger", tag.tagger},
		{"message", &StringCoder{tag.message}},
	}.Readers()...)
}

func (tag *Tag) load() {
	buffer, err := ioutil.ReadAll(tag.Reader())
	if err != nil {
		Die(err)
	}
	tag.buffer = buffer
}

func (tag *Tag) setBuffer(bytes []byte) {
	tag.buffer = bytes
}

// Decode reads from an io.Reader, attempting to decode the stream as
// a serialized Tag object. If any part of the stream is improperly formatted,
// an error is returned.
func (tag *Tag) Decode(reader io.Reader) error {
	return decodeFields(tag, reader)
}

func (tag *Tag) loadFields(fields fieldslice) error {
	for _, field := range fields {
		var err error
		s := field.Value

		switch field.Name {
		case "object":
			v := &sha1field{}
			err = v.Decode(s.Reader())
			tag.object = *v
		case "type":
			v := &StringCoder{}
			err = v.Decode(s.Reader())
			tag.objectType = v.string
		case "tag":
			v := &StringCoder{}
			err = v.Decode(s.Reader())
			tag.name = v.string
		case "tagger":
			v := &Person{}
			err = v.Decode(s.Reader())
			tag.tagger = v
		case "message":
			v := &StringCoder{}
			err = v.Decode(s.Reader())
			tag.message = v.string
		default:
			Die(Errorf("unrecognized tag field %s", field.Name))
		}

		if err != nil {
			return err
		}
	}

	return nil
}
