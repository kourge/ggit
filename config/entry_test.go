package config

import (
	"strings"
	"testing"
)

type entryFixture struct {
	Entry
	String           string
	NormalizedString string
}

var (
	_fixtureEntry1 entryFixture = entryFixture{
		Entry{"diff", "auto"}, "diff =  auto", "diff = auto",
	}
	_fixtureEntry2 entryFixture = entryFixture{
		Entry{"name", "John Doe"}, "name = John Doe", "name = \"John Doe\"",
	}
	_fixtureEntry3 entryFixture = entryFixture{
		Entry{"ignorecase", true}, " ignorecase=true", "ignorecase = true",
	}
	_fixtureEntry4 entryFixture = entryFixture{
		Entry{"repositoryformatversion", int64(0)}, "repositoryformatversion  = 0", "repositoryformatversion = 0",
	}
)

func TestEntry_String(t *testing.T) {
	for _, fixture := range []entryFixture{
		_fixtureEntry1, _fixtureEntry2, _fixtureEntry3, _fixtureEntry4,
	} {
		var actual string = fixture.Entry.String()
		var expected string = fixture.NormalizedString

		if actual != expected {
			t.Errorf("entry.String() = %v, want %v", actual, expected)
		}
	}
}

func TestEntry_Decode(t *testing.T) {
	for _, fixture := range []entryFixture{
		_fixtureEntry1, _fixtureEntry2, _fixtureEntry3, _fixtureEntry4,
	} {
		var actual *Entry = &Entry{}
		var expected *Entry = &fixture.Entry

		actual.Decode(strings.NewReader(fixture.String))

		if *actual != *expected {
			t.Errorf("entry.Decode() produced %v, want %v", *actual, *expected)
		}
	}
}
