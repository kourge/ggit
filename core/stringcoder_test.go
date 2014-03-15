package core

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var (
	_fixtureStringCoder *StringCoder = &StringCoder{"foobar"}
)

func TestStringCoder_Reader(t *testing.T) {
	var actual []byte
	var expected []byte = []byte(_fixtureStringCoder.string)

	actual, err := ioutil.ReadAll(_fixtureStringCoder.Reader())
	if err != nil {
		t.Errorf("got error %v while reading from Reader", err)
	}

	if !bytes.Equal(actual, expected) {
		t.Error("stream.Reader() did not generate same byte sequence")
	}
}

func TestStringCoder_Decode(t *testing.T) {
	var actual *StringCoder = &StringCoder{""}
	var expected *StringCoder = _fixtureStringCoder

	err := actual.Decode(bytes.NewBufferString(_fixtureStringCoder.string))
	if err != nil {
		t.Errorf("stringcoder.Decode() returned error %v", err)
	}

	if *actual != *expected {
		t.Errorf("stringcoder.string = %v, want %v", *actual, *expected)
	}
}
