package skiplist

import (
	"math/rand"
	"sort"
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

func TestSkiplist(t *testing.T) {
	sl := New(WithMaxLevel(defaultLevel))
	nums := []key{12, 6, 17, 21, 3, 7, 9, 26, 25, 19}
	toRemove := []key{21, 9, 25}
	for _, n := range nums {
		sl.Insert(n, n*2)
	}
	if sl.Length() != len(nums) {
		t.Errorf("Insert: expected len %v, actual len %v", len(nums), sl.Length())
	}
	for _, n := range toRemove {
		v := sl.Remove(n)
		i, ok := v.(key)
		if !ok || !i.Equal(n*2) {
			t.Errorf("Remove: expected %v, actual %v", n, v)
		}
	}
	if sl.Length() != len(nums) {
		t.Errorf("Remove: expected len %v, actual len %v", len(nums), sl.Length())
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
	sl.Insert(key(12), key(12))
	// Delete already remove
	v := sl.Remove(key(21))
	if v != nil {
		t.Errorf("Remove: already removed, expected nil, actual %v", v)
	}
	for _, n := range nums {
		v := sl.Search(n)
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
	t.Log("\n" + sl.String())
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

type sliceNode struct {
	Key   KeyType
	Value interface{}
}

func sliceSearch(slice []sliceNode, target key) {
	for _, n := range slice {
		if n.Key.Equal(target) {
			return
		}
	}
}

func BenchmarkSkiplistSearch(b *testing.B) {
	benches := []struct {
		name     string
		totalNum int
	}{
		{"dense 1k", 1000},
		{"dense 10k", 10000},
	}
	for _, bb := range benches {
		b.Run(bb.name+"/skiplist", func(b *testing.B) {
			nums := randNums(bb.totalNum)
			sl := New(WithMaxLevel(15))
			for _, n := range nums {
				sl.Insert(n, n)
			}
			targets := randTargets(bb.totalNum, b.N)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = sl.Search(targets[i])
			}
		})
		b.Run(bb.name+"/slice", func(b *testing.B) {
			nums := randNums(bb.totalNum)
			slice := make([]sliceNode, 0)
			for _, n := range nums {
				slice = append(slice, sliceNode{Key: n, Value: n})
			}
			sort.Slice(slice, func(i, j int) bool { return slice[i].Key.Less(slice[j].Key) })
			targets := randTargets(bb.totalNum, b.N)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sliceSearch(slice, targets[i])
			}
		})
	}
}
