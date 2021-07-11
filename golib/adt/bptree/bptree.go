package bptree

import (
	"fmt"
	"strings"

	"github.com/atriw/lib/golib/adt"
)

const debug = false
const validate = false

type node struct {
	leaf     bool
	n        int
	keys     keys
	values   values
	children children
}

type keys []adt.Key

func (k keys) insert(key adt.Key, idx int, limit int) {
	for i := limit; i > idx; i-- {
		k[i] = k[i-1]
	}
	k[idx] = key
}

func (k keys) delete(idx int, limit int) {
	for i := idx; i < limit-1; i++ {
		k[i] = k[i+1]
	}
}

type values []interface{}

func (v values) insert(value interface{}, idx int, limit int) {
	for i := limit; i > idx; i-- {
		v[i] = v[i-1]
	}
	v[idx] = value
}

func (v values) delete(idx int, limit int) {
	for i := idx; i < limit-1; i++ {
		v[i] = v[i+1]
	}
}

type children []*node

func (c children) insert(child *node, idx int, limit int) {
	for i := limit; i > idx; i-- {
		c[i] = c[i-1]
	}
	c[idx] = child
}

func (c children) delete(idx int, limit int) {
	for i := idx; i < limit-1; i++ {
		c[i] = c[i+1]
	}
}

func newNode(leaf bool, order int) *node {
	return &node{
		leaf: leaf,
		// The last key of non-leaf node is always nil.
		keys:     make(keys, order),
		values:   make(values, order),
		children: make(children, order),
	}
}

func (n *node) size() int {
	if n.leaf {
		return n.n
	}
	return n.n + 1
}

func (n *node) isFull() bool {
	return n.size() == len(n.keys)
}

func (n *node) leafInsert(key adt.Key, value interface{}) {
	idx, exact := find(n.keys, key, n.n)
	if exact {
		n.values[idx] = value
		return
	}
	n.keys.insert(key, idx, n.n)
	n.values.insert(value, idx, n.n)
	n.n++
}

func (n *node) internalInsert(key adt.Key, child *node) {
	idx, exact := find(n.keys, key, n.n)
	if exact {
		panic(fmt.Sprintf("duplicate internal insert: key %v, node:\n%v\nchild:\n%v\n", key, adt.PrintMultiWayTreeDepth(n, 1), adt.PrintMultiWayTreeDepth(child, 1)))
	}
	n.keys.insert(key, idx, n.n)
	n.children.insert(child, idx+1, n.n+1)
	n.n++
}

func (n *node) String() string {
	if n == nil {
		return "<nil>"
	}
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
	if n == nil {
		return nil
	}
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
	if validate {
		backup := adt.PrintMultiWayTree(t.root)
		defer func() {
			if !t.propertyHalfFull(t.root) {
				panic(fmt.Sprintf("insert %v, not half full, tree: \n%v\nbackup: \n%v\n", key, adt.PrintMultiWayTree(t.root), backup))
			}
			if !t.propertySameHeight() {
				panic(fmt.Sprintf("insert %v, not same height, tree: \n%v\nbackup: \n%v\n", key, adt.PrintMultiWayTree(t.root), backup))
			}
		}()
	}
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
	if debug {
		fmt.Println("insert", key, n.children[idx])
	}
	s, l := t.insert(n.children[idx], key, value)
	if s == nil {
		return nil, nil
	}
	return t.insertInternal(n, l, s)
}

func (t *BPTree) insertLeaf(n *node, key adt.Key, value interface{}) (split *node, lastKey adt.Key) {
	idx, exact := find(n.keys, key, n.n)
	if exact {
		n.values[idx] = value
		return
	}
	if !n.isFull() {
		n.keys.insert(key, idx, n.n)
		n.values.insert(value, idx, n.n)
		n.n++
		return nil, nil
	}
	split = newNode(true, t.order)
	keys, values := make(keys, n.n+1), make(values, n.n+1)
	copy(keys, n.keys)
	copy(values, n.values)
	keys.insert(key, idx, n.n)
	values.insert(value, idx, n.n)
	half := len(keys) / 2
	copy(n.keys, keys[:half])
	copy(n.values, values[:half])
	copy(split.keys, keys[half:])
	copy(split.values, values[half:])
	n.n, split.n = half, len(keys)-half
	return split, n.lastKey()
}

func (t *BPTree) insertInternal(n *node, l adt.Key, s *node) (split *node, lastKey adt.Key) {
	idx, exact := find(n.keys, l, n.n)
	if exact {
		panic("duplicate internal insert")
	}
	if !n.isFull() {
		n.keys.insert(l, idx, n.n)
		n.children.insert(s, idx+1, n.n+1)
		n.n++
		return nil, nil
	}
	if debug {
		fmt.Println("split", l, adt.PrintMultiWayTreeDepth(n, 1), adt.PrintMultiWayTreeDepth(s, 1))
	}
	split = newNode(false, t.order)
	keys, children := make(keys, n.n+1), make(children, n.n+2)
	copy(keys, n.keys)
	copy(children, n.children)
	keys.insert(l, idx, n.n)
	children.insert(s, idx+1, n.n+1)
	half := len(keys) / 2
	lastKey = keys[half]
	copy(n.keys, keys[:half])
	copy(n.children, children[:half+1])
	copy(split.keys, keys[half+1:])
	copy(split.children, children[half+1:])
	n.n, split.n = half, len(keys)-half-1
	if debug {
		fmt.Println("split result", adt.PrintMultiWayTreeDepth(n, 1), adt.PrintMultiWayTreeDepth(split, 1))
	}
	return split, lastKey
}

func (t *BPTree) Delete(key adt.Key) interface{} {
	if validate {
		backup := adt.PrintMultiWayTree(t.root)
		defer func() {
			if !t.propertyHalfFull(t.root) {
				panic(fmt.Sprintf("delete %v, not half full, tree: \n%v\nbackup: \n%v\n", key, adt.PrintMultiWayTree(t.root), backup))
			}
			if !t.propertySameHeight() {
				panic(fmt.Sprintf("delete %v, not same height, tree: \n%v\nbackup: \n%v\n", key, adt.PrintMultiWayTree(t.root), backup))
			}
		}()
	}
	if t.root == nil {
		return nil
	}
	var deleted interface{}
	t.delete(t.root, key, &deleted)
	if t.root.n == 0 {
		t.root = t.root.children[0]
	}
	return deleted
}

func (t *BPTree) delete(n *node, key adt.Key, deleted *interface{}) (underflow bool) {
	if n == nil {
		return
	}
	if n.leaf {
		n.leafDelete(key, deleted)
		return n.underflow()
	}
	// Finds the first key which is greater or equal than the needed key.
	idx, _ := find(n.keys, key, n.n)
	// Recurs into its left child.
	if debug {
		fmt.Println("delete", key, adt.PrintMultiWayTree(n.children[idx]))
	}
	underflow = t.delete(n.children[idx], key, deleted)
	if !underflow {
		return
	}
	if debug {
		fmt.Println("underflow", adt.PrintMultiWayTreeDepth(n.children[idx], 1))
	}
	if idx == n.n {
		return t.underflow(n, idx-1, true)
	}
	return t.underflow(n, idx, false)
}

func (t *BPTree) underflow(n *node, idx int, last bool) (underflow bool) {
	// Merge right sibling into the underflow node.
	// And delete the sibling and the key whose right child is the sibling.
	left := n.children[idx]
	right := n.children[idx+1]
	if right == nil {
		panic(adt.PrintMultiWayTreeDepth(n, 1))
	}
	if (!last && t.shouldMerge(right)) || (last && t.shouldMerge(left)) {
		n.merge(n.keys[idx], left, right)
		n.internalDelete(idx)
		return n.underflow()
	}
	from, to := right, left
	if last {
		from, to = to, from
	}
	n.keys[idx] = n.transfer(n.keys[idx], from, to, last)
	if debug {
		fmt.Println("transfer result", adt.PrintMultiWayTreeDepth(n, 1), adt.PrintMultiWayTreeDepth(from, 1), adt.PrintMultiWayTreeDepth(to, 1))
	}
	return false
}

func (t *BPTree) half() int {
	return (t.order + 1) / 2
}

func (t *BPTree) shouldMerge(n *node) bool {
	return n.size() <= t.order-t.half()+1
}

// merge merges right into left.
func (n *node) merge(mid adt.Key, left, right *node) {
	defer func() {
		if debug {
			fmt.Println("merge result", adt.PrintMultiWayTreeDepth(left, 1))
		}
	}()
	if debug {
		fmt.Println("merge", adt.PrintMultiWayTreeDepth(left, 1), adt.PrintMultiWayTreeDepth(right, 1))
	}
	// Internal merge.
	if !left.leaf {
		// Append mid key from parent to left keys.
		left.internalAppendRightKey(mid)
		// Move all keys from right to left.
		left.internalAppendRight(right.n, right.keys, right.children)
		// Move the last children
		left.children[left.n] = right.children[right.n]
		return
	}
	left.leafAppendRight(right.n, right.keys, right.values)
}

// Delete the idxth key and its right child in an internal node.
func (n *node) internalDelete(idx int) {
	n.keys.delete(idx, n.n)
	n.children.delete(idx+1, n.n+1)
	n.n--
}

// Find and delete a key in a leaf node.
func (n *node) leafDelete(key adt.Key, deleted *interface{}) {
	for i := 0; i < n.n; i++ {
		if key.Equal(n.keys[i]) {
			*deleted = n.values[i]
			n.keys.delete(i, n.n)
			n.values.delete(i, n.n)
			n.n--
			return
		}
		// TODO: early return
	}
}

func (n *node) underflow() bool {
	return n.numToFillUnderflow() > 0
}

func (n *node) numToFillUnderflow() int {
	return (len(n.keys)+1)/2 - n.size()
}

func (n *node) transfer(mid adt.Key, from, to *node, leftToRight bool) adt.Key {
	if debug {
		fmt.Println("transfer", mid, adt.PrintMultiWayTreeDepth(n, 1), adt.PrintMultiWayTreeDepth(from, 1), adt.PrintMultiWayTreeDepth(to, 1))
	}
	total := to.numToFillUnderflow()
	if leftToRight {
		return transferLeftRight(mid, from, to, total)
	}
	return transferRightLeft(mid, from, to, total)
}

func (n *node) leafAppendRight(limit int, keys []adt.Key, values []interface{}) {
	for i := 0; i < limit; i++ {
		n.keys[n.n+i] = keys[i]
		n.values[n.n+i] = values[i]
	}
	n.n += limit
}

func (n *node) leafAppendLeft(limit int, keys []adt.Key, values []interface{}) {
	fill := limit
	for i := fill + n.n - 1; i >= fill; i-- {
		n.keys[i] = n.keys[i-fill]
		n.values[i] = n.values[i-fill]
	}
	for i := 0; i < fill; i++ {
		n.keys[i] = keys[i]
		n.values[i] = values[i]
	}
	n.n += fill
}

func (n *node) internalAppendRight(limit int, keys []adt.Key, children []*node) {
	for i := 0; i < limit; i++ {
		n.keys[n.n+i] = keys[i]
		n.children[n.n+i] = children[i]
	}
	n.n += limit
}

func (n *node) internalAppendLeft(limit int, keys []adt.Key, children []*node) {
	fill := limit
	for i := fill + n.n - 1; i >= fill; i-- {
		n.keys[i] = n.keys[i-fill]
		n.children[i] = n.children[i-fill]
	}
	for i := 0; i < fill; i++ {
		n.keys[i] = keys[i]
		n.children[i] = children[i]
	}
	n.n += fill
}

func (n *node) internalAppendRightKey(key adt.Key) {
	n.keys[n.n] = key
	n.n++
}

func (n *node) internalAppendLeftKey(key adt.Key) {
	for i := n.n; i >= 1; i-- {
		n.keys[i] = n.keys[i-1]
	}
	n.keys[0] = key
	n.n++
}

func (n *node) internalPopRightKey() adt.Key {
	n.n--
	return n.keys[n.n]
}

func (n *node) internalPopLeftKey() adt.Key {
	mid := n.keys[0]
	for i := 0; i < n.n-1; i++ {
		n.keys[i] = n.keys[i+1]
	}
	n.n--
	return mid
}

func (n *node) leafPopLeft(num int) (keys []adt.Key, values []interface{}) {
	if n.n < num {
		panic("pop")
	}
	for i := 0; i < num; i++ {
		keys = append(keys, n.keys[i])
		values = append(values, n.values[i])
	}
	for i := 0; i < n.n-num; i++ {
		n.keys[i] = n.keys[i+num]
		n.values[i] = n.values[i+num]
	}
	n.n -= num
	return keys, values
}

func (n *node) leafPopRight(num int) (keys []adt.Key, values []interface{}) {
	for i := n.n - num; i < n.n; i++ {
		keys = append(keys, n.keys[i])
		values = append(values, n.values[i])
	}
	n.n -= num
	return keys, values
}

func (n *node) internalPopLeft(num int) (keys []adt.Key, children []*node) {
	if n.n < num {
		panic("pop")
	}
	for i := 0; i < num; i++ {
		keys = append(keys, n.keys[i])
		children = append(children, n.children[i])
	}
	for i := 0; i < n.n-num; i++ {
		n.keys[i] = n.keys[i+num]
		n.children[i] = n.children[i+num]
	}
	n.children[n.n-num] = n.children[n.n]
	n.n -= num
	return keys, children
}

func (n *node) internalPopRight(num int) (keys []adt.Key, children []*node) {
	for i := n.n - num; i < n.n; i++ {
		keys = append(keys, n.keys[i])
		children = append(children, n.children[i+1])
	}
	n.n -= num
	return keys, children
}

func (n *node) lastKey() adt.Key {
	return n.keys[n.n-1]
}

// Transfer from right to left
func transferRightLeft(mid adt.Key, from, to *node, num int) adt.Key {
	if from.leaf {
		keys, values := from.leafPopLeft(num)
		to.leafAppendRight(num, keys, values)
		return to.lastKey()
	}
	to.internalAppendRightKey(mid)
	keys, children := from.internalPopLeft(num)
	to.internalAppendRight(num, keys, children)
	return to.internalPopRightKey()
}

func transferLeftRight(mid adt.Key, from, to *node, num int) adt.Key {
	if from.leaf {
		keys, values := from.leafPopRight(num)
		to.leafAppendLeft(num, keys, values)
		return from.lastKey()
	}
	to.internalAppendLeftKey(mid)
	keys, children := from.internalPopRight(num)
	to.internalAppendLeft(num, keys, children)
	return to.internalPopLeftKey()
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
	half := t.half()
	if n == nil {
		return true
	}
	if n.leaf {
		if n == t.root {
			return n.n >= 1
		}
		return n.n >= half
	}
	if n == t.root {
		if n.size() < 2 {
			return false
		}
	} else {
		if n.size() < half {
			return false
		}
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
