package skiplist_test

import (
	"testing"

	"github.com/atriw/lib/golib/adt"
	. "github.com/atriw/lib/golib/adt/skiplist"
)

func TestSkiplist(t *testing.T) {
	sl := New()
	adt.XTestADT(t, sl)
}

func BenchmarkSkiplistSearch(b *testing.B) {
	sl := New(WithMaxLevel(15))
	adt.XBenchSearch(b, sl)
}
