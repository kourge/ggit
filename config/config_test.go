package config

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

type configEntryFixture struct {
	ConfigEntry
	String           string
	NormalizedString string
}

type configSectionFixture struct {
	ConfigSection
	String           string
	NormalizedString string
}

type configFixture struct {
	Config
	String           string
	NormalizedString string
}

var (
	_fixtureEntry1 configEntryFixture = configEntryFixture{
		ConfigEntry{"diff", "auto"}, "diff =  auto", "diff = auto",
	}
	_fixtureEntry2 configEntryFixture = configEntryFixture{
		ConfigEntry{"name", "John Doe"}, "name = John Doe", "name = \"John Doe\"",
	}
	_fixtureEntry3 configEntryFixture = configEntryFixture{
		ConfigEntry{"ignorecase", true}, " ignorecase=true", "ignorecase = true",
	}
	_fixtureEntry4 configEntryFixture = configEntryFixture{
		ConfigEntry{"repositoryformatversion", int64(0)}, "repositoryformatversion  = 0", "repositoryformatversion = 0",
	}

	_fixtureSection1 configSectionFixture = configSectionFixture{
		ConfigSection: ConfigSection{"core", []ConfigEntry{
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
		Config: Config{[]ConfigSection{
			{"user", []ConfigEntry{
				{"name", "Jane Doe"},
				{"email", "jane@example.com"},
			}},
			{"core", []ConfigEntry{
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

func TestConfigEntry_String(t *testing.T) {
	for _, fixture := range []configEntryFixture{
		_fixtureEntry1, _fixtureEntry2, _fixtureEntry3, _fixtureEntry4,
	} {
		var actual string = fixture.ConfigEntry.String()
		var expected string = fixture.NormalizedString

		if actual != expected {
			t.Errorf("entry.String() = %v, want %v", actual, expected)
		}
	}
}

func TestConfigEntry_Decode(t *testing.T) {
	for _, fixture := range []configEntryFixture{
		_fixtureEntry1, _fixtureEntry2, _fixtureEntry3, _fixtureEntry4,
	} {
		var actual *ConfigEntry = &ConfigEntry{}
		var expected *ConfigEntry = &fixture.ConfigEntry

		actual.Decode(strings.NewReader(fixture.String))

		if *actual != *expected {
			t.Errorf("entry.Decode() produced %v, want %v", *actual, *expected)
		}
	}
}

func TestConfigSection_String(t *testing.T) {
	var actual string = _fixtureSection1.ConfigSection.String()
	var expected string = _fixtureSection1.NormalizedString

	if actual != expected {
		t.Errorf("section.String() = %v, want %v", actual, expected)
	}
}

func TestConfigSection_Decode(t *testing.T) {
	var actual *ConfigSection = &ConfigSection{}
	var expected *ConfigSection = &_fixtureSection1.ConfigSection

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
		[]ConfigSection{
			{"core", []ConfigEntry{
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
