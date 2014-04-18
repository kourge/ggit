package config

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/kourge/ggit/core"
)

// A Config represents a INI-style Git config file, composed of multiple
// sections, each of them with their own list of key-value pairs. Includes are
// not supported at the moment.
type Config map[string]Section

var _ core.EncodeDecoder = Config{}

// Reader returns an io.Reader that, when read, prints out each section,
// each separated by a blank line.
func (config Config) Reader() io.Reader {
	readers := make([]io.Reader, len(config)*2)

	i := 0
	for _, section := range config {
		readers[i*2] = section.Reader()
		readers[i*2+1] = bytes.NewReader([]byte{'\n'})
		i += 1
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
func (config Config) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)
	lines := new(bytes.Buffer)

	flush := func(realloc bool) error {
		section := &Section{}
		if err := section.Decode(lines); err != nil {
			return err
		}
		config[section.Name] = *section

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
