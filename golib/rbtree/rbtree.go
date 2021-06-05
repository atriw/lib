package rbtree

import (
	"fmt"
	"strings"
)

type KeyType interface {
	Less(interface{}) bool
	Equal(interface{}) bool
}

type color int

const (
	colorRed color = iota
	colorBlack
)

type node struct {
	Key   KeyType
	Value interface{}

	parent *node
	left   *node
	right  *node
	color  color
}

func (n *node) isLeft() bool {
	return n.parent != nil && n.parent.left == n
}

func (n *node) isRight() bool {
	return n.parent != nil && n.parent.right == n
}

func (n *node) brother() *node {
	if n.parent == nil {
		return nil
	}
	if n.isLeft() {
		return n.parent.right
	}
	return n.parent.left
}

func (n *node) isRed() bool {
	return n != nil && n.color == colorRed
}

func (n *node) isBlack() bool {
	return n == nil || n.color == colorBlack
}

func (n *node) String() string {
	if n == nil {
		return "<nil>"
	}
	color := "r"
	if n.color == colorBlack {
		color = "b"
	}
	return fmt.Sprintf("[%v:%v:%v]", n.Key, n.Value, color)
}

func (n *node) print(sb *strings.Builder, prefix, childPrefix string) {
	sb.WriteString(prefix)
	sb.WriteString(n.String())
	sb.WriteString("\n")
	if n == nil {
		return
	}
	n.left.print(sb, childPrefix+"├── ", childPrefix+"│   ")
	n.right.print(sb, childPrefix+"└── ", childPrefix+"    ")
}

func (n *node) setDir(dir direction, c *node) {
	if dir == self {
		panic("wrong direction")
	}
	if dir == left {
		n.left = c
	} else {
		n.right = c
	}
}

func (n *node) dir(dir direction) *node {
	if dir == self {
		return n
	}
	if dir == left {
		return n.left
	}
	return n.right
}

type RBTree struct {
	length int
	root   *node
}

func New() *RBTree {
	return &RBTree{root: nil}
}

func (t *RBTree) Length() int {
	return t.length
}

func (t *RBTree) Search(key KeyType) interface{} {
	p, dir := t.search(key)
	if p == nil {
		return nil
	}
	n := p.dir(dir)
	if n == nil {
		return nil
	}
	return n.Value
}

type direction int

const (
	self direction = iota
	left
	right
)

func (t *RBTree) search(key KeyType) (parent *node, dir direction) {
	n := t.root
	for n != nil {
		parent = n
		if key.Equal(n.Key) {
			dir = self
			break
		}
		if key.Less(n.Key) {
			n = n.left
			dir = left
		} else {
			n = n.right
			dir = right
		}
	}
	return parent, dir
}

func (t *RBTree) Insert(key KeyType, value interface{}) {
	parent, dir := t.search(key)
	if parent != nil && dir == self {
		return
	}
	t.length++
	if parent == nil {
		t.root = &node{Key: key, Value: value, color: colorBlack}
		return
	}
	n := &node{Key: key, Value: value, parent: parent, color: colorRed}
	parent.setDir(dir, n)
	for n != t.root && n.parent.isRed() {
		uncle := n.parent.brother()
		// case 1:
		//      b              r <--- next iteration
		//     / \            / \
		//    p   r   ----> p(b) b
		//   / \            / \
		//  n   b          n   b
		if uncle.isRed() {
			n.parent.color = colorBlack
			uncle.color = colorBlack
			n.parent.parent.color = colorRed
			n = n.parent.parent
			continue
		}

		// case 2:
		//      b                      b
		//     / \    leftRotate      / \
		//    p   b   ---------->    n   b
		//   / \                    / \
		//  b   n    case 3 ---->  p   b
		//                        /
		//                       b
		//
		//      b                      b
		//     / \    rightRotate     / \
		//    b   p   ---------->    b   n
		//       / \                    / \
		//      n   b                  b   p  <---- case 3
		//                                  \
		//                                   b
		if n.parent.isLeft() && n.isRight() {
			n = n.parent
			t.leftRotate(n)
		} else if n.parent.isRight() && n.isLeft() {
			n = n.parent
			t.rightRotate(n)
		}

		// case 3:
		//      b             p(b)
		//     / \            / \
		//    p   b  ---->   n   r
		//   / \                / \
		//  n   b              b   b
		//
		//      b               p(b)
		//     / \              / \
		//    b   p   ---->    r   n
		//       / \          / \
		//      b   n        b   b
		n.parent.color = colorBlack
		n.parent.parent.color = colorRed
		if n.parent.isLeft() {
			t.rightRotate(n.parent.parent)
		} else {
			t.leftRotate(n.parent.parent)
		}
	}
	t.root.color = colorBlack
}

// newChild         newParent
//    |     ------->    |
// newParent        newChild
func (t *RBTree) reparent(newChild, newParent *node) {
	newParent.parent = newChild.parent
	if newChild.parent == nil {
		t.root = newParent
	} else {
		if newChild.isLeft() {
			newChild.parent.left = newParent
		} else {
			newChild.parent.right = newParent
		}
	}
	newChild.parent = newParent
}

//		l                      r
//     / \     leftRotate     / \
//    A   r    --------->    l   C
//       / \   <---------   / \
//      B   C  rightRotate A   B
func (t *RBTree) leftRotate(l *node) {
	r := l.right
	l.right = r.left
	if r.left != nil {
		r.left.parent = l
	}
	t.reparent(l, r)
	r.left = l
}

func (t *RBTree) rightRotate(r *node) {
	l := r.left
	r.left = l.right
	if l.right != nil {
		l.right.parent = r
	}
	t.reparent(r, l)
	l.right = r
}

func (t *RBTree) Remove(key KeyType) interface{} {
	return nil
}

func (t *RBTree) String() string {
	var sb strings.Builder
	t.root.print(&sb, "", "")
	return sb.String()
}
