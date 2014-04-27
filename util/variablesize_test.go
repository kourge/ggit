package util

import (
	"bytes"
	"fmt"
)

func ExampleVariableSize() {
	magic := bytes.NewReader([]byte{0x80, 0xd5, 0x7f})
	v := NewVariableSize(0)
	v.Decode(magic)
	fmt.Printf("%v", v)

	// Output:
	// 2091648
}
