package core

import (
	"io"
	"io/ioutil"
)

// A Commit is a Git object type that points to a Tree and one or more Commits
// as its parent(s). Furthermore, an author, a committer, and a message are
// associated with a Commit.
type Commit struct {
	tree      sha1field
	parents   []sha1field
	author    *Person
	committer *Person
	message   string

	buffer []byte
}

var _ Object = &Commit{}

// NewCommit returns a new Commit with the given attributes.
func NewCommit(
	tree Sha1,
	parents []Sha1,
	author, committer Person,
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

func (commit *Commit) Author() Person {
	return *commit.author
}

func (commit *Commit) Committer() Person {
	return *commit.committer
}

func (commit *Commit) Message() string {
	return commit.message
}

// ParentsEqual returns true if and only if this commit and the other commit
// have the same number of parents and the same parents in the same order.
// Otherwise false is returned.
func (commit *Commit) ParentsEqual(other *Commit) bool {
	a, b := commit.Parents(), other.Parents()

	if len(a) != len(b) {
		return false
	}

	for i, parent := range a {
		if b[i] != parent {
			return false
		}
	}
	return true
}

// Equal returns true if this commit and the other commit have the same tree,
// parents, author, committer, and message. Otherwise, false is returned.
func (commit *Commit) Equal(other *Commit) bool {
	return commit.Tree() == other.Tree() &&
		commit.ParentsEqual(other) &&
		commit.Author().Equal(other.Author()) &&
		commit.Committer().Equal(other.Committer()) &&
		commit.Message() == other.Message()
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
	n := len(commit.parents)
	fields := make(fieldslice, n+4)

	fields[0] = field{"tree", &commit.tree}
	for i, parent := range commit.parents {
		fields[i+1] = field{"parent", &parent}
	}
	fields[n+1] = field{"author", commit.author}
	fields[n+2] = field{"committer", commit.committer}
	fields[n+3] = field{"message", &StringCoder{commit.message}}

	return io.MultiReader(fields.Readers()...)
}

func (commit *Commit) load() {
	buffer, err := ioutil.ReadAll(commit.Reader())
	if err != nil {
		Die(err)
	}
	commit.buffer = buffer
}

func (commit *Commit) setBuffer(bytes []byte) {
	commit.buffer = bytes
}

// Decode reads from an io.Reader, attempting to decode the stream as a
// serialized Commit object. If any part of the stream is improperly formatted,
// an error is returned.
func (commit *Commit) Decode(reader io.Reader) error {
	return decodeFields(commit, reader)
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
			v := &Person{}
			err = v.Decode(s.Reader())
			commit.author = v
		case "committer":
			v := &Person{}
			err = v.Decode(s.Reader())
			commit.committer = v
		case "message":
			v := &StringCoder{}
			err = v.Decode(s.Reader())
			commit.message = v.string
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
