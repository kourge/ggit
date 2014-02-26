package config

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

type configSectionFixture struct {
	Section
	String           string
	NormalizedString string
}

type configFixture struct {
	Config
	String           string
	NormalizedString string
}

var (
	_fixtureSection1 configSectionFixture = configSectionFixture{
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

	_fixtureConfig configFixture = configFixture{
		Config: Config{[]Section{
			{"user", []Entry{
				{"name", "Jane Doe"},
				{"email", "jane@example.com"},
			}},
			{"core", []Entry{
				{"repositoryformatversion", int64(0)},
				{"filemode", true},
				{"diff", "auto"},
				{"bare", false},
			}},
		}},
		String: `
[user]
	name =   Jane Doe
	email =     jane@example.com

	[core]

repositoryformatversion = 0
	filemode = true
diff=auto
	bare=false
`,
		NormalizedString: `[user]
	name = "Jane Doe"
	email = jane@example.com

[core]
	repositoryformatversion = 0
	filemode = true
	diff = auto
	bare = false
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

func TestConfig_String(t *testing.T) {
	var actual string = _fixtureConfig.Config.String()
	var expected string = _fixtureConfig.NormalizedString

	if actual != expected {
		t.Errorf("config.String() = %v, want %v", actual, expected)
	}
}

func TestConfig_Decode(t *testing.T) {
	var actual *Config = &Config{}
	var expected *Config = &_fixtureConfig.Config

	actual.Decode(strings.NewReader(_fixtureConfig.String))

	if !reflect.DeepEqual(*actual, *expected) {
		t.Errorf("config.Decode() produced %v, want %v", *actual, *expected)
	}
}

func ExampleConfig_Reader() {
	defaultConfig := Config{
		[]Section{
			{"core", []Entry{
				{"repositoryformatversion", int64(0)},
				{"filemode", true},
				{"bare", false},
				{"logallrefupdates", true},
				{"ignorecase", true},
				{"precomposeunicode", false},
			}},
		},
	}

	io.Copy(os.Stdout, defaultConfig.Reader())
	// Output:
	// [core]
	// 	repositoryformatversion = 0
	// 	filemode = true
	// 	bare = false
	// 	logallrefupdates = true
	// 	ignorecase = true
	// 	precomposeunicode = false
}
