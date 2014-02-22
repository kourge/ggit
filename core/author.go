package core

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// An Author consists of a Name and an Email, both of which should be valid
// UTF-8 strings. Additionally, the Email must not contain the rune '<' or '>'.
type Author struct {
	Name  string
	Email string
}

var _ EncodeDecoder = &Author{}

// Reader returns an io.Reader that formats the author into a UTF-8 string
// in the human-readable form of "Name <Email>".
func (author Author) Reader() io.Reader {
	return io.MultiReader(
		strings.NewReader(author.Name),
		bytes.NewReader([]byte(" <")),
		strings.NewReader(author.Email),
		bytes.NewReader([]byte{'>'}),
	)
}

func (author Author) String() string {
	return fmt.Sprintf("%s <%s>", author.Name, author.Email)
}

// Decode parses a serialized Author assumed to be in the format of
// "Name <Email>". The whitespace between the name and the email may be
// arbitrarily long, but must not be absent.
func (author *Author) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)

	if name, err := r.ReadString(byte('<')); err != nil {
		return err
	} else {
		author.Name = strings.TrimSpace(name[0 : len(name)-1])
	}

	if email, err := r.ReadString(byte('>')); err != nil {
		return err
	} else {
		author.Email = email[0 : len(email)-1]
	}

	return nil
}
