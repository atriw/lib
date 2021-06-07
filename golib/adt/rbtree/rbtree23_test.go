package rbtree_test

import (
	"testing"

	"github.com/atriw/lib/golib/adt"
	. "github.com/atriw/lib/golib/adt/rbtree"
)

func TestRBTree23(t *testing.T) {
	rbt := New23()
	adt.XTestADT(t, rbt)
}

func BenchmarkRBTree23Search(b *testing.B) {
	rbt := New23()
	adt.XBenchSearch(b, rbt)
}
