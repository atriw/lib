package rbtree_test

import (
	"testing"

	"github.com/atriw/lib/golib/adt"
	. "github.com/atriw/lib/golib/adt/rbtree"
)

func TestLLRBTree(t *testing.T) {
	rbt := NewLL()
	adt.XTestADT(t, rbt)
}

func BenchmarkLLRBTreeSearch(b *testing.B) {
	adt.XBenchSearch(b, func() adt.ADT { return NewLL() })
}
