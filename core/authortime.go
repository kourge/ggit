package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// AuthorTime encapsulates both an Author and a time.Time, a construct usually
// used to indicate an author or committer.
type AuthorTime struct {
	Author
	time.Time
}

var _ EncodeDecoder = &AuthorTime{}

// NewAuthorTime returns an AuthorTime with the specified name and email as its
// Author and the given sec and offset as its Time. The provided sec is treated
// as a Unix time represending the number of seconds elapsed since the Unix
// epoch and the offset is the time zone numeric offset from UTC in the unit of
// seconds.
func NewAuthorTime(name, email string, sec int64, offset int) AuthorTime {
	t := time.Unix(sec, 0)
	tz := time.FixedZone("", offset)
	return AuthorTime{Author{name, email}, t.In(tz)}
}

// Reader returns an io.Reader that formats the AuthorTime into a UTF-8 string
// in the human-readable form of "Name <Email> UnixTime Timezone", where
// UnixTime is a Unix time (the number of seconds elapsed since January 1, 1970
// UTC) and Timezone is a numeric offset from UTC in the format of "±hhmm".
func (a AuthorTime) Reader() io.Reader {
	return io.MultiReader(
		a.Author.Reader(),
		bytes.NewReader([]byte{' '}),
		strings.NewReader(a.time()),
	)
}

func (a AuthorTime) time() string {
	return fmt.Sprintf("%d %s", a.Unix(), a.Format("-0700"))
}

// String returns this AuthorTime in the format of "Name <Email> UnixTime
// Timezone", where UnixTime is a decimal representation of a Unix time in the
// unit of seconds and Timezon is a numeric offset from UTC in the format of
// "±hhmm".
func (a AuthorTime) String() string {
	return fmt.Sprintf("%s <%s> %s", a.Name, a.Email, a.time())
}

// Decode parses a serialized AuthorTime assumed to be in the format of "Name
// <Email> UnixTime Timezone". The whitespace between the name and the email may
// be arbitrarily long, but must not be absent. The UnixTime must be an integer
// and is interpreted as the number of seconds elapsed since January 1, 1970
// UTF. The Timezone must be a numeric offset from UTC in the format of "±hhmm".
// If any of these components is malformed, an error is returned.
func (a *AuthorTime) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)

	authorString, err := r.ReadString(byte('>'))
	if err != nil {
		return err
	}

	author := &Author{}
	if err := author.Decode(strings.NewReader(authorString)); err != nil {
		return err
	}

	rest, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	parts := strings.Split(strings.TrimSpace(string(rest)), " ")
	if len(parts) != 2 {
		return errors.New("time component contains more than two fields")
	}

	secString, offsetString := parts[0], parts[1]
	var t time.Time

	sec, err := strconv.ParseInt(strings.TrimSpace(secString), 10, 64)
	if err != nil {
		return err
	}
	t = time.Unix(sec, 0)

	tz, err := time.Parse("-0700", offsetString)
	if err != nil {
		return err
	}
	t = t.In(tz.Location())

	a.Author = *author
	a.Time = t
	return nil
}

// Equal returns true if this AuthorTime and the other AuthorTime being compared
// share the same Author and have equal Time. When considering Time equality,
// the time zone is ignored.
func (a AuthorTime) Equal(b AuthorTime) bool {
	return a.Author == b.Author && a.Time.Equal(b.Time)
}
