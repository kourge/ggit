package config

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/kourge/goit/core"
)

// A Section represents a section in a Git config file. It has its own name and
// a list of key-value pairs under it.
type Section struct {
	Name    string
	Entries []Entry
}

var _ core.Encoder = Section{}

func (section Section) bytesBuffer() *bytes.Buffer {
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
// first line. Subsequent lines are the underlying Entry structs serialized in
// order, each indented by a single horizontal tab rune '\t' and each separated
// by a new line rune '\n'.
func (section Section) Reader() io.Reader {
	return section.bytesBuffer()
}

func (section Section) lazyReader() io.Reader {
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

// String returns a string that is the result of draining the io.Reader returned
// by Reader().
func (section Section) String() string {
	return section.bytesBuffer().String()
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
	entries := make([]Entry, 0)

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

		entry := &Entry{}
		entry.Decode(strings.NewReader(line))
		entries = append(entries, *entry)
	}

	section.Entries = entries
	return nil
}

// WalkSectionFunc is the type of the function called for each key-value pair
// iterated by Walk. Its return value is inspected to determine whether the
// iteration is short-circuited; when a WalkSectionFunc returns false, Walk
// breaks out of the iteration loop.
type WalkSectionFunc func(key string, value interface{}) (stop bool)

// Walk iterates though the entries contained in this section, calling walkFn
// for each key-value pair.
func (section Section) Walk(walkFn WalkSectionFunc) {
	for _, entry := range section.Entries {
		stop := walkFn(entry.Key, entry.Value)
		if stop {
			break
		}
	}
}

// Len returns the number of key-value pairs in this section.
func (section Section) Len() int {
	return len(section.Entries)
}

// Get returns a pair in the form of (value, exists) given a key. If there
// exists a value for the given key, exists will be true. Otherwise, it will be
// false and value will be nil.
func (section Section) Get(key string) (value interface{}, exists bool) {
	for _, entry := range section.Entries {
		if entry.Key == key {
			return entry.Value, true
		}
	}

	return
}

// Set either adds a key-value pair or, if the key already exists, associates
// the given key to the provided value instead. If the former case is true,
// Set returns true. Otherwise, it returns false.
func (section *Section) Set(key string, value interface{}) (added bool) {
	for i, entry := range section.Entries {
		if entry.Key == key {
			section.Entries[i].Value = value
			return false
		}
	}

	section.Entries = append(section.Entries, Entry{key, value})
	return true
}

// Del removes the key-value pair for a given key or does nothing if no
// key-value pair exists for the given key. If the former case is true, Del
// returns true. Otherwise, it returns false.
func (section *Section) Del(key string) (deleted bool) {
	index := -1
	for i, entry := range section.Entries {
		if entry.Key == key {
			index = i
			break
		}
	}

	deleted = index != -1
	if deleted {
		for i, length := index, len(section.Entries); i < length; i++ {
			if i != length - 1 {
				section.Entries[i] = section.Entries[i + 1]
			}
		}
		section.Entries = section.Entries[:len(section.Entries)-1]
	}

	return
}
