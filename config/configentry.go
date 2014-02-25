package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/kourge/goit/core"
)

// A ConfigEntry represents a key-value pair in a Git config file. The key is an
// alphanumeric string, and the value is a decimal integer, a boolean value, or
// a string.
type ConfigEntry struct {
	Key   string
	Value interface{}
}

var _ core.EncodeDecoder = &ConfigEntry{}

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
