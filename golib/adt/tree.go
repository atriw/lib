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
