package core

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var (
	_fixtureTime             time.Time      = time.Unix(1234567890, 0).UTC()
	_fixtureTimezone         *time.Location = time.FixedZone("", -7*3600)
	_fixtureAuthorTime       AuthorTime     = AuthorTime{_fixtureAuthor, _fixtureTime}
	_fixtureAuthorTimeString string         = "John Doe <john@example.com> 1234567890 +0000"
	_fixtureAuthorTime2      AuthorTime     = AuthorTime{
		Author{"Jane Doe", "jane@example.com"},
		time.Unix(1111111111, 0).In(_fixtureTimezone),
	}
	_fixtureAuthorTime2String string = "Jane Doe <jane@example.com> 1111111111 -0700"
)

func TestAuthorTime_Reader(t *testing.T) {
	for _, s := range []struct {
		aString    string
		authortime AuthorTime
	}{
		{_fixtureAuthorTimeString, _fixtureAuthorTime},
		{_fixtureAuthorTime2String, _fixtureAuthorTime2},
	} {
		var actual []byte
		var expected []byte = []byte(s.aString)

		buffer := new(bytes.Buffer)
		buffer.ReadFrom(s.authortime.Reader())
		actual = buffer.Bytes()

		if !bytes.Equal(actual, expected) {
			t.Error("authortime.Reader() did not generate same byte sequence")
		}
	}
}

func TestAuthorTime_String(t *testing.T) {
	for _, s := range []struct {
		actual   string
		expected string
	}{
		{_fixtureAuthorTime.String(), _fixtureAuthorTimeString},
		{_fixtureAuthorTime2.String(), _fixtureAuthorTime2String},
	} {
		if s.actual != s.expected {
			t.Errorf("authortime.String() = %s, want %s", s.actual, s.expected)
		}
	}
}

func TestAuthorTime_Decode(t *testing.T) {
	for _, s := range []struct {
		expected *AuthorTime
		aString  string
	}{
		{&_fixtureAuthorTime, _fixtureAuthorTimeString},
		{&_fixtureAuthorTime2, _fixtureAuthorTime2String},
	} {
		var actual *AuthorTime = &AuthorTime{}
		var expected *AuthorTime = s.expected
		err := actual.Decode(strings.NewReader(s.aString))

		if err != nil {
			t.Errorf("authortime.Decode() returned error %v", err)
		}

		if !actual.Equal(*expected) {
			t.Errorf("authortime.Decode() produced %#v, want %#v", actual, expected)
		}
	}
}

func TestAuthorTime_Equal(t *testing.T) {
	if !_fixtureAuthorTime.Equal(_fixtureAuthorTime) {
		t.Errorf("authortime %#v does not equal itself", _fixtureAuthorTime)
	}

	if _fixtureAuthorTime.Equal(_fixtureAuthorTime2) {
		t.Errorf("%#v should not equal %#v", _fixtureAuthorTime)
	}
}

func TestNewAuthorTime(t *testing.T) {
	var actual AuthorTime = NewAuthorTime(
		"Jane Doe", "jane@example.com", 1111111111, -7*3600,
	)
	var expected AuthorTime = _fixtureAuthorTime2

	if !actual.Equal(expected) {
		t.Errorf("%#v should equal %#v", actual, expected)
	}
}
