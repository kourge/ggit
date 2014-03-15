package core

import (
	"bytes"
	"io"
	"strings"
)

// StringCoder wraps a string so that it implements the Encoder interface.
type StringCoder struct {
	string
}

var _ EncodeDecoder = &StringCoder{""}

// Reader returns an io.Reader that yields the underlying string as bytes.
func (s *StringCoder) Reader() io.Reader {
	return strings.NewReader(s.string)
}

// Decode reads from an io.Reader and treats the bytes read as an UTF-8 string.
func (s *StringCoder) Decode(reader io.Reader) error {
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(reader); err != nil {
		return err
	}
	s.string = buffer.String()
	return nil
}
