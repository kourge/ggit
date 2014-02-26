package config

import (
	"reflect"
	"strings"
	"testing"
)

type sectionFixture struct {
	Section
	String           string
	NormalizedString string
}

var (
	_fixtureSection1 sectionFixture = sectionFixture{
		Section: Section{"core", []Entry{
			{"repositoryformatversion", int64(0)},
			{"filemode", true},
			{"diff", "auto"},
			{"bare", false},
			{"name", "John Doe"},
		}},
		String: `[core]
	repositoryformatversion = 0
	# comment 1
	filemode = true

	diff = auto
; comment 2
	bare = false
	name = "John Doe"`,
		NormalizedString: `[core]
	repositoryformatversion = 0
	filemode = true
	diff = auto
	bare = false
	name = "John Doe"
`,
	}
)

func TestSection_String(t *testing.T) {
	var actual string = _fixtureSection1.Section.String()
	var expected string = _fixtureSection1.NormalizedString

	if actual != expected {
		t.Errorf("section.String() = %v, want %v", actual, expected)
	}
}

func TestSection_Decode(t *testing.T) {
	var actual *Section = &Section{}
	var expected *Section = &_fixtureSection1.Section

	actual.Decode(strings.NewReader(_fixtureSection1.String))

	if !reflect.DeepEqual(*actual, *expected) {
		t.Errorf("section.Decode() produced %v, want %v", *actual, *expected)
	}
}

func TestSection_Walk_Iterate(t *testing.T) {
	section := &_fixtureSection1.Section
	i := 0
	section.Walk(func(k string, v interface{}) (stop bool) {
		actualK, expectedK := k, section.Entries[i].Key
		if actualK != expectedK {
			t.Errorf("section.Walk() gave k = %v for index %d, want %v", actualK, i, expectedK)
		}

		actualV, expectedV := v, section.Entries[i].Value
		if actualV != expectedV {
			t.Errorf("section.Walk() gave v = %v for index %d, want %v", actualV, i, expectedV)
		}

		i += 1
		return false
	})
}

func TestSection_Walk_Break(t *testing.T) {
	section := &_fixtureSection1.Section
	i := 0
	section.Walk(func(k string, v interface{}) (stop bool) {
		i += 1
		return true
	})

	actual, expected := i, 1
	if actual != expected {
		t.Errorf("section.Walk() iterated %d time(s), want %d time(s)", actual, expected)
	}
}

func TestSection_Get(t *testing.T) {
	var actual interface{}
	var expected interface{} = "John Doe"

	section := &_fixtureSection1.Section
	actual, exists := section.Get("name")

	if !exists {
		t.Errorf(`section.Get("name")[1] = %v, want %v`, exists, true)
	}

	if actual != expected {
		t.Errorf(`section.Get("name")[0] = %v, want %v`, actual, expected)
	}
}

func TestSection_Get_Nonexist(t *testing.T) {
	var actual interface{}
	var expected interface{} = nil

	section := &_fixtureSection1.Section
	actual, exists := section.Get("foobar")

	if exists {
		t.Errorf(`section.Get("foobar")[1] = %v, want %v`, exists, false)
	}

	if actual != expected {
		t.Errorf(`section.Get("foobar")[0] = %v, want %v`, actual, expected)
	}
}

func TestSection_Set(t *testing.T) {
	section := _fixtureSection1.Section

	func(section *Section) {
		var actual bool = section.Set("foo", "bar")
		var expected bool = true

		if actual != expected {
			t.Errorf(`section.Set("foo", "bar") = %v, want %v`, actual, expected)
		}
	}(&section)

	func(section *Section) {
		actual, _ := section.Get("foo")
		expected := "bar"

		if actual != expected {
			t.Errorf(`section.Get("foo") = %v, want %v`, actual, expected)
		}
	}(&section)
}

func TestSection_Set_Exist(t *testing.T) {
	var section Section = _fixtureSection1.Section

	func(section *Section) {
		var actual bool = section.Set("bare", true)
		var expected bool = false

		if actual != expected {
			t.Errorf(`section.Set("bare", true) = %v, want %v`, actual, expected)
		}
	}(&section)

	func(section *Section) {
		actual, _ := section.Get("bare")
		expected := true

		if actual != expected {
			t.Errorf(`section.Get("bare") = %v, want %v`, actual, expected)
		}
	}(&section)
}

func TestSection_Del(t *testing.T) {
	var section Section = _fixtureSection1.Section

	func(section *Section) {
		actual, expected := section.Del("bare"), true

		if actual != expected {
			t.Errorf(`section.Del("bare") = %v, want %v`, actual, expected)
		}
	}(&section)

	func(section *Section) {
		actual1, actual2 := section.Get("bare")
		var expected1 interface{} = nil
		var expected2 bool = false

		if actual1 != expected1 {
			t.Errorf(`section.Get("bare")[0] = %v, want %v`, actual1, expected1)
		}
		if actual2 != expected2 {
			t.Errorf(`section.Get("bare")[1] = %v, want %v`, actual2, expected2)
		}
	}(&section)
}

func TestSection_Del_Nonexist(t *testing.T) {
	var section Section = _fixtureSection1.Section

	func(section *Section) {
		actual, expected := section.Del("foo"), false

		if actual != expected {
			t.Errorf(`section.Del("bare") = %v, want %v`, actual, expected)
		}
	}(&section)

	func(section *Section) {
		actual1, actual2 := section.Get("foo")
		var expected1 interface{} = nil
		var expected2 bool = false

		if actual1 != expected1 {
			t.Errorf(`section.Get("foo")[0] = %v, want %v`, actual1, expected1)
		}
		if actual2 != expected2 {
			t.Errorf(`section.Get("foo")[1] = %v, want %v`, actual2, expected2)
		}
	}(&section)
}
