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
			t.Error("Validate: the adt does not hold expected properties.")
		}
	}
	for _, n := range toRemove {
		t.Logf("Delete key %v", n)
		v := adt.Delete(n)
		i, ok := v.(key)
		if !ok || !i.Equal(n*2) {
			t.Errorf("Remove: expected %v, actual %v", n, v)
		}
		if x, ok := adt.(interface{ String() string }); ok {
			t.Log("\n" + x.String())
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
		if n == 12 {
			i, ok := v.(key)
			if !ok || !i.Equal(key(12)) {
				t.Errorf("Search: expected 12, actual %v", v)
			}
			continue
		}
		if inRemove(n) {
			if v != nil {
				t.Errorf("Search: expected nil, actual %v", v)
			}
		} else {
			i, ok := v.(key)
			if !ok || !i.Equal(n*2) {
				t.Errorf("Search: expected %v, actual %v", n*2, v)
			}
		}
	}
	if x, ok := adt.(interface{ String() string }); ok {
		t.Log("\n" + x.String())
	}
	if x, ok := adt.(Validate); ok {
		if !x.Validate() {
			t.Error("Validate: the adt does not hold expected properties.")
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

func XBenchSearch(b *testing.B, f func() ADT) {
	benches := []struct {
		name     string
		totalNum int
		sparse   bool
	}{
		{name: "dense 1k", totalNum: 1000},
		{name: "dense 10k", totalNum: 10000},
		{name: "sparse 1k", totalNum: 1000, sparse: true},
		{name: "sparse 10k", totalNum: 10000, sparse: true},
	}
	for _, bb := range benches {
		adt := f()
		b.Run(bb.name+"/"+reflect.TypeOf(adt).Elem().Name(), func(b *testing.B) {
			nums := randNums(bb.totalNum)
			for _, n := range nums {
				adt.Insert(n, n)
			}
			var validate = func() {
				if x, ok := adt.(Validate); ok {
					if !x.Validate() {
						b.Error("Validate: the adt does not hold expected properties.")
					}
				}
			}
			validate()
			targets := randTargets(bb.totalNum, b.N)
			if bb.sparse {
				toRemove := randTargets(bb.totalNum, bb.totalNum/2)
				for _, n := range toRemove {
					adt.Delete(n)
				}
			}
			validate()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = adt.Search(targets[i])
			}
		})
	}
}

func XBenchInsert(b *testing.B, f func() ADT) {
	benches := []struct {
		name     string
		totalNum int
		sparse   bool
		print    bool
	}{
		{name: "dense 1k", totalNum: 1000, print: true},
		{name: "dense 10k", totalNum: 10000},
		{name: "sparse 1k", totalNum: 1000, sparse: true},
		{name: "sparse 10k", totalNum: 10000, sparse: true},
	}
	for _, bb := range benches {
		adt := f()
		b.Run(bb.name+"/"+reflect.TypeOf(adt).Elem().Name(), func(b *testing.B) {
			var nums []key
			if !bb.sparse {
				nums = randNums(bb.totalNum)
			} else {
				nums = randNums(bb.totalNum * 2)[:bb.totalNum]
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				n := nums[i%len(nums)]
				adt.Insert(n, nil)
			}
			if x, ok := adt.(Validate); ok {
				if !x.Validate() {
					b.Error("Validate: the adt does not hold expected properties.")
					if x, ok := adt.(interface{ String() string }); ok && bb.print {
						b.Log("\n" + x.String())
					}
				}
			}
		})
	}
}

func XBenchDelete(b *testing.B, f func() ADT) {
	benches := []struct {
		name     string
		totalNum int
		sparse   bool
	}{
		{name: "dense 1k", totalNum: 1000},
		{name: "dense 10k", totalNum: 10000},
		{name: "sparse 1k", totalNum: 1000, sparse: true},
		{name: "sparse 10k", totalNum: 10000, sparse: true},
	}
	for _, bb := range benches {
		adt := f()
		b.Run(bb.name+"/"+reflect.TypeOf(adt).Elem().Name(), func(b *testing.B) {
			numBuckets := b.N / bb.totalNum
			if numBuckets < 1 {
				numBuckets = 1
			}
			buckets := make([]ADT, numBuckets)
			for i := 0; i < numBuckets; i++ {
				buckets[i] = f()
			}
			var nums []key
			var targets []key
			if !bb.sparse {
				nums = randNums(bb.totalNum)
				targets = randTargets(bb.totalNum, b.N)
			} else {
				nums = randNums(bb.totalNum * 2)[:bb.totalNum]
				targets = randTargets(bb.totalNum*2, b.N)
			}
			for _, n := range nums {
				for _, adt := range buckets {
					adt.Insert(n, n)
				}
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = buckets[i%len(buckets)].Delete(targets[i])
			}
		})
	}
}
