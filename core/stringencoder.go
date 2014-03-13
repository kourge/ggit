package core

import (
	"io"
	"strings"
)

// StringEncoder wraps a string so that it implements the Encoder interface.
type StringEncoder string

var _ Encoder = StringEncoder("")

// Reader returns an io.Reader that yields the underlying string as bytes.
func (s StringEncoder) Reader() io.Reader {
	return strings.NewReader(string(s))
}
