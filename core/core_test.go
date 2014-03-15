package core

import (
	"errors"
	"testing"
)

var (
	_fixtureMessage string = "PC Load Letter"
)

func TestDie(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			expected := _fixtureMessage
			if r != expected {
				t.Errorf("Expected die() to panic error message: %v, got %v instead", expected, r)
			}
		} else {
			t.Error("Expected die() to panic")
		}
	}()
	Die(errors.New(_fixtureMessage))
}
