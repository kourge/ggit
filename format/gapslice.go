package format

import (
	"bytes"
	"container/list"
	"fmt"
)

type GapSlice struct {
	list *list.List
	size int
}

type gapSliceChunk []interface{}

func NewGapSlice(existing []interface{}) *GapSlice {
	s := &GapSlice{list.New(), len(existing)}
	s.list.PushBack(gapSliceChunk(existing))
	return s
}

func (s *GapSlice) Len() int {
	return s.size
}

func (s *GapSlice) locate(i int) (elem *list.Element, offset int) {
	if i >= s.size {
		return nil, 0
	}

	remaining := i
	list := s.list
	for e := list.Front(); e != nil; e = e.Next() {
		switch value := e.Value.(type) {
		case gapSliceChunk:
			if chunkSize := len(value); remaining < chunkSize {
				return e, remaining
			} else {
				remaining -= chunkSize
			}
		default:
			if remaining == 0 {
				return e, 0
			} else {
				remaining -= 1
			}
		}
	}

	panic(Errorf("could not find element for index %d in GapSlice", i))
}

func (s *GapSlice) Get(i int) (v interface{}, exists bool) {
	if i >= s.size {
		return nil, false
	}

	e, offset := s.locate(i)
	if e == nil {
		return nil, false
	}

	if slice, isSlice := e.Value.(gapSliceChunk); isSlice {
		return slice[offset], true
	}

	return e.Value, true
}

func (s *GapSlice) Set(i int, v interface{}) (success bool) {
	s.RemoveAt(i)
	return s.InsertAt(i, v)
}

func (s *GapSlice) PushBack(v interface{}) {
	s.InsertAt(s.size, v)
}

// InsertAt will always result in split chunks.
func (s *GapSlice) InsertAt(i int, v interface{}) (success bool) {
	if i > s.size {
		return false
	}

	list := s.list

	// Special case: inserting in the very end of this gap slice.
	if i == s.size {
		list.PushBack(v)
		s.size += 1
		return true
	}

	e, offset := s.locate(i)
	if e == nil {
		return false
	}

	if slice, isSlice := e.Value.(gapSliceChunk); isSlice {
		if offset == 0 {
			list.InsertBefore(v, e)
		} else {
			a, b := slice[:offset], slice[offset:]
			e.Value = a
			e = list.InsertAfter(v, e)
			list.InsertAfter(b, e)
		}
		s.size += 1
		return true
	}

	// The given index is before or after an item.
	list.InsertBefore(v, e)
	s.size += 1
	return true
}

func (s *GapSlice) RemoveAt(i int) (v interface{}, exists bool) {
	if i > s.size {
		return nil, false
	}

	list := s.list

	e, offset := s.locate(i)
	if e == nil {
		return nil, false
	}

	// The given index is within a gap slice chunk.
	if slice, isSlice := e.Value.(gapSliceChunk); isSlice {
		v = slice[offset]
		// If the index is at the beginning or end of the chunk, simply reslice.
		// Else, split the chunk.
		switch offset {
		case 0:
			e.Value = slice[1:]
		case len(slice) - 1:
			e.Value = slice[:len(slice)-1]
		default:
			a, b := slice[:offset], slice[offset+1:]
			e.Value = a
			list.InsertAfter(b, e)
		}
		s.size -= 1
		return v, true
	}

	v = list.Remove(e)
	s.size -= 1
	return v, true
}

type GapSliceWalkFunc func(v interface{}) error

func (s *GapSlice) Walk(walkFn GapSliceWalkFunc) error {
	list := s.list
	for e := list.Front(); e != nil; e = e.Next() {
		switch value := e.Value.(type) {
		case gapSliceChunk:
			for _, v := range value {
				if err := walkFn(v); err != nil {
					return err
				}
			}
		default:
			if err := walkFn(value); err != nil {
				return nil
			}
		}
	}
	return nil
}

func (s *GapSlice) collect() []interface{} {
	slice := make([]interface{}, s.Len())
	i := 0
	s.Walk(func(v interface{}) error {
		slice[i] = v
		i += 1
		return nil
	})
	return slice
}

func (s *GapSlice) Pack() []interface{} {
	packed := s.collect()
	s.list = new(list.List)
	s.PushBack(gapSliceChunk(packed))
	return packed
}

func (s *GapSlice) String() string {
	list := s.list
	b := bytes.NewBufferString("(")
	for e := list.Front(); e != nil; e = e.Next() {
		switch v := e.Value.(type) {
		case gapSliceChunk:
			last := len(v) - 1
			for i, item := range v {
				b.WriteString(fmt.Sprintf("%v", item))
				if i != last {
					b.WriteRune(' ')
				}
			}
		default:
			b.WriteString(fmt.Sprintf("%v", v))
		}
		if e != list.Back() {
			b.WriteRune(' ')
		}
	}
	b.WriteRune(')')
	return b.String()
}

func (s *GapSlice) GoString() string {
	list := s.list
	b := bytes.NewBufferString("(")
	for e := list.Front(); e != nil; e = e.Next() {
		if chunk, isChunk := e.Value.(gapSliceChunk); isChunk {
			b.WriteString(fmt.Sprintf("%v", chunk))
		} else {
			b.WriteString(fmt.Sprintf("%#v", e.Value))
		}
		if e != list.Back() {
			b.WriteRune(' ')
		}
	}
	b.WriteRune(')')
	return b.String()
}
