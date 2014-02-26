package config

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// A OrderedSection represents a section in a Git config file. It has its own
// name and a list of key-value pairs under it.
type OrderedSection struct {
	name    string
	entries []Entry
}

var _ Section = &OrderedSection{}

// NewOrderedSection returns a OrderedSection named name that contains the items
// in entries.
func NewOrderedSection(name string, entries ...Entry) Section {
	return &OrderedSection{name: name, entries: entries}
}

// NewEmptyOrderedSection returns an blank OrderedSection with no name and no
// entries.
func NewEmptyOrderedSection() Section {
	return &OrderedSection{}
}

// Name returns the name of this section.
func (section *OrderedSection) Name() string {
	return section.name
}

// SetName sets the name of this section.
func (section *OrderedSection) SetName(name string) {
	section.name = name
}

func (section *OrderedSection) bytesBuffer() *bytes.Buffer {
	buffer := new(bytes.Buffer)

	buffer.WriteRune('[')
	buffer.WriteString(section.name)
	buffer.WriteRune(']')
	buffer.WriteRune('\n')

	for _, entry := range section.entries {
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
func (section *OrderedSection) Reader() io.Reader {
	return section.bytesBuffer()
}

func (section *OrderedSection) lazyReader() io.Reader {
	offset := 3
	readers := make([]io.Reader, len(section.entries)*3+offset)
	readers[0] = bytes.NewReader([]byte{'['})
	readers[1] = strings.NewReader(section.name)
	readers[2] = bytes.NewReader([]byte{']', '\n'})

	for i, entry := range section.entries {
		readers[i*3+offset+0] = bytes.NewReader([]byte{'\t'})
		readers[i*3+offset+1] = entry.Reader()
		readers[i*3+offset+2] = bytes.NewReader([]byte{'\n'})
	}

	return io.MultiReader(readers...)
}

// String returns a string that is the result of draining the io.Reader returned
// by Reader().
func (section *OrderedSection) String() string {
	return section.bytesBuffer().String()
}

// Decode parses bytes into a OrderedSection. This stream of bytes is assumed to
// contain only a single section. If multiple sections appear, all entries will
// be treated as if they were in a single section, and the last seen section's
// name is considered to be this single section's name.
//
// Lines that start with the pound sign rune '#' or the semicolon rune ';' will
// be treated as comment and ignored. Completely whitespace lines or blank lines
// will also be ignored.
func (section *OrderedSection) Decode(reader io.Reader) error {
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
			section.name = strings.TrimSpace(line[1:i])
			continue
		}

		entry := &Entry{}
		entry.Decode(strings.NewReader(line))
		entries = append(entries, *entry)
	}

	section.entries = entries
	return nil
}

// Walk iterates though the entries contained in this section, calling walkFn
// for each key-value pair.
func (section *OrderedSection) Walk(walkFn WalkSectionFunc) {
	for _, entry := range section.entries {
		stop := walkFn(entry.Key, entry.Value)
		if stop {
			break
		}
	}
}

// Len returns the number of key-value pairs in this section.
func (section *OrderedSection) Len() int {
	return len(section.entries)
}

// Get returns a pair in the form of (value, exists) given a key. If there
// exists a value for the given key, exists will be true. Otherwise, it will be
// false and value will be nil.
func (section *OrderedSection) Get(key string) (value interface{}, exists bool) {
	for _, entry := range section.entries {
		if entry.Key == key {
			return entry.Value, true
		}
	}

	return
}

// Set either adds a key-value pair or, if the key already exists, associates
// the given key to the provided value instead. If the former case is true, Set
// returns true. Otherwise, it returns false.
func (section *OrderedSection) Set(key string, value interface{}) (added bool) {
	for i, entry := range section.entries {
		if entry.Key == key {
			section.entries[i].Value = value
			return false
		}
	}

	section.entries = append(section.entries, Entry{key, value})
	return true
}

// Del removes the key-value pair for a given key or does nothing if no
// key-value pair exists for the given key. If the former case is true, Del
// returns true. Otherwise, it returns false.
func (section *OrderedSection) Del(key string) (deleted bool) {
	index := -1
	for i, entry := range section.entries {
		if entry.Key == key {
			index = i
			break
		}
	}

	deleted = index != -1
	if deleted {
		for i, length := index, len(section.entries); i < length; i++ {
			if i != length - 1 {
				section.entries[i] = section.entries[i + 1]
			}
		}
		section.entries = section.entries[:len(section.entries)-1]
	}

	return
}
