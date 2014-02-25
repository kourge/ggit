package core

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// A Config represents a INI-style Git config file, composed of multiple
// sections, each of them with their own list of key-value pairs. Includes are
// not supported at the moment.
type Config struct {
	Sections []ConfigSection
}

var _ EncodeDecoder = &Config{}

// Reader returns an io.Reader that, when read, prints out each section,
// each separated by a blank line.
func (config Config) Reader() io.Reader {
	readers := make([]io.Reader, len(config.Sections)*2)

	for i, section := range config.Sections {
		readers[i*2] = section.Reader()
		readers[i*2+1] = bytes.NewReader([]byte{'\n'})
	}

	return io.MultiReader(readers[:len(readers)-1]...)
}

// String returns a string that is the result of draining the io.Reader
// returned by Reader().
func (config Config) String() string {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(config.Reader())
	return buffer.String()
}

// Decode parses bytes into a Config. This stream of bytes is assumed to contain
// multiple sections.
//
// Lines that start with the pound sign rune '#' or the semicolon rune ';' will
// be treated as comment and ignored. Completely whitespace lines or blank lines
// will also be ignored.
func (config *Config) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)
	lines := new(bytes.Buffer)

	flush := func(realloc bool) error {
		section := &ConfigSection{}
		if err := section.Decode(lines); err != nil {
			return err
		}
		config.Sections = append(config.Sections, *section)

		if realloc {
			lines = new(bytes.Buffer)
		}
		return nil
	}

	reachedEof := false
	for !reachedEof {
		line, err := r.ReadString(byte('\n'))
		if err != nil {
			if err != io.EOF {
				return err
			}
			reachedEof = true
		}
		line = strings.TrimSpace(line)

		if len(line) == 0 || line[0] == '#' || line[0] == ';' {
			continue
		}

		// A new section is starting. Unless it's already empty, flush the
		// current line buffer and decode the lines into a new section.
		if line[0] == '[' && lines.Len() != 0 {
			if err := flush(true); err != nil {
				return err
			}
		}

		if _, err := lines.WriteString(line); err != nil {
			return err
		}
		if err := lines.WriteByte('\n'); err != nil {
			return err
		}
	}

	if err := flush(false); err != nil {
		return err
	}
	return nil
}

// A ConfigSection represents a section in a Git config file. It has its own
// name and a list of key-value pairs under it.
type ConfigSection struct {
	Name    string
	Entries []ConfigEntry
}

var _ Encoder = ConfigSection{}

func (section ConfigSection) bytesBuffer() *bytes.Buffer {
	buffer := new(bytes.Buffer)

	buffer.WriteRune('[')
	buffer.WriteString(section.Name)
	buffer.WriteRune(']')
	buffer.WriteRune('\n')

	for _, entry := range section.Entries {
		buffer.WriteRune('\t')
		buffer.WriteString(entry.String())
		buffer.WriteRune('\n')
	}

	return buffer
}

// Returns an io.Reader that yields the section name in square brackets as the
// first line. Subsequent lines are the underlying ConfigEntry structs
// serialized in order, each indented by a single horizontal tab rune '\t' and
// each separated by a new line rune '\n'.
func (section ConfigSection) Reader() io.Reader {
	return section.bytesBuffer()
}

func (section ConfigSection) lazyReader() io.Reader {
	offset := 3
	readers := make([]io.Reader, len(section.Entries)*3+offset)
	readers[0] = bytes.NewReader([]byte{'['})
	readers[1] = strings.NewReader(section.Name)
	readers[2] = bytes.NewReader([]byte{']', '\n'})

	for i, entry := range section.Entries {
		readers[i*3+offset+0] = bytes.NewReader([]byte{'\t'})
		readers[i*3+offset+1] = entry.Reader()
		readers[i*3+offset+2] = bytes.NewReader([]byte{'\n'})
	}

	return io.MultiReader(readers...)
}

// String returns a string that is the result of draining the io.Reader
// returned by Reader().
func (section ConfigSection) String() string {
	return section.bytesBuffer().String()
}

// Decode parses bytes into a ConfigSection. This stream of bytes is assumed
// to contain only a single section. If multiple sections appear, all entries
// will be treated as if they were in a single section, and the last seen
// section's name is considered to be this single section's name.
//
// Lines that start with the pound sign rune '#' or the semicolon rune ';' will
// be treated as comment and ignored. Completely whitespace lines or blank lines
// will also be ignored.
func (section *ConfigSection) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)
	entries := make([]ConfigEntry, 0)

	reachedEof := false
	for !reachedEof {
		line, err := r.ReadString(byte('\n'))
		if err != nil {
			if err != io.EOF {
				return err
			}
			reachedEof = true
		}
		line = strings.TrimSpace(line)

		if len(line) == 0 || line[0] == '#' || line[0] == ';' {
			continue
		}

		if line[0] == '[' {
			i := strings.IndexByte(line, ']')
			section.Name = strings.TrimSpace(line[1:i])
			continue
		}

		entry := &ConfigEntry{}
		entry.Decode(strings.NewReader(line))
		entries = append(entries, *entry)
	}

	section.Entries = entries
	return nil
}

// A ConfigEntry represents a key-value pair in a Git config file. The key is an
// alphanumeric string, and the value is a decimal integer, a boolean value, or
// a string.
type ConfigEntry struct {
	Key   string
	Value interface{}
}

var _ EncodeDecoder = &ConfigEntry{}

// Reader returns an io.Reader that wraps around the string returned by String().
func (entry ConfigEntry) Reader() io.Reader {
	return strings.NewReader(entry.String())
}

// String returns the ConfigEntry converted to the string form of
// "<Key> = <Value>". If Value is an int or a boolean, their literal forms are
// used. If Value is a string, it is quoted with double quotes unless it
// contains no whitespace, as determined by unicode.IsSpace.
func (entry ConfigEntry) String() string {
	if value, ok := entry.Value.(string); ok {
		if strings.IndexFunc(value, unicode.IsSpace) != -1 {
			return fmt.Sprintf("%s = %s", entry.Key, strconv.Quote(value))
		}
	}
	return fmt.Sprintf("%s = %v", entry.Key, entry.Value)
}

// Decode parses a single line into a ConfigEntry, assuming the form
// "<Key> = <Value>". The Key should be alphanumeric, but must not contain the
// rune '='. The Value may have extra whitespace before, after, or both, as they
// are stripped in the decoding process.
//
// If the value starts with the double-quote rune '"', then it is assumed to
// begin and end with that rune; it is then treated as a quoted string with
// possible whitespace escape sequences that mirror the format of a Go string
// literal.
//
// If the value starts with a rune that is considered a number, as determined by
// unicode.IsDigit, then it is treated like a base-10 integer and an attempt is
// made to parse it as such. The parsed integer will always be of type int64.
//
// If the value is the literal "true" or "false", it is considered a boolean.
//
// In all other cases, the Value is treated as a string.
func (entry *ConfigEntry) Decode(reader io.Reader) error {
	var key string
	var value interface{}
	r := bufio.NewReader(reader)

	if line, err := r.ReadString(byte('=')); err != nil {
		return err
	} else {
		key = strings.TrimSpace(line[:len(line)-1])
	}

	rest := new(bytes.Buffer)
	rest.ReadFrom(r)
	restString := strings.TrimSpace(rest.String())
	switch {
	case restString[0] == '"':
		s, err := strconv.Unquote(restString)
		if err != nil {
			return err
		}
		value = s
	case unicode.IsDigit(rune(restString[0])):
		i, err := strconv.ParseInt(restString, 10, 64)
		if err != nil {
			return err
		}
		value = i
	case restString == "true":
		value = true
	case restString == "false":
		value = false
	default:
		value = restString
	}

	entry.Key = key
	entry.Value = value
	return nil
}
