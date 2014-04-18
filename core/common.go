package core

import (
	"errors"
	"fmt"
)

// Die takes an error and panics with the error's string.
func Die(err error) {
	panic(err.Error())
}

// Errorf is a wrapper around errors.New(fmt.Sprintf(format, rest...)).
func Errorf(format string, rest ...interface{}) error {
	return errors.New(fmt.Sprintf(format, rest...))
}
