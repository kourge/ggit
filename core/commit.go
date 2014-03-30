package core

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
)

// A Commit is a Git object type that points to a Tree and one or more Commits
// as its parent(s). Furthermore, an author, a committer, and a message are
// associated with a Commit.
type Commit struct {
	tree      sha1field
	parents   []sha1field
	author    *AuthorTime
	committer *AuthorTime
	message   string

	buffer []byte
}

var _ Object = &Commit{}

// NewCommit returns a new Commit with the given attributes.
func NewCommit(
	tree Sha1,
	parents []Sha1,
	author, committer AuthorTime,
	message string,
) *Commit {
	ps := make([]sha1field, len(parents))
	for i, parent := range parents {
		ps[i] = sha1field{parent}
	}
	commit := &Commit{sha1field{tree}, ps, &author, &committer, message, nil}
	commit.load()
	return commit
}

func (commit *Commit) Tree() Sha1 {
	return commit.tree.Sha1
}

func (commit *Commit) Parents() []Sha1 {
	parents := make([]Sha1, len(commit.parents))
	for i, parent := range commit.parents {
		parents[i] = parent.Sha1
	}
	return parents
}

func (commit *Commit) Author() AuthorTime {
	return *commit.author
}

func (commit *Commit) Committer() AuthorTime {
	return *commit.committer
}

func (commit *Commit) Message() string {
	return commit.message
}

func (commit *Commit) Type() string {
	return "commit"
}

func (commit *Commit) Size() int {
	if len(commit.buffer) == 0 {
		commit.load()
	}
	return len(commit.buffer)
}

// Reader returns an io.Reader that yields this commit in a serialized format.
func (commit *Commit) Reader() io.Reader {
	fields := fieldslice{
		{"tree", &commit.tree},
	}

	parents := make(fieldslice, len(commit.parents))
	for i, parent := range commit.parents {
		parents[i] = field{"parent", &parent}
	}

	fields = append(fields, parents...)
	fields = append(fields, fieldslice{
		{"author", commit.author},
		{"committer", commit.committer},
	}...)

	readers := append(fields.Readers(), []io.Reader{
		bytes.NewReader([]byte("\n\n")),
		strings.NewReader(commit.message),
		bytes.NewReader([]byte{'\n'}),
	}...)

	return io.MultiReader(readers...)
}

func (commit *Commit) load() {
	buffer, err := ioutil.ReadAll(commit.Reader())
	if err != nil {
		Die(err)
	}
	commit.buffer = buffer
}

// Decode reads from an io.Reader, attempting to decode the stream as a
// serialized Commit object. If any part of the stream is improperly formatted,
// an error is returned.
func (commit *Commit) Decode(reader io.Reader) error {
	lines, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	parts := bytes.SplitN(lines, []byte("\n\n"), 2)
	if len(parts) != 2 {
		return Errorf("received %d parts while decoding commit", len(parts))
	}

	fieldBytes, message := parts[0], parts[1]
	fields, err := fieldsliceDecode(bytes.NewReader(fieldBytes))
	if err != nil {
		return err
	}

	if err := commit.loadFields(fields); err != nil {
		return err
	}

	commit.message = strings.TrimRightFunc(string(message), unicode.IsSpace)
	commit.buffer = lines
	return nil
}

func (commit *Commit) loadFields(fields fieldslice) error {
	var parents []sha1field

	for _, field := range fields {
		var err error
		s := field.Value

		switch field.Name {
		case "tree":
			v := &sha1field{}
			err = v.Decode(s.Reader())
			commit.tree = *v
		case "parent":
			v := &sha1field{}
			err = v.Decode(s.Reader())
			parents = append(parents, *v)
		case "author":
			v := &AuthorTime{}
			err = v.Decode(s.Reader())
			commit.author = v
		case "committer":
			v := &AuthorTime{}
			err = v.Decode(s.Reader())
			commit.committer = v
		default:
			Die(Errorf("unrecognized commit field %s", field.Name))
		}

		if err != nil {
			return err
		}
	}

	commit.parents = parents
	return nil
}
