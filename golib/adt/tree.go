package adt

import (
	"fmt"
	"strings"
)

type TreeNode interface {
	Left() TreeNode
	Right() TreeNode
	External() bool
}

func PrintTree(n TreeNode) string {
	var sb strings.Builder
	printTree(n, &sb, "", "")
	return sb.String()
}

func printTree(n TreeNode, sb *strings.Builder, prefix, childPrefix string) {
	sb.WriteString(prefix)
	if x, ok := n.(interface{ String() string }); ok {
		sb.WriteString(x.String())
	} else {
		sb.WriteString(fmt.Sprint(n))
	}
	sb.WriteString("\n")
	if n.External() {
		return
	}
	printTree(n.Right(), sb, childPrefix+"├── ", childPrefix+"│   ")
	printTree(n.Left(), sb, childPrefix+"└── ", childPrefix+"    ")
}

type MultiWayTreeNode interface {
	Iterator() Iterator
}

type Iterator interface {
	HasNext() bool
	Next() MultiWayTreeNode
}

func PrintMultiWayTree(n MultiWayTreeNode) string {
	var sb strings.Builder
	printMultiWayTree(n, &sb, "", "")
	return sb.String()
}

func printMultiWayTree(n MultiWayTreeNode, sb *strings.Builder, prefix, childPrefix string) {
	sb.WriteString(prefix)
	if x, ok := n.(interface{ String() string }); ok {
		sb.WriteString(x.String())
	} else {
		sb.WriteString(fmt.Sprint(n))
	}
	sb.WriteString("\n")
	it := n.Iterator()
	if it == nil {
		return
	}
	for it.HasNext() {
		next := it.Next()
		if it.HasNext() {
			printMultiWayTree(next, sb, childPrefix+"├── ", childPrefix+"│   ")
		} else {
			printMultiWayTree(next, sb, childPrefix+"└── ", childPrefix+"    ")
		}
	}
}
