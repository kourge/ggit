package util

import (
	"fmt"
	"testing"
)

func ExampleGapSlice() {
	a := []interface{}{1.0, 3.0, 5.0, 7.0}
	s := NewGapSlice(a)

	fmt.Printf("1: %#v\n", s)

	s.InsertAt(3, 6.0)
	fmt.Printf("2: %#v\n", s)

	s.InsertAt(2, 4.0)
	fmt.Printf("3: %#v\n", s)

	s.InsertAt(1, 2.0)
	fmt.Printf("4: %#v\n", s)

	s.Pack()
	fmt.Printf("5: %#v\n", s)

	// Output:
	// 1: ([1 3 5 7])
	// 2: ([1 3 5] 6 [7])
	// 3: ([1 3] 4 [5] 6 [7])
	// 4: ([1] 2 [3] 4 [5] 6 [7])
	// 5: ([1 2 3 4 5 6 7])
}

func TestGapSlice_Len(t *testing.T) {
	s := NewGapSlice([]interface{}{1.0, 3.0, 5.0, 7.0})
	if s.Len() != 4 {
		t.Errorf("Len of %#v is %d, want %d", s, s.Len(), 4)
	}

	s.InsertAt(1, 2.0)
	if s.Len() != 5 {
		t.Errorf("Len of %#v is %d, want %d", s, s.Len(), 5)
	}

	s.RemoveAt(4)
	if s.Len() != 4 {
		t.Errorf("Len of %#v is %d, want %d", s, s.Len(), 4)
	}
}

func TestGapSlice_Get(t *testing.T) {
	s := NewGapSlice([]interface{}{1.0, 3.0, 5.0, 7.0})
	s.InsertAt(2, 4.0)

	for _, test := range []struct {
		i int
		v float64
	}{{1, 3.0}, {2, 4.0}, {3, 5.0}} {
		if v, exists := s.Get(test.i); !exists {
			t.Errorf("s.Get(%d) should exist", test.i)
		} else if f := v.(float64); f != test.v {
			t.Errorf("s.Get(%d) = %v, want %v", test.i, f, test.v)
		}
	}

	for _, i := range []int{-1, 5} {
		if _, exists := s.Get(i); exists {
			t.Errorf("s.Get(%d) should not exist", i)
		}
	}
}

func TestGapSlice_Set(t *testing.T) {
	a := []interface{}{1.0, 5.0}
	s := NewGapSlice(a)
	s.InsertAt(1, 0.0)
	for _, test := range []struct {
		i int
		v float64
	}{{0, 2.0}, {1, 3.0}, {2, 4.0}} {
		if ok := s.Set(test.i, test.v); !ok {
			t.Errorf("s.Set(%d, %v) = false, want true", test.i, test.v)
		}
	}

	if ok := s.Set(3, 5.0); !ok {
		t.Errorf("s.Set(3, _) should succeed")
	}

	if a[0] != 1.0 || a[1] != 5.0 {
		t.Errorf("s.Set(_, _) mutated underlying slice")
	}
}

func TestGapSlice_PushBack(t *testing.T) {
	a := []interface{}{2.718}
	s := NewGapSlice(a)
	v := 3.14
	s.PushBack(v)

	if v, exists := s.Get(1); !exists {
		t.Errorf("s.PushBack(%v) failed", v)
	} else if v.(float64) != v {
		t.Errorf("s.PushBack(%v) failed", v)
	}

	if a[0] != 2.718 {
		t.Errorf("s.PushBack(_) mutated underlying slice")
	}
}

func TestGapSlice_InsertAt(t *testing.T) {
	s := NewGapSlice([]interface{}{1.0, 3.0, 5.0, 7.0})
	s.InsertAt(3, 6.0)
	s.InsertAt(2, 4.0)
	s.InsertAt(1, 2.0)

	var actual = s.GoString()
	var expected = "([1] 2 [3] 4 [5] 6 [7])"
	if actual != expected {
		t.Errorf("s ended up = %v, want %v", actual, expected)
	}
}

func TestGapSlice_RemoveAt(t *testing.T) {
	s := NewGapSlice([]interface{}{1.0, 3.0, 5.0, 7.0, 9.0})

	for _, test := range []struct {
		i int
		v float64
	}{{4, 9.0}, {1, 3.0}, {0, 1.0}} {
		if v, exists := s.RemoveAt(test.i); !exists {
			t.Errorf("s.RemoveAt(%d) failed", test.i)
		} else if v != test.v {
			t.Errorf("s.RemoveAt(%d) returned %v, want %v", test.i, v, test.v)
		}
	}

	var actual = s.GoString()
	var expected = "([] [5 7])"
	if actual != expected {
		t.Errorf("s ended up = %v, want %v", actual, expected)
	}
}

func TestGapSlice_Walk(t *testing.T) {
	s := NewGapSlice([]interface{}{1.0, 3.0, 5.0, 7.0})
	s.InsertAt(3, 6.0)
	s.InsertAt(2, 4.0)
	s.InsertAt(1, 2.0)

	a := make([]interface{}, 0)
	s.Walk(func(v interface{}) error {
		a = append(a, v)
		return nil
	})

	b := []interface{}{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}
	for i, v := range a {
		if v != b[i] {
			t.Errorf("s.Walk() produced a[%d] = %v, want %v", i, v, b[i])
		}
	}
}

func TestGapSlice_Pack(t *testing.T) {
	s := NewGapSlice([]interface{}{1.0, 3.0, 5.0, 7.0})
	s.InsertAt(3, 6.0)
	s.InsertAt(2, 4.0)
	s.InsertAt(1, 2.0)
	s.Pack()

	var actual = s.GoString()
	var expected = "([1 2 3 4 5 6 7])"
	if actual != expected {
		t.Errorf("s ended up = %v, want %v", actual, expected)
	}
}
