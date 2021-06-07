package adt_test

import (
	. "github.com/atriw/lib/golib/adt"
	"testing"
)

type sliceEntry struct {
	key   Key
	value interface{}
}

type slice struct {
	s []*sliceEntry
}

func (s *slice) Insert(key Key, value interface{}) {
	s.s = append(s.s, &sliceEntry{key: key, value: value})
}

func (s *slice) Search(key Key) interface{} {
	for _, e := range s.s {
		if e.key.Equal(key) {
			return e.value
		}
	}
	return nil
}

func (s *slice) Delete(key Key) interface{} {
	for i, e := range s.s {
		if e.key.Equal(key) {
			s.s = append(s.s[:i], s.s[i+1:]...)
			return e.value
		}
	}
	return nil
}

func (s *slice) Length() int {
	return len(s.s)
}

func TestSlice(t *testing.T) {
	s := &slice{}
	XTestADT(t, s)
}

func BenchmarkSliceSearch(b *testing.B) {
	s := &slice{}
	XBenchSearch(b, s)
}
