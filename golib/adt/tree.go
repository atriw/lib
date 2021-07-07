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
	printTree(n, &sb, "", "", -1)
	return sb.String()
}

func printTree(n TreeNode, sb *strings.Builder, prefix, childPrefix string, depth int) {
	sb.WriteString(prefix)
	if x, ok := n.(interface{ String() string }); ok {
		sb.WriteString(x.String())
	} else {
		sb.WriteString(fmt.Sprint(n))
	}
	sb.WriteString("\n")
	if n.External() || depth == 0 {
		return
	}
	printTree(n.Right(), sb, childPrefix+"├── ", childPrefix+"│   ", depth-1)
	printTree(n.Left(), sb, childPrefix+"└── ", childPrefix+"    ", depth-1)
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
	printMultiWayTree(n, &sb, "", "", -1)
	return sb.String()
}

func PrintMultiWayTreeDepth(n MultiWayTreeNode, depth int) string {
	var sb strings.Builder
	printMultiWayTree(n, &sb, "", "", depth)
	return sb.String()
}

func printMultiWayTree(n MultiWayTreeNode, sb *strings.Builder, prefix, childPrefix string, depth int) {
	sb.WriteString(prefix)
	if x, ok := n.(interface{ String() string }); ok {
		sb.WriteString(x.String())
	} else {
		sb.WriteString(fmt.Sprint(n))
	}
	sb.WriteString("\n")
	it := n.Iterator()
	if it == nil || depth == 0 {
		return
	}
	for it.HasNext() {
		next := it.Next()
		if it.HasNext() {
			printMultiWayTree(next, sb, childPrefix+"├── ", childPrefix+"│   ", depth-1)
		} else {
			printMultiWayTree(next, sb, childPrefix+"└── ", childPrefix+"    ", depth-1)
		}
	}
}
