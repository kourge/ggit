package core

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type field struct {
	Name  string
	Value EncodeDecoder
}

var _ EncodeDecoder = &field{}

func (f *field) Reader() io.Reader {
	return io.MultiReader(
		strings.NewReader(f.Name),
		bytes.NewReader([]byte{' '}),
		f.Value.Reader(),
	)
}

func (f *field) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)

	if name, err := r.ReadString(byte(' ')); err != nil {
		return err
	} else {
		f.Name = name[:len(name)-1]
	}

	s := &StringCoder{}
	if err := s.Decode(r); err != nil {
		return err
	} else {
		f.Value = s
	}

	return nil
}

type fieldslice []field

var _ Encoder = fieldslice{}

func (fs fieldslice) Readers() []io.Reader {
	var message io.Reader = nil

	n := len(fs)
	last := n-1
	size := n*2-1

	if lastField := fs[last]; lastField.Name == "message" {
		message = lastField.Value.Reader()
		fs = fs[:last]
		last -= 1
		size = (n-1)*2-1+3
	}

	readers := make([]io.Reader, size)
	for i, f := range fs {
		readers[i*2] = f.Reader()
		if i != last {
			readers[i*2+1] = bytes.NewReader([]byte{'\n'})
		}
	}

	if message != nil {
		readers[last*2+1] = bytes.NewReader([]byte("\n\n"))
		readers[last*2+2] = message
		readers[last*2+3] = bytes.NewReader([]byte{'\n'})
	}

	return readers
}

func (fs fieldslice) Reader() io.Reader {
	return io.MultiReader(fs.Readers()...)
}

func fieldsliceDecode(reader io.Reader) (fields fieldslice, err error) {
	r := bufio.NewReader(reader)

	reachedEof := false
	for !reachedEof {
		field := &field{}
		line, err := r.ReadBytes(byte('\n'))

		if err == io.EOF {
			reachedEof = true
		} else if err != nil {
			return nil, err
		} else {
			line = line[:len(line)-1]
		}

		if err := field.Decode(bytes.NewReader(line)); err == io.EOF {
			reachedEof = true
		} else if err != nil {
			return nil, err
		}

		fields = append(fields, *field)
	}

	return
}
