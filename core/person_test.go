package core

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var (
	_fixtureTime         time.Time      = time.Unix(1234567890, 0).UTC()
	_fixtureTimezone     *time.Location = time.FixedZone("", -7*3600)
	_fixturePerson       Person         = Person{_fixtureAuthor, _fixtureTime}
	_fixturePersonString string         = "John Doe <john@example.com> 1234567890 +0000"
	_fixturePerson2      Person         = Person{
		Author{"Jane Doe", "jane@example.com"},
		time.Unix(1111111111, 0).In(_fixtureTimezone),
	}
	_fixturePerson2String string = "Jane Doe <jane@example.com> 1111111111 -0700"
)

func TestPerson_Reader(t *testing.T) {
	for _, s := range []struct {
		aString string
		person  Person
	}{
		{_fixturePersonString, _fixturePerson},
		{_fixturePerson2String, _fixturePerson2},
	} {
		var actual []byte
		var expected []byte = []byte(s.aString)

		buffer := new(bytes.Buffer)
		buffer.ReadFrom(s.person.Reader())
		actual = buffer.Bytes()

		if !bytes.Equal(actual, expected) {
			t.Error("person.Reader() did not generate same byte sequence")
		}
	}
}

func TestPerson_String(t *testing.T) {
	for _, s := range []struct {
		actual   string
		expected string
	}{
		{_fixturePerson.String(), _fixturePersonString},
		{_fixturePerson2.String(), _fixturePerson2String},
	} {
		if s.actual != s.expected {
			t.Errorf("person.String() = %s, want %s", s.actual, s.expected)
		}
	}
}

func TestPerson_Decode(t *testing.T) {
	for _, s := range []struct {
		expected *Person
		aString  string
	}{
		{&_fixturePerson, _fixturePersonString},
		{&_fixturePerson2, _fixturePerson2String},
	} {
		var actual *Person = &Person{}
		var expected *Person = s.expected
		err := actual.Decode(strings.NewReader(s.aString))

		if err != nil {
			t.Errorf("person.Decode() returned error %v", err)
		}

		if !actual.Equal(*expected) {
			t.Errorf("person.Decode() produced %#v, want %#v", actual, expected)
		}
	}
}

func TestPerson_Equal(t *testing.T) {
	if !_fixturePerson.Equal(_fixturePerson) {
		t.Errorf("person %#v does not equal itself", _fixturePerson)
	}

	if _fixturePerson.Equal(_fixturePerson2) {
		t.Errorf("%#v should not equal %#v", _fixturePerson)
	}
}

func TestNewPerson(t *testing.T) {
	var actual Person = NewPerson(
		"Jane Doe", "jane@example.com", 1111111111, -7*3600,
	)
	var expected Person = _fixturePerson2

	if !actual.Equal(expected) {
		t.Errorf("%#v should equal %#v", actual, expected)
	}
}
