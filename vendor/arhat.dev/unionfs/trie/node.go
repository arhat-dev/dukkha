package trie

import "sort"

func newNode(parent *Node, elemKey string) *Node {
	return &Node{
		elemKey: elemKey,

		parent:   parent,
		children: make(map[string]*Node),
	}
}

type Node struct {
	elemKey string
	value   interface{}

	parent   *Node
	children map[string]*Node
}

func (n *Node) Key() string {
	if n.parent == nil {
		return "" + n.elemKey
	}

	return n.parent.Key() + n.elemKey
}

func (n *Node) ElementKey() string {
	return n.elemKey
}

func (n *Node) Value() interface{} {
	return n.value
}

func (n *Node) Children() []*Node {
	var ret []*Node
	for k := range n.children {
		ret = append(ret, n.children[k])
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].elemKey < ret[j].elemKey
	})

	return ret
}

func (n *Node) add(node *Node) {
	n.children[node.elemKey] = node
}

func (n *Node) get(elemKey string) *Node {
	return n.children[elemKey]
}

func (n *Node) deleteSelf() {
	if n.parent != nil {
		delete(n.parent.children, n.elemKey)
	}
}
