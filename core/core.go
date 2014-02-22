// The core package contains all the low-level primitives of Git objects and
// provides ways to encode and decode them to and from byte streams.
package core

func die(err error) {
	panic("fatal: " + err.Error())
}
