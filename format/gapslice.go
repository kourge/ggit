package format

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
)

// A GapSlice is akin to a gap buffer data structure that is commonly used by
// text editors. Its base structure is a linked list where each node is either
// a value node or a slice node. The GapSlice starts out as a single slice node.
//
// A subsequent insertion, if occurring at an index before or after the slice,
// is done by either prepending a value node before or append a value node after
// the slice node. If the insertion occurs in the middle of a slice node, then
// it is split into two slice nodes, and a value node is inserted between the
// two slice nodes.
//
// Similarly, a removal, if occurring at an index before or after a slice, is
// done by simply reslicing the slice node. If the removal occurs in the middle
// of a slice node, then it is split into two slice nodes, with the first slice
// node resliced to exclude the removed index.
type GapSlice struct {
	list *list.List
	size int
}

// The type that is stored in a linked list node when it is a slice node.
type gapSliceChunk []interface{}

// NewGapSlice constructs a GapSlice based on the given existing slice. The
// underlying slice itself will never be modified.
func NewGapSlice(existing []interface{}) *GapSlice {
	s := &GapSlice{list.New(), len(existing)}
	s.list.PushBack(gapSliceChunk(existing))
	return s
}

// Len returns the total number of items in this GapSlice.
func (s *GapSlice) Len() int {
	return s.size
}

// locate returns the node on which the item at the logical index should reside.
// If the index lands on a slice node, the offset is the physical index relative
// to the slice. If the index lands on a value node, then the offset is
// meaningless and set to 0.
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

// Get attempts to fetch the item logically residing at index i. If the index
// is within bounds, then v is the item retried and exists is set to true. If
// the index is out of bounds, then v is set to nil and exists will be false.
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

// Set is a convenience method for removing an item at index i and then
// inserting another item at the same index. Note that while this logically
// changes the GapSlice, the underlying slice used to initialize the GapSlice
// will not be mutated.
func (s *GapSlice) Set(i int, v interface{}) (success bool) {
	s.RemoveAt(i)
	return s.InsertAt(i, v)
}

// PushBack is a convenience method for inserting an item at index i where i is
// the length of this GapSlice, i.e. it appends an item.
func (s *GapSlice) PushBack(v interface{}) {
	s.InsertAt(s.size, v)
}

// InsertAt attempts to insert the value v at the given logical index i. It
// returns false if i is out of bounds. Note that a value of i equal to the
// length of this GapSlice is not considered out of bounds; it is handled as
// a special case of appending to this GapSlice.
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

	list.InsertBefore(v, e)
	s.size += 1
	return true
}

// RemoveAt attempts to remove the value v at the given logical index i. The
// return value v is the item that was found at i. If i is out of bounds, then
// v is set to nil and exists will be false. Otherwise, exists will be true.
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

// GapSliceWalkFunc is the type of the function called for each item visited by
// Walk.
type GapSliceWalkFunc func(v interface{}) error

// Walk iterates through this GapSlice, calling walkFn for each item visited by
// Walk.
//
// If walkFun returns an error, the iteration is immediately halted, and in turn,
// that error is returned by Walk.
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

// collect consolidates all items in this GapSlice onto a newly-allocated single
// slice. The existing nodes are untouched.
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

// Pack consolidates all items in this GapSlice into a single slice node. The
// slice is then returned.
func (s *GapSlice) Pack() []interface{} {
	packed := s.collect()
	s.list = new(list.List)
	s.PushBack(gapSliceChunk(packed))
	return packed
}

// String formats this GapSlice as a space-separated list of items surrounded
// by parentheses.
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

// GoString formats this GapSlice as a space-separated list of nodes surrounded
// by parentheses. Value nodes are represented by the GoString format of the
// value itself, while slice nodes are presented as a space-separated list of
// values in GoString format surrounded by square brackets.
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
