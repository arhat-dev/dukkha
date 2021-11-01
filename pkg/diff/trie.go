package diff

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func newNode(parent *Node, elemKey string) *Node {
	return &Node{
		elemKey: elemKey,

		parent: parent,
	}
}

type Node struct {
	elemKey string
	parent  *Node

	raw *yaml.Node

	scalarData *yaml.Node

	children   []*Node
	childIndex map[string]int
}

func (n *Node) MarshalYAML() (interface{}, error) {
	return n.raw, nil
}

func (n *Node) UnmarshalYAML(yn *yaml.Node) error {
	n.raw = yn

	switch yn.Kind {
	case yaml.MappingNode:
		if n.childIndex == nil {
			n.childIndex = make(map[string]int)
		}

		// TODO: handle merging and rendering suffix
		for i := 0; i < len(yn.Content); i += 2 {
			k := yn.Content[i].Value
			if strings.IndexByte(k, '.') != -1 {
				k = strconv.Quote(k)
			}
			k = "." + k
			v := yn.Content[i+1]

			child := newNode(n, k)
			err := child.UnmarshalYAML(v)
			if err != nil {
				return err
			}

			n.children = append(n.children, child)
			n.childIndex[k] = i / 2
		}

		return nil
	case yaml.SequenceNode:
		if n.childIndex == nil {
			n.childIndex = make(map[string]int)
		}

		// TODO: handle merging
		for i, v := range yn.Content {
			idx := "[" + strconv.FormatInt(int64(i), 10) + "]"

			child := newNode(n, idx)
			err := child.UnmarshalYAML(v)
			if err != nil {
				return err
			}

			n.children = append(n.children, child)
			n.childIndex[idx] = i
		}

		return nil
	case yaml.ScalarNode:
		n.scalarData = yn
		return nil
	default:
		return fmt.Errorf("unsupported kind %q", yn.Kind)
	}
}

func (n *Node) Key() []string {
	if n.parent == nil {
		return []string{n.elemKey}
	}

	return append(n.parent.Key(), n.elemKey)
}

func (n *Node) ElementKey() string {
	return n.elemKey
}

// Get a trie node according to the key sequence in order
// exact is set to true when there
// key always refer to node's children
func (n *Node) Get(key []string) (_ *Node, exact bool) {
	switch {
	case len(key) == 0:
		// no key for this node (selected by upper level)
		// then this node is the exact one
		return n, true
	case len(n.children) == 0:
		return n, false
	default:
		i, ok := n.childIndex[key[0]]
		if !ok {
			return n, false
		}

		return n.children[i].Get(key[1:])
	}
}
