package rbtree_test

import (
	. "github.com/atriw/lib/golib/rbtree"
	"testing"
)

type key int

func (k key) Less(other interface{}) bool {
	i, ok := other.(key)
	return ok && k < i
}

func (k key) Equal(other interface{}) bool {
	i, ok := other.(key)
	return ok && k == i
}

func TestRBTree(t *testing.T) {
	rbt := New()
	nums := []key{12, 6, 17, 21, 3, 7, 9, 26, 25, 19}
	toRemove := []key{21, 9, 25}
	for _, n := range nums {
		rbt.Insert(n, n*2)
	}
	if rbt.Length() != len(nums) {
		t.Errorf("Insert: expected len %v, actual len %v", len(nums), rbt.Length())
	}
	t.Log("\n" + rbt.String())
	for _, n := range toRemove {
		v := rbt.Remove(n)
		i, ok := v.(key)
		if !ok || !i.Equal(n*2) {
			t.Errorf("Remove: expected %v, actual %v", n, v)
		}
	}
	if rbt.Length() != len(nums) {
		t.Errorf("Remove: expected len %v, actual len %v", len(nums), rbt.Length())
	}
	inRemove := func(n key) bool {
		for _, v := range toRemove {
			if n.Equal(v) {
				return true
			}
		}
		return false
	}
	// Insert already exist
	rbt.Insert(key(12), key(12))
	// Delete already remove
	v := rbt.Remove(key(21))
	if v != nil {
		t.Errorf("Remove: already removed, expected nil, actual %v", v)
	}
	for _, n := range nums {
		v := rbt.Search(n)
		if inRemove(n) {
			if v != nil {
				t.Errorf("Search: expected nil, actual %v", v)
			}
		} else {
			i, ok := v.(key)
			if !ok || !i.Equal(n*2) {
				t.Errorf("Search: expected %v, actual %v", n, v)
			}
		}
	}
	t.Log("\n" + rbt.String())
}
