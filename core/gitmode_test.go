package core

import (
	"testing"
)

func TestGitmodeStringRegular(t *testing.T) {
	var mode GitMode = GitModeRegular | GitModeReadWritable
	var actual string = mode.String()
	var expected string = "100644"

	if actual != expected {
		t.Errorf("mode.String() = %v, want %v", actual, expected)
	}
}

func TestGitmodeStringDir(t *testing.T) {
	var mode GitMode = GitModeDir | GitModeNullPerm
	var actual string = mode.String()
	var expected string = "040000"

	if actual != expected {
		t.Errorf("mode.String() = %v, want %v", actual, expected)
	}
}
