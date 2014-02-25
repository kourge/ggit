package config

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/kourge/goit/core"
)

// A ConfigSection represents a section in a Git config file. It has its own
// name and a list of key-value pairs under it.
type ConfigSection struct {
	Name    string
	Entries []ConfigEntry
}

var _ core.Encoder = ConfigSection{}

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
