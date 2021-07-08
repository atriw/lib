package bptree_test

import (
	"testing"

	"github.com/atriw/lib/golib/adt"
	. "github.com/atriw/lib/golib/adt/bptree"
)

func TestBPTree(t *testing.T) {
	bpt := New(WithOrder(4))
	adt.XTestADT(t, bpt)
}

func BenchmarkBPTreeSearch(b *testing.B) {
	adt.XBenchSearch(b, func() adt.ADT { return New(WithOrder(11)) })
}

func BenchmarkBPTreeInsert(b *testing.B) {
	adt.XBenchInsert(b, func() adt.ADT { return New(WithOrder(11)) })
}

func BenchmarkBPTreeDelete(b *testing.B) {
	adt.XBenchDelete(b, func() adt.ADT { return New(WithOrder(11)) })
}
