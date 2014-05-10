package core

import (
	"bytes"
	"testing"
)

func TestGitmode_String_Regular(t *testing.T) {
	var mode GitMode = GitModeRegular | GitModeReadWritable
	var actual string = mode.String()
	var expected string = "100644"

	if actual != expected {
		t.Errorf("mode.String() = %v, want %v", actual, expected)
	}
}

func TestGitmode_String_Dir(t *testing.T) {
	var mode GitMode = GitModeDir | GitModeNullPerm
	var actual string = mode.String()
	var expected string = "040000"

	if actual != expected {
		t.Errorf("mode.String() = %v, want %v", actual, expected)
	}
}

func TestGitmode_Reader_Regular(t *testing.T) {
	var mode GitMode = GitModeRegular | GitModeReadWritable
	var actual []byte
	var expected []byte = []byte("100644")

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(mode.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Errorf("mode.Reader() gave %v, want %v", actual, expected)
	}
}

func TestGitmode_Reader_Dir(t *testing.T) {
	var mode GitMode = GitModeDir | GitModeNullPerm
	var actual []byte
	var expected []byte = []byte("40000")

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(mode.Reader())
	actual = buffer.Bytes()

	if !bytes.Equal(actual, expected) {
		t.Errorf("mode.String() gave %v, want %v", actual, expected)
	}
}
