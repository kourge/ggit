package config

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

type configFixture struct {
	Config
	String           string
	NormalizedString string
}

var (
	_fixtureConfig configFixture = configFixture{
		Config: Config{
			SectionMaker: NewEmptyOrderedSection,
			Sections: []Section{
				NewOrderedSection("user", []Entry{
					{"name", "Jane Doe"},
					{"email", "jane@example.com"},
				}...),
				NewOrderedSection("core", []Entry{
					{"repositoryformatversion", int64(0)},
					{"filemode", true},
					{"diff", "auto"},
					{"bare", false},
				}...),
			},
		},
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
		SectionMaker: NewEmptyOrderedSection,
		Sections: []Section{
			NewOrderedSection("core", []Entry{
				{"repositoryformatversion", int64(0)},
				{"filemode", true},
				{"bare", false},
				{"logallrefupdates", true},
				{"ignorecase", true},
				{"precomposeunicode", false},
			}...),
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
