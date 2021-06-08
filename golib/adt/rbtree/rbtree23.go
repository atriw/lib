package rbtree

import (
	"fmt"

	"github.com/atriw/lib/golib/adt"
)

type node23 struct {
	Key   adt.Key
	Value interface{}

	left  *node23
	right *node23
	color color
	n     int
}

func (n *node23) size() int {
	if n == nil {
		return 0
	}
	return n.n
}

func (n *node23) isRed() bool {
	if n == nil {
		return false
	}
	return n.color == colorRed
}

func (n *node23) String() string {
	if n == nil {
		return "<nil>"
	}
	color := "r"
	if n.color == colorBlack {
		color = "b"
	}
	return fmt.Sprintf("[%v:%v:%v]", n.Key, n.Value, color)
}

func (n *node23) Left() adt.TreeNode {
	return n.left
}

func (n *node23) Right() adt.TreeNode {
	return n.right
}

func (n *node23) External() bool {
	return n == nil
}

type RBTree23 struct {
	root *node23
}

func New23() *RBTree23 {
	return &RBTree23{root: nil}
}

func (t *RBTree23) Search(key adt.Key) interface{} {
	return t.search(t.root, key)
}

func (t *RBTree23) search(n *node23, key adt.Key) interface{} {
	if n == nil {
		return nil
	}
	if key.Equal(n.Key) {
		return n.Value
	}
	if key.Less(n.Key) {
		return t.search(n.left, key)
	}
	return t.search(n.right, key)
}

func (t *RBTree23) Insert(key adt.Key, value interface{}) {
	t.root = t.insert(t.root, key, value)
	t.root.color = colorBlack
}

func (t *RBTree23) insert(n *node23, key adt.Key, value interface{}) *node23 {
	if n == nil {
		return &node23{Key: key, Value: value, color: colorRed, n: 1}
	}
	if key.Equal(n.Key) {
		n.Value = value
	} else if key.Less(n.Key) {
		n.left = t.insert(n.left, key, value)
	} else {
		n.right = t.insert(n.right, key, value)
	}

	if !n.left.isRed() && n.right.isRed() {
		n = t.leftRotate(n)
	}
	if n.left.isRed() && n.left.left.isRed() {
		n = t.rightRotate(n)
	}
	if n.left.isRed() && n.right.isRed() {
		t.flipColor(n)
	}

	n.n = n.left.size() + n.right.size() + 1
	return n
}

func (t *RBTree23) leftRotate(n *node23) *node23 {
	r := n.right
	n.right = r.left
	r.left = n
	n.color, r.color = r.color, n.color
	r.n = n.n
	n.n = n.left.size() + n.right.size() + 1
	return r
}

func (t *RBTree23) rightRotate(n *node23) *node23 {
	l := n.left
	n.left = l.right
	l.right = n
	n.color, l.color = l.color, n.color
	l.n = n.n
	n.n = n.left.size() + n.right.size() + 1
	return l
}

func (t *RBTree23) flipColor(n *node23) {
	n.left.color = complement(n.left.color)
	n.right.color = complement(n.right.color)
	n.color = complement(n.color)
}

func (t *RBTree23) Delete(key adt.Key) interface{} {
	if t.root == nil {
		return nil
	}
	if !t.root.left.isRed() && t.root.right.isRed() {
		t.root.color = colorRed
	}
	deleted := &node23{}
	t.root = t.delete(t.root, key, deleted)
	if t.root != nil {
		t.root.color = colorBlack
	}
	return deleted.Value
}

func (t *RBTree23) delete(n *node23, key adt.Key, deleted *node23) *node23 {
	if n == nil {
		return nil
	}
	if key.Less(n.Key) {
		if !n.left.isRed() && n.left != nil && !n.left.left.isRed() {
			n = t.moveRedLeft(n)
		}
		n.left = t.delete(n.left, key, deleted)
	} else {
		if n.left.isRed() {
			n = t.rightRotate(n)
		}
		if key.Equal(n.Key) && n.right == nil {
			*deleted = *n
			return nil
		}
		if !n.right.isRed() && n.right != nil && !n.right.left.isRed() {
			n = t.moveRedRight(n)
		}
		if key.Equal(n.Key) {
			min := &node23{}
			n.right = t.deleteMin(n.right, min)
			deleted.Value = n.Value
			n.Key = min.Key
			n.Value = min.Value
		} else {
			n.right = t.delete(n.right, key, deleted)
		}
	}
	return t.balance(n)
}

func (t *RBTree23) DeleteMin() (adt.Key, interface{}) {
	if t.root == nil {
		return nil, nil
	}
	if !t.root.left.isRed() {
		// Invariant: current node is not 2-node.
		t.root.color = colorRed
	}
	min := &node23{}
	t.root = t.deleteMin(t.root, min)
	if t.root != nil {
		t.root.color = colorBlack
	}
	return min.Key, min.Value
}

func (t *RBTree23) deleteMin(n *node23, min *node23) *node23 {
	if n.left == nil {
		*min = *n
		return nil
	}
	if !n.left.isRed() && !n.left.left.isRed() {
		n = t.moveRedLeft(n)
	}
	n.left = t.deleteMin(n.left, min)
	return t.balance(n)
}

func (t *RBTree23) DeleteMax() (adt.Key, interface{}) {
	if t.root == nil {
		return nil, nil
	}
	if !t.root.left.isRed() {
		t.root.color = colorRed
	}
	max := &node23{}
	t.root = t.deleteMax(t.root, max)
	if t.root != nil {
		t.root.color = colorBlack
	}
	return max.Key, max.Value
}

func (t *RBTree23) deleteMax(n *node23, max *node23) *node23 {
	if n.left.isRed() {
		// Make 3-node right-leaned.
		n = t.rightRotate(n)
	}
	if n.right == nil {
		*max = *n
		return nil
	}
	if !n.right.isRed() && n.right.left.isRed() {
		n = t.moveRedRight(n)
	}
	n.right = t.deleteMax(n.right, max)
	return t.balance(n)
}

func (t *RBTree23) moveRedLeft(n *node23) *node23 {
	//     n(r)                n(b)
	//    /  \                /  \
	//   x(b) y(b)  ---->    x(r) y(r)  2-node -> 4-node
	//  /                   /
	// z(b)                z(b)
	t.flipColor(n)
	if n.right.left.isRed() {
		//     n(r)            w(r)
		//    /  \            /  \
		//   x(b) y(b)  ---> n(b) y(b)  2-node + 3/4-node -> 3-node + 2/3-node
		//  /    /          /
		// z(b) w(r)       x(r)
		//                /
		//               z(b)
		n.right = t.rightRotate(n.right)
		n = t.leftRotate(n)
		t.flipColor(n)
	}
	return n
}

func (t *RBTree23) moveRedRight(n *node23) *node23 {
	//     n(r)                n(b)
	//    /  \                /  \
	//   y(b) x(b)  ---->    y(r) x(r)  2-node -> 4-node
	//		 /                   /
	//      z(b)                z(b)
	t.flipColor(n)
	if n.left.left.isRed() {
		//     n(r)                  y(r)
		//    /  \                  /  \
		//   y(b) x(b)  ---->      w(b) n(b)    2-node + 3/4-node -> 3-node(right-leaned) + 2/3-node
		//  /    /                       \
		// w(r) z(b)                     x(r)
		//                              /
		//                             z(b)
		n = t.rightRotate(n)
		t.flipColor(n)
	}
	return n
}

func (t *RBTree23) balance(n *node23) *node23 {
	if n == nil {
		return nil
	}

	if n.right.isRed() && !n.left.isRed() {
		// Make right-leaned red link left-leaned.
		n = t.leftRotate(n)
	}
	if n.left.isRed() && n.left.left.isRed() {
		// Make doulbe red link a 4-node.
		n = t.rightRotate(n)
	}
	if n.left.isRed() && n.right.isRed() {
		// Split 4-node.
		t.flipColor(n)
	}
	n.n = n.left.size() + n.right.size() + 1
	return n
}

func (t *RBTree23) Length() int {
	return t.root.size()
}

func (t *RBTree23) String() string {
	return adt.PrintTree(t.root)
}
