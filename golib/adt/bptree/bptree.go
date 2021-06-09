package bptree

import (
	"fmt"
	"strings"

	"github.com/atriw/lib/golib/adt"
)

type node struct {
	leaf     bool
	n        int
	keys     []adt.Key
	values   []interface{}
	children []*node
}

func newNode(leaf bool, order int) *node {
	var values []interface{}
	var children []*node
	if leaf {
		values = make([]interface{}, order)
	} else {
		children = make([]*node, order)
	}
	return &node{
		leaf: leaf,
		// The last key of non-leaf node is always nil.
		keys:     make([]adt.Key, order),
		values:   values,
		children: children,
	}
}

func (n *node) isFull() bool {
	if n.leaf {
		return n.n == len(n.keys)
	}
	return n.n == len(n.keys)-1
}

func (n *node) leafInsert(key adt.Key, value interface{}) {
	idx, exact := find(n.keys, key, n.n)
	if exact {
		n.values[idx] = value
		return
	}
	for i := n.n; i > idx; i-- {
		n.keys[i] = n.keys[i-1]
		n.values[i] = n.values[i-1]
	}
	n.keys[idx] = key
	n.values[idx] = value
	n.n++
}

func (n *node) internalInsert(key adt.Key, child *node) {
	idx, exact := find(n.keys, key, n.n)
	if exact {
		panic("duplicate internal insert")
	}
	for i := n.n; i > idx; i-- {
		n.keys[i] = n.keys[i-1]
		n.children[i+1] = n.children[i]
	}
	n.keys[idx] = key
	n.children[idx+1] = child
	n.n++
}

func (n *node) String() string {
	var sb strings.Builder
	sb.WriteString("[<")
	if n.leaf {
		sb.WriteString("l")
	} else {
		sb.WriteString("i")
	}
	sb.WriteString(fmt.Sprintf(":%v>:", n.n))
	for i := 0; i < n.n; i++ {
		if i != 0 {
			sb.WriteString(":")
		}
		if n.leaf {
			sb.WriteString("<")
		}
		sb.WriteString(fmt.Sprint(n.keys[i]))
		if n.leaf {
			sb.WriteString(fmt.Sprintf(":%v>", n.values[i]))
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (n *node) Iterator() adt.Iterator {
	if n.leaf {
		return nil
	}
	return &iterator{nodes: n.children, n: n.n + 1}
}

type iterator struct {
	nodes []*node
	idx   int
	n     int
}

func (i *iterator) HasNext() bool {
	return i.idx < i.n
}

func (i *iterator) Next() adt.MultiWayTreeNode {
	n := i.nodes[i.idx]
	i.idx++
	return n
}

type BPTree struct {
	root  *node
	order int
}

const defaultOrder = 128

func New(opts ...Option) *BPTree {
	t := &BPTree{root: nil, order: defaultOrder}
	for _, opt := range opts {
		opt(t)
	}
	if t.order <= 3 {
		// TODO: support order 2 and 3
		panic("order should be at least 4")
	}
	return t
}

type Option func(*BPTree)

func WithOrder(order int) Option {
	return func(t *BPTree) {
		t.order = order
	}
}

func (t *BPTree) Search(key adt.Key) interface{} {
	return t.search(t.root, key)
}

func (t *BPTree) search(n *node, key adt.Key) interface{} {
	if n == nil {
		return nil
	}
	idx, exact := find(n.keys, key, n.n)
	if n.leaf {
		if !exact {
			return nil
		}
		return n.values[idx]
	}
	return t.search(n.children[idx], key)
}

func (t *BPTree) Insert(key adt.Key, value interface{}) {
	if t.root == nil {
		t.root = newNode(true, t.order)
		t.root.leafInsert(key, value)
		return
	}
	split, lastKey := t.insert(t.root, key, value)
	if split == nil {
		return
	}
	newRoot := newNode(false, t.order)
	newRoot.children[0] = t.root
	newRoot.internalInsert(lastKey, split)
	t.root = newRoot
}

func (t *BPTree) insert(n *node, key adt.Key, value interface{}) (split *node, lastKey adt.Key) {
	if n == nil {
		return nil, nil
	}
	if n.leaf {
		return t.insertLeaf(n, key, value)
	}
	idx, _ := find(n.keys, key, n.n)
	s, l := t.insert(n.children[idx], key, value)
	if s == nil {
		return nil, nil
	}
	return t.insertInternal(n, l, s)
}

func (t *BPTree) insertLeaf(n *node, key adt.Key, value interface{}) (split *node, lastKey adt.Key) {
	if !n.isFull() {
		n.leafInsert(key, value)
		return nil, nil
	}
	split = newNode(true, t.order)
	half := n.n / 2
	for i := 0; i < half; i++ {
		split.keys[i] = n.keys[i+n.n-half]
		split.values[i] = n.values[i+n.n-half]
	}
	n.n -= half
	split.n = half
	if key.Less(n.keys[half]) {
		n.leafInsert(key, value)
	} else {
		split.leafInsert(key, value)
	}
	return split, n.keys[n.n-1]
}

func (t *BPTree) insertInternal(n *node, l adt.Key, s *node) (split *node, lastKey adt.Key) {
	if !n.isFull() {
		n.internalInsert(l, s)
		return nil, nil
	}
	split = newNode(false, t.order)
	half := n.n / 2
	for i := 0; i < half; i++ {
		split.keys[i] = n.keys[i+n.n-half]
		split.children[i] = n.children[i+n.n-half]
	}
	split.children[half] = n.children[n.n]
	n.n -= half
	split.n = half
	lastKey = n.keys[n.n-1]
	n.n--
	if l.Less(lastKey) {
		n.internalInsert(l, s)
	} else {
		split.internalInsert(l, s)
	}
	return split, lastKey
}

func (t *BPTree) Delete(key adt.Key) interface{} {
	return nil
}

func (t *BPTree) Length() int {
	return t.length(t.root)
}

func (t *BPTree) length(n *node) int {
	if n == nil {
		return 0
	}
	if n.leaf {
		return n.n
	}
	sum := 0
	for i := 0; i <= n.n; i++ {
		sum += t.length(n.children[i])
	}
	return sum
}

func (t *BPTree) String() string {
	return adt.PrintMultiWayTree(t.root)
}

func (t *BPTree) Validate() bool {
	return t.propertySameHeight() && t.propertyHalfFull(t.root)
}

func (t *BPTree) propertySameHeight() bool {
	_, ok := t.height(t.root)
	return ok
}

func (t *BPTree) height(n *node) (int, bool) {
	if n == nil || n.leaf {
		return 0, true
	}
	h := -1
	for i := 0; i <= n.n; i++ {
		ch, ok := t.height(n.children[i])
		if !ok {
			return 0, false
		}
		if h >= 0 && ch != h {
			return 0, false
		}
		h = ch
	}
	return h, true
}

func (t *BPTree) propertyHalfFull(n *node) bool {
	half := (t.order + 1) / 2
	if n == nil {
		return true
	}
	if n.leaf {
		if n == t.root {
			return n.n >= 1
		}
		return n.n >= half
	}
	if n == t.root && n.n+1 < 2 {
		return false
	}
	if n.n+1 < half {
		return false
	}
	for i := 0; i <= n.n; i++ {
		if !t.propertyHalfFull(n.children[i]) {
			return false
		}
	}
	return true
}

func find(keys []adt.Key, target adt.Key, limit int) (idx int, exact bool) {
	for ; idx < limit; idx++ {
		if keys[idx].Less(target) {
			continue
		}
		if keys[idx].Equal(target) {
			return idx, true
		}
		break
	}
	return idx, false
}
