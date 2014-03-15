// The core package contains all the low-level primitives of Git objects and
// provides ways to encode and decode them to and from byte streams.
package core

import (
	"errors"
	"fmt"
)

func Die(err error) {
	panic("fatal: " + err.Error())
}

// Errorf is a wrapper around errors.New(fmt.Sprintf(format, rest...)).
func Errorf(format string, rest ...interface{}) error {
	return errors.New(fmt.Sprintf(format, rest...))
}
