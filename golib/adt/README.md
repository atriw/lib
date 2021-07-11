# Abstract data type

## Interfaces

```go
type Key interface {
	Less(interface{}) bool
	Equal(interface{}) bool
}

type ADT interface {
	Insert(key Key, value interface{})
	Search(key Key) interface{}
	Delete(key Key) interface{}
	Length() int
}
```

## Implementations

- Skiplist
- Red-Black tree
- Left-leaning red-black tree
- B+ tree

## Benchmarks

Skiplist with 15 as max level.

B+ tree with 10 as order.

```bash
goos: linux
goarch: amd64
pkg: github.com/atriw/lib/golib/adt
BenchmarkSearch/dense_1k/slice-12                 709290              1655 ns/op
BenchmarkSearch/dense_10k/slice-12                 70647             16809 ns/op
BenchmarkSearch/sparse_1k/slice-12                741074              1594 ns/op
BenchmarkSearch/sparse_10k/slice-12                58040             23650 ns/op
BenchmarkSearch/dense_1k/RBTree-12               9299102               127 ns/op
BenchmarkSearch/dense_10k/RBTree-12              5234247               225 ns/op
BenchmarkSearch/sparse_1k/RBTree-12              9728966               122 ns/op
BenchmarkSearch/sparse_10k/RBTree-12             5189743               228 ns/op
BenchmarkSearch/dense_1k/LLRBTree-12             8647956               135 ns/op
BenchmarkSearch/dense_10k/LLRBTree-12            5013253               237 ns/op
BenchmarkSearch/sparse_1k/LLRBTree-12            8898558               131 ns/op
BenchmarkSearch/sparse_10k/LLRBTree-12           4836254               245 ns/op
BenchmarkSearch/dense_1k/Skiplist-12             2830119               419 ns/op
BenchmarkSearch/dense_10k/Skiplist-12             255639              4659 ns/op
BenchmarkSearch/sparse_1k/Skiplist-12            3138306               372 ns/op
BenchmarkSearch/sparse_10k/Skiplist-12            432619              2771 ns/op
BenchmarkSearch/dense_1k/BPTree-12               8334624               139 ns/op
BenchmarkSearch/dense_10k/BPTree-12              4567568               254 ns/op
BenchmarkSearch/sparse_1k/BPTree-12              8727386               135 ns/op
BenchmarkSearch/sparse_10k/BPTree-12             4797978               245 ns/op
BenchmarkInsert/dense_1k/slice-12                 660908              1716 ns/op
BenchmarkInsert/dense_10k/slice-12                 66852             17291 ns/op
BenchmarkInsert/sparse_1k/slice-12                593624              3067 ns/op
BenchmarkInsert/sparse_10k/slice-12                66772             30527 ns/op
BenchmarkInsert/dense_1k/RBTree-12               9118976               129 ns/op
BenchmarkInsert/dense_10k/RBTree-12              5266610               223 ns/op
BenchmarkInsert/sparse_1k/RBTree-12              8338718               145 ns/op
BenchmarkInsert/sparse_10k/RBTree-12             4840795               252 ns/op
BenchmarkInsert/dense_1k/LLRBTree-12             5937380               196 ns/op
BenchmarkInsert/dense_10k/LLRBTree-12            3746398               316 ns/op
BenchmarkInsert/sparse_1k/LLRBTree-12            5392424               224 ns/op
BenchmarkInsert/sparse_10k/LLRBTree-12           3405024               355 ns/op
BenchmarkInsert/dense_1k/Skiplist-12             2532507               472 ns/op
BenchmarkInsert/dense_10k/Skiplist-12             503961              4435 ns/op
BenchmarkInsert/sparse_1k/Skiplist-12            1982985               688 ns/op
BenchmarkInsert/sparse_10k/Skiplist-12            441786              9035 ns/op
BenchmarkInsert/dense_1k/BPTree-12               7605034               151 ns/op
BenchmarkInsert/dense_10k/BPTree-12              4710601               252 ns/op
BenchmarkInsert/sparse_1k/BPTree-12              7248897               164 ns/op
BenchmarkInsert/sparse_10k/BPTree-12             4309918               284 ns/op
BenchmarkDelete/dense_1k/slice-12                 438283              9638 ns/op
BenchmarkDelete/dense_10k/slice-12                 70291             26461 ns/op
BenchmarkDelete/sparse_1k/slice-12                293746             13428 ns/op
BenchmarkDelete/sparse_10k/slice-12                49143             30881 ns/op
BenchmarkDelete/dense_1k/RBTree-12               1631854               810 ns/op
BenchmarkDelete/dense_10k/RBTree-12              1234207              1012 ns/op
BenchmarkDelete/sparse_1k/RBTree-12              1538397               842 ns/op
BenchmarkDelete/sparse_10k/RBTree-12             1210663              1043 ns/op
BenchmarkDelete/dense_1k/LLRBTree-12             1311289              1093 ns/op
BenchmarkDelete/dense_10k/LLRBTree-12            1000000              1425 ns/op
BenchmarkDelete/sparse_1k/LLRBTree-12            1000000              1154 ns/op
BenchmarkDelete/sparse_10k/LLRBTree-12           1000000              1493 ns/op
BenchmarkDelete/dense_1k/Skiplist-12             1000000              3067 ns/op
BenchmarkDelete/dense_10k/Skiplist-12             358872             13256 ns/op
BenchmarkDelete/sparse_1k/Skiplist-12            1000000              3542 ns/op
BenchmarkDelete/sparse_10k/Skiplist-12            322694             15015 ns/op
BenchmarkDelete/dense_1k/BPTree-12               1817608               768 ns/op
BenchmarkDelete/dense_10k/BPTree-12              1250306               990 ns/op
BenchmarkDelete/sparse_1k/BPTree-12              1607680               822 ns/op
BenchmarkDelete/sparse_10k/BPTree-12             1241611              1000 ns/op
PASS
ok      github.com/atriw/lib/golib/adt  147.801s
```