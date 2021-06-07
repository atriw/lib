package adt

import (
	"math/rand"
	"reflect"
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

func XTestADT(t *testing.T, adt ADT) {
	nums := []key{12, 6, 17, 21, 3, 7, 9, 26, 25, 19}
	toRemove := []key{21, 9, 25}
	for _, n := range nums {
		adt.Insert(n, n*2)
	}
	if adt.Length() != len(nums) {
		t.Errorf("Insert: expected len %v, actual len %v", len(nums), adt.Length())
	}
	if x, ok := adt.(interface{ String() string }); ok {
		t.Log("\n" + x.String())
	}
	if x, ok := adt.(Validate); ok {
		if !x.Validate() {
			t.Error("Validate: the tree doesn't hold red-black-tree's properties")
		}
	}
	for _, n := range toRemove {
		v := adt.Delete(n)
		i, ok := v.(key)
		if !ok || !i.Equal(n*2) {
			t.Errorf("Remove: expected %v, actual %v", n, v)
		}
	}
	if adt.Length() != len(nums)-len(toRemove) {
		t.Errorf("Remove: expected len %v, actual len %v", len(nums)-len(toRemove), adt.Length())
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
	adt.Insert(key(12), key(12))
	// Delete already remove
	v := adt.Delete(key(21))
	if v != nil {
		t.Errorf("Remove: already removed, expected nil, actual %v", v)
	}
	for _, n := range nums {
		v := adt.Search(n)
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
	if x, ok := adt.(interface{ String() string }); ok {
		t.Log("\n" + x.String())
	}
	if x, ok := adt.(Validate); ok {
		if !x.Validate() {
			t.Error("Validate: the tree doesn't hold red-black-tree's properties")
		}
	}
}

func randNums(n int) []key {
	var nums []key
	for i := 0; i < n; i++ {
		nums = append(nums, key(i))
	}
	rand.Shuffle(n, func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
	return nums
}

func randTargets(n, total int) []key {
	var targets []key
	for i := 0; i < total; i++ {
		targets = append(targets, key(rand.Intn(n)))
	}
	return targets
}

func XBenchSearch(b *testing.B, adt ADT) {
	benches := []struct {
		name     string
		totalNum int
	}{
		{"dense 1k", 1000},
		{"dense 10k", 10000},
	}
	for _, bb := range benches {
		b.Run(bb.name+"/"+reflect.TypeOf(adt).Elem().Name(), func(b *testing.B) {
			nums := randNums(bb.totalNum)
			for _, n := range nums {
				adt.Insert(n, n)
			}
			targets := randTargets(bb.totalNum, b.N)
			if x, ok := adt.(Validate); ok {
				if !x.Validate() {
					b.Error("Validate: the adt does not hold expected properties.")
				}
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = adt.Search(targets[i])
			}
		})
	}
}
