package diff

import (
	"fmt"
	"strconv"
	"strings"
	_ "unsafe" // for go:linkname

	_ "arhat.dev/rs" // add required references for go:linkname
	"gopkg.in/yaml.v3"
)

func newNode(parent *Node, elemKey string) *Node {
	return &Node{
		elemKey: elemKey,

		parent: parent,
	}
}

type RendererSpec struct {
	Name     string
	TypeHint string
	Patch    bool
}

func (r *RendererSpec) Clone() *RendererSpec {
	return &RendererSpec{
		Name:     r.Name,
		TypeHint: r.TypeHint,
		Patch:    r.Patch,
	}
}

// maintain the same behavior as
// https://github.com/arhat-dev/rs/blob/master/field.go#L403
func parseRenderingSuffix(suffix string) []*RendererSpec {
	var (
		parts = strings.Split(suffix, "|")
		ret   []*RendererSpec
	)

	for _, part := range parts {
		size := len(part)
		if size == 0 {
			continue
		}

		spec := &RendererSpec{
			Patch: part[size-1] == '!',
		}

		if spec.Patch {
			part = part[:size-1]
			// size-- // size not used any more
		}

		if idx := strings.LastIndexByte(part, '?'); idx >= 0 {
			spec.TypeHint = part[idx+1:]
			part = part[:idx]
		}

		spec.Name = part
		ret = append(ret, spec)
	}

	return ret
}

type Node struct {
	elemKey string
	parent  *Node

	Renderers []*RendererSpec
	RawNode   *yaml.Node

	scalarData *yaml.Node

	children   []*Node
	childIndex map[string]int
}

func (n *Node) MarshalYAML() (interface{}, error) {
	return n.RawNode, nil
}

//go:linkname unmarshalMap arhat.dev/rs.unmarshalYamlMap
func unmarshalMap(content []*yaml.Node) ([][]*yaml.Node, error)

func (n *Node) UnmarshalYAML(yn *yaml.Node) error {
	n.RawNode = yn

	switch yn.Kind {
	case yaml.MappingNode:
		if n.childIndex == nil {
			n.childIndex = make(map[string]int)
		}

		pairs, err := unmarshalMap(yn.Content)
		if err != nil {
			return err
		}

		for _, pair := range pairs {
			k := pair[0].Value

			var rsSpecs []*RendererSpec
			// handle rendering suffix first
			// maintain the same behavior as
			// https://github.com/arhat-dev/rs/blob/master/unmarshal.go#L40
			suffixStart := strings.LastIndexByte(k, '@')
			if suffixStart != -1 {
				rsSpecs = parseRenderingSuffix(k[suffixStart+1:])
				k = k[:suffixStart]
			}

			if strings.IndexByte(k, '.') != -1 {
				k = strconv.Quote(k)
			}
			k = "." + k
			v := pair[1]

			child := newNode(n, k)
			err := child.UnmarshalYAML(v)
			if err != nil {
				return err
			}

			child.Renderers = rsSpecs
			n.children = append(n.children, child)
			n.childIndex[k] = len(n.children) - 1
		}

		return nil
	case yaml.SequenceNode:
		if n.childIndex == nil {
			n.childIndex = make(map[string]int)
		}

		for i, v := range yn.Content {
			for v.Alias != nil {
				v = v.Alias
			}

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
func (n *Node) Get(key []string) (_ *Node, tailKey []string) {
	switch {
	case len(key) == 0:
		// no key for this node (selected by upper level)
		// then this node is the exact one
		return n, nil
	case len(n.children) == 0:
		return n, key
	default:
		i, ok := n.childIndex[key[0]]
		if !ok {
			return n, key
		}

		return n.children[i].Get(key[1:])
	}
}

func (n *Node) Clone() *Node {
	clone := &Node{
		elemKey:    n.elemKey,
		parent:     n.parent,
		Renderers:  nil,
		RawNode:    n.RawNode,
		scalarData: n.scalarData,
		children:   nil,
		childIndex: nil,
	}

	for _, rdr := range n.Renderers {
		clone.Renderers = append(clone.Renderers, rdr.Clone())
	}

	for _, child := range n.children {
		clone.children = append(clone.children, child.Clone())
	}

	for k, v := range n.childIndex {
		if clone.childIndex == nil {
			clone.childIndex = make(map[string]int)
		}

		clone.childIndex[k] = v
	}

	return clone
}

func (n *Node) Append(other *Node) {
	if n.childIndex == nil {
		n.childIndex = make(map[string]int)
	}

	n.children = append(n.children, other)
}
