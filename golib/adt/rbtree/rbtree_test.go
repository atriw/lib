package rbtree_test

import (
	"testing"

	"github.com/atriw/lib/golib/adt"
	. "github.com/atriw/lib/golib/adt/rbtree"
)

func TestRBTree(t *testing.T) {
	rbt := New()
	adt.XTestADT(t, rbt)
}

func BenchmarkRBTreeSearch(b *testing.B) {
	adt.XBenchSearch(b, func() adt.ADT { return New() })
}
