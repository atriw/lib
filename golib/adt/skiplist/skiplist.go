package skiplist

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/atriw/lib/golib/adt"
)

type node struct {
	key     adt.Key
	value   interface{}
	forward []*node
}

func (n *node) advance(level int, target adt.Key) *node {
	for next := n.forward[level]; next != nil && next.key.Less(target); {
		n = next
		next = next.forward[level]
	}
	return n
}

type nodeList []*node

func (nl nodeList) next() *node {
	return nl[0].forward[0]
}

func (nl nodeList) assertNext(key adt.Key) bool {
	next := nl.next()
	return next != nil && next.key.Equal(key)
}

type headerKey struct{}

func (hk *headerKey) Less(interface{}) bool  { return true }
func (hk *headerKey) Equal(interface{}) bool { return false }

// Skiplist is an instructive skiplist implementation without optimization or concurrency safety
type Skiplist struct {
	length   int
	maxLevel int
	header   *node
	level    int
}

const defaultLevel = 5

// Option is Skiplist initialization options
type Option func(*Skiplist)

// WithMaxLevel sets the max level of Skiplist
func WithMaxLevel(l int) Option {
	return func(sl *Skiplist) {
		sl.maxLevel = l
	}
}

// New returns an empty Skiplist
func New(opts ...Option) *Skiplist {
	sl := &Skiplist{maxLevel: defaultLevel, header: &node{key: &headerKey{}}}
	for _, o := range opts {
		o(sl)
	}
	for i := 0; i < sl.maxLevel; i++ {
		sl.header.forward = append(sl.header.forward, nil)
	}

	return sl
}

func (sl *Skiplist) prevNodes(key adt.Key) nodeList {
	prev := make(nodeList, sl.level+1)
	node := sl.header
	for i := sl.level; i >= 0; i-- {
		node = node.advance(i, key)
		prev[i] = node
	}
	return prev
}

// Search returns the value of key if exists, else nil
func (sl *Skiplist) Search(key adt.Key) interface{} {
	prev := sl.prevNodes(key)
	if prev.assertNext(key) {
		return prev.next().value
	}
	return nil
}

func (sl *Skiplist) randLevel() int {
	return rand.Intn(sl.maxLevel)
}

// Insert inserts key, value into Skiplist
func (sl *Skiplist) Insert(key adt.Key, value interface{}) {
	prev := sl.prevNodes(key)
	if prev.assertNext(key) {
		return
	}
	newLevel := sl.randLevel()
	if newLevel > sl.level {
		sl.level++
		newLevel = sl.level
		prev = append(prev, sl.header)
	}
	newNode := &node{key: key, value: value, forward: make([]*node, newLevel+1)}
	for i := newLevel; i >= 0; i-- {
		node := prev[i]
		newNode.forward[i] = node.forward[i]
		node.forward[i] = newNode
	}
	sl.length++
}

// Remove removes and returns value of key
func (sl *Skiplist) Delete(key adt.Key) interface{} {
	prev := sl.prevNodes(key)
	if !prev.assertNext(key) {
		return nil
	}
	sl.length--
	node := prev.next()
	for i, n := range node.forward {
		prev[i].forward[i] = n
	}
	return node.value
}

// Length returns total number of elements
func (sl *Skiplist) Length() int {
	return sl.length
}

func (sl *Skiplist) String() string {
	var sb strings.Builder
	zeroIndex := make(map[*node]int)
	indexLength := make(map[int]int)
	var buf []string
	for i := 0; i <= sl.level; i++ {
		sb.WriteString(fmt.Sprintf("%v: header->", i))
		j := 0
		for node := sl.header.forward[i]; node != nil; node = node.forward[i] {
			s := fmt.Sprintf("[%v:%v]", node.key, node.value)
			if i == 0 {
				zeroIndex[node] = j
				indexLength[j] = len(s)
			} else {
				for ; j < zeroIndex[node]; j++ {
					sb.WriteString(strings.Repeat("-", indexLength[j]+2))
				}
			}
			sb.WriteString(s)
			sb.WriteString("--")
			j++
		}
		for ; j < sl.length; j++ {
			sb.WriteString(strings.Repeat("-", indexLength[j]+2))
		}
		sb.WriteString("<nil>")
		buf = append(buf, sb.String())
		sb.Reset()
	}
	for i := len(buf) - 1; i >= 0; i-- {
		sb.WriteString(buf[i])
		if i != 0 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
