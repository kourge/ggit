package config

import (
	"github.com/kourge/goit/core"
)

type Section interface {
	Name() string
	SetName(name string)
	Len() int
	Get(key string) (value interface{}, exists bool)
	Set(key string, value interface{}) (added bool)
	Del(key string) (deleted bool)
	Walk(WalkSectionFunc)
	core.EncodeDecoder
}

// WalkSectionFunc is the type of the function called for each key-value pair
// iterated by Walk. Its return value is inspected to determine whether the
// iteration is short-circuited; when a WalkSectionFunc returns false, Walk
// breaks out of the iteration loop.
type WalkSectionFunc func(key string, value interface{}) (stop bool)

// EmptySectionMaker is the type of the function called to make a new struct
// that implements Section.
type EmptySectionMaker func() Section

