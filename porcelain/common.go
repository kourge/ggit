package porcelain

import (
	"os"
)

type fileFunc func(f *os.File) error

var fileNop fileFunc = func(f *os.File) error {
	return nil
}

type fileCreation struct {
	Name string
	Do   fileFunc
}
