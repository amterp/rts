package rts

import (
	"fmt"
	"strings"

	ts "github.com/tree-sitter/go-tree-sitter"
)

var escapedReplacer = strings.NewReplacer(
	"\n", "\\n",
	"\r", "\\r",
	"\t", "\\t",
)

func (rt *RtsTree) Dump() string {
	node := rt.root.RootNode()
	maxByte, maxPosRow, maxPosCol := findMaxRanges(node, 0, 0, 0)

	byteLen := len(fmt.Sprintf("%d", maxByte))
	rowLen := len(fmt.Sprintf("%d", maxPosRow))
	colLen := len(fmt.Sprintf("%d", maxPosCol))
	fmtString := fmt.Sprintf("B: [%%%dd, %%%dd] PS: [%%%dd, %%%dd], PE: [%%%dd, %%%dd] %%s%%s",
		byteLen, byteLen, rowLen, colLen, rowLen, colLen)

	var sb strings.Builder
	rt.recurseAppendString(&sb, fmtString, node, 0)

	return sb.String()
}

func findMaxRanges(node *ts.Node, maxByte uint, maxPosRow uint, maxPosCol uint) (uint, uint, uint) {
	if node.EndByte() > maxByte {
		maxByte = node.EndByte()
	}
	if node.EndPosition().Row > maxPosRow {
		maxPosRow = node.EndPosition().Row
	}
	if node.EndPosition().Column > maxPosCol {
		maxPosCol = node.EndPosition().Column
	}

	children := node.Children(node.Walk())
	for _, child := range children {
		maxByte, maxPosRow, maxPosCol = findMaxRanges(&child, maxByte, maxPosRow, maxPosCol)
	}
	return maxByte, maxPosRow, maxPosCol
}

func (rt *RtsTree) recurseAppendString(sb *strings.Builder, fmtString string, node *ts.Node, treeLevel int) {
	indent := strings.Repeat("  ", treeLevel)
	sb.WriteString(fmt.Sprintf(fmtString,
		node.StartByte(), node.EndByte(),
		node.StartPosition().Row, node.StartPosition().Column,
		node.EndPosition().Row, node.EndPosition().Column,
		indent,
		node.Kind(),
	))

	children := node.Children(node.Walk())
	if len(children) == 0 {
		src := rt.src[node.StartByte():node.EndByte()]
		sb.WriteString(fmt.Sprintf(" `%s`\n", escapedReplacer.Replace(src)))
		return
	} else {
		sb.WriteString("\n")
	}

	for _, child := range children {
		rt.recurseAppendString(sb, fmtString, &child, treeLevel+1)
	}
}
