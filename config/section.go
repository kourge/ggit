package config

import (
	"bufio"
	"bytes"
	"io"
	"sort"
	"strings"
)

// A Dict represents a map with string keys and arbitrary-typed values.
type Dict map[string]interface{}

// A Section represents a section in a Git config file. It has its own name and
// a list of key-value pairs under it.
type Section struct {
	Name string
	Dict
}

// Returns an io.Reader that yields the section name in square brackets as the
// first line. Subsequent lines are the underlying Entry structs serialized in
// order, each indented by a single horizontal tab rune '\t' and each separated
// by a new line rune '\n'.
func (section *Section) Reader() io.Reader {
	offset := 3
	readers := make([]io.Reader, len(section.Dict)*3+offset)
	readers[0] = bytes.NewReader([]byte{'['})
	readers[1] = strings.NewReader(section.Name)
	readers[2] = bytes.NewReader([]byte{']', '\n'})

	keys := make([]string, len(section.Dict))
	{
		i := 0
		for k, _ := range section.Dict {
			keys[i] = k
			i++
		}
	}
	sort.Strings(keys)

	for i, k := range keys {
		v := section.Dict[k]
		readers[i*3+offset+0] = bytes.NewReader([]byte{'\t'})
		readers[i*3+offset+1] = (&Entry{k, v}).Reader()
		readers[i*3+offset+2] = bytes.NewReader([]byte{'\n'})
	}

	return io.MultiReader(readers...)
}

// String returns a string that is the result of draining the io.Reader returned
// by Reader().
func (section *Section) String() string {
	var buffer = new(bytes.Buffer)
	buffer.ReadFrom(section.Reader())
	return buffer.String()
}

// Decode parses bytes into a Section. This stream of bytes is assumed to
// contain only a single section. If multiple sections appear, all entries will
// be treated as if they were in a single section, and the last seen section's
// name is considered to be this single section's name.
//
// Lines that start with the pound sign rune '#' or the semicolon rune ';' will
// be treated as comment and ignored. Completely whitespace lines or blank lines
// will also be ignored.
func (section *Section) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)
	section.Dict = make(Dict)
	entry := &Entry{}

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

		entry.Decode(strings.NewReader(line))
		section.Dict[entry.Key] = entry.Value
	}

	return nil
}
