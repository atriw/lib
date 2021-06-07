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
	return r
}

func (t *RBTree23) rightRotate(n *node23) *node23 {
	l := n.left
	n.left = l.right
	l.right = n
	n.color, l.color = l.color, n.color
	return l
}

func (t *RBTree23) flipColor(n *node23) {
	n.left.color = colorBlack
	n.right.color = colorBlack
	n.color = colorRed
}

func (t *RBTree23) Delete(key adt.Key) interface{} {
	return nil
}

func (t *RBTree23) Length() int {
	return t.root.size()
}

func (t *RBTree23) String() string {
	return adt.PrintTree(t.root)
}
