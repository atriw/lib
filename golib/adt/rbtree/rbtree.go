package rbtree

import (
	"fmt"

	"github.com/atriw/lib/golib/adt"
)

type color int

const (
	colorRed color = iota
	colorBlack
	colorDoubleBlack
)

type node struct {
	Key   adt.Key
	Value interface{}

	parent *node
	left   *node
	right  *node
	color  color
}

func newInternalNode(key adt.Key, value interface{}) *node {
	n := &node{Key: key, Value: value, color: colorRed}
	e1 := newExternalNode(n)
	e2 := newExternalNode(n)
	n.left = e1
	n.right = e2
	return n
}

func newExternalNode(parent *node) *node {
	return &node{parent: parent, color: colorBlack}
}

func (n *node) isLeft() bool {
	return n.parent != nil && n.parent.left == n
}

func (n *node) isRight() bool {
	return n.parent != nil && n.parent.right == n
}

func (n *node) isExternal() bool {
	return n.left == nil && n.right == nil
}

// isFull can only be called on internal node.
func (n *node) isFull() bool {
	return !n.left.isExternal() && !n.right.isExternal()
}

// child can only be called on internal node.
func (n *node) child() *node {
	if !n.left.isExternal() {
		return n.left
	}
	return n.right
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

// successor can only be called on internal node.
func (n *node) successor() *node {
	succ := n.right
	for !succ.isExternal() && !succ.left.isExternal() {
		succ = succ.left
	}
	return succ
}

func (n *node) isRed() bool {
	return !n.isExternal() && n.color == colorRed
}

func (n *node) isBlack() bool {
	return n.isExternal() || n.color == colorBlack
}

func (n *node) String() string {
	if n.isExternal() {
		return "[Ext]"
	}
	color := "r"
	if n.color == colorBlack {
		color = "b"
	}
	height, _ := n.blackHeight()
	return fmt.Sprintf("[%v:%v:%v:%v]", n.Key, n.Value, color, height)
}

func (n *node) Left() adt.TreeNode {
	return n.left
}

func (n *node) Right() adt.TreeNode {
	return n.right
}

func (n *node) External() bool {
	return n.isExternal()
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
	return &RBTree{root: newExternalNode(nil)}
}

func (t *RBTree) Length() int {
	return t.length
}

func (t *RBTree) Search(key adt.Key) interface{} {
	p, dir := t.search(key)
	if p.isExternal() {
		return nil
	}
	n := p.dir(dir)
	if n.isExternal() {
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

// search searches for the given key.
// If
// - the root is external, returns root, self.
// - finds the exact node, returns node, self.
// - no exact match found, stops at external node, returns the parent of the external node, and direction refering to it.
func (t *RBTree) search(key adt.Key) (p *node, dir direction) {
	n := t.root
	p = n
	for !n.isExternal() {
		p = n
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
	return p, dir
}

func (t *RBTree) Insert(key adt.Key, value interface{}) {
	p, dir := t.search(key)
	// Find existing key.
	if !p.isExternal() && dir == self {
		p.Value = value
		return
	}
	t.length++
	// Find root that is external
	if p.isExternal() {
		t.root = newInternalNode(key, value)
		t.root.color = colorBlack
		return
	}

	n := newInternalNode(key, value)
	n.parent = p
	p.setDir(dir, n)
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

// rechild sets n's parent's child to c, but not sets c'parent to n's parent.
func (t *RBTree) rechild(n, c *node) {
	if n.parent == nil {
		t.root = c
	} else {
		if n.isLeft() {
			n.parent.left = c
		} else {
			n.parent.right = c
		}
	}
}

// reparent sets c'parent to p, p's parent to c and rechild p, c.
//    p                 c
//    |     ------->    |
//    c                 p
func (t *RBTree) reparent(p, c *node) {
	c.parent = p.parent
	t.rechild(p, c)
	p.parent = c
}

//      l                      r
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

func (t *RBTree) Delete(key adt.Key) interface{} {
	n, dir := t.search(key)
	if dir != self {
		return nil
	}
	if n.isExternal() {
		return nil
	}
	t.length--
	v := n.Value
	t.delete(n)
	return v
}

func (t *RBTree) delete(n *node) {
	if n.isFull() {
		succ := n.successor()
		n.Key, n.Value = succ.Key, succ.Value
		n = succ
	}
	child := n.child()
	t.reparent(n, child)
	if n.isRed() {
		return
	}
	n = child
	if n.isRed() {
		n.color = colorBlack
		return
	}
	// n.color = colorDoubleBlack
	for n != t.root {
		brother := n.brother()
		var opposite, same direction
		var oppAdj, sameAdj func(*node)
		if n.isRight() {
			opposite, same = left, right
			oppAdj, sameAdj = t.leftRotate, t.rightRotate
		} else {
			opposite, same = right, left
			oppAdj, sameAdj = t.rightRotate, t.leftRotate
		}
		if brother.isBlack() {
			if brother.dir(opposite).isRed() {
				n.color = colorBlack
				brother.dir(opposite).color = colorBlack
				brother.color = n.parent.color
				n.parent.color = colorBlack
				sameAdj(n.parent)
				return
			}
			if brother.dir(same).isRed() {
				n.color = colorBlack
				brother.dir(same).color = n.parent.color
				n.parent.color = colorBlack
				oppAdj(brother)
				sameAdj(n.parent)
				return
			}
			if n.parent.isRed() {
				n.color = colorBlack
				brother.color = colorRed
				n.parent.color = colorBlack
				return
			}
			n.color = colorBlack
			brother.color = colorRed
			// n.parent.color = colorDoubleBlack
			n = n.parent
			continue
		}
		brother.color = n.parent.color
		n.parent.color = colorRed
		sameAdj(n.parent)
	}
}

func (t *RBTree) String() string {
	return adt.PrintTree(t.root)
}

func (t *RBTree) Validate() bool {
	return t.root.propertyRedHasNoRedChildren() && t.root.propertyBlackHeightEqual()
}

func (n *node) propertyRedHasNoRedChildren() bool {
	if n.isExternal() {
		return true
	}
	if n.isRed() {
		if n.left.isRed() || n.right.isRed() {
			return false
		}
	}
	return n.left.propertyRedHasNoRedChildren() && n.right.propertyRedHasNoRedChildren()
}

func (n *node) propertyBlackHeightEqual() bool {
	_, t := n.blackHeight()
	return t
}

func (n *node) blackHeight() (int, bool) {
	if n.isExternal() {
		return 0, true
	}
	lbh, ok := n.left.blackHeight()
	if !ok {
		return 0, false
	}
	rbh, ok := n.right.blackHeight()
	if !ok {
		return 0, false
	}
	if lbh != rbh {
		return 0, false
	}
	if n.isBlack() {
		lbh += 1
	}
	return lbh, true
}
