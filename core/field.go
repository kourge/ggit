package core

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
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
	last := n - 1
	size := n * 2

	if lastField := fs[last]; lastField.Name == "message" {
		message = lastField.Value.Reader()
		fs = fs[:last]
		last -= 1
		size = (n-1)*2 + 3
	}

	readers := make([]io.Reader, size)
	for i, f := range fs {
		readers[i*2] = f.Reader()
		readers[i*2+1] = bytes.NewReader([]byte{'\n'})
	}

	if message != nil {
		readers[last*2+2] = bytes.NewReader([]byte{'\n'})
		readers[last*2+3] = message
		readers[last*2+4] = bytes.NewReader([]byte{'\n'})
	}

	return readers
}

func (fs fieldslice) Reader() io.Reader {
	return io.MultiReader(fs.Readers()...)
}

type fieldcoder interface {
	Decoder
	loadFields(fields fieldslice) error
	setBuffer(bytes []byte)
}

func decodeFields(object fieldcoder, reader io.Reader) error {
	buffer := new(bytes.Buffer)
	r := io.TeeReader(reader, buffer)

	if fields, err := fieldsliceDecode(r); err != nil {
		return err
	} else if err := object.loadFields(fields); err != nil {
		return err
	}

	object.setBuffer(buffer.Bytes())
	return nil
}

func fieldsliceDecode(reader io.Reader) (fields fieldslice, err error) {
	r := bufio.NewReader(reader)

	reachedEof := false
	hasMessage := false
	for !reachedEof {
		f := &field{}
		line, err := r.ReadBytes(byte('\n'))

		if err == io.EOF {
			reachedEof = true
		} else if err != nil {
			return fields, err
		} else if len(line) > 0 && !hasMessage {
			line = line[:len(line)-1]
		}

		if len(line) == 0 {
			hasMessage = true
			continue
		} else if !hasMessage {
			if err := f.Decode(bytes.NewReader(line)); err == io.EOF {
				reachedEof = true
			} else if err != nil {
				return fields, err
			}
			fields = append(fields, *f)
		} else {
			rest, err := ioutil.ReadAll(r)
			if err != nil {
				return fields, err
			}
			line = append(line, rest...)
			message := strings.TrimRightFunc(string(line), unicode.IsSpace)
			fields = append(fields, field{"message", &StringCoder{message}})
		}
	}

	return
}
