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

type rsSpec struct {
	renderer string
	typeHint string
	patch    bool
}

// maintain the same behavior as
// https://github.com/arhat-dev/rs/blob/master/field.go#L403
func parseRenderingSuffix(suffix string) []*rsSpec {
	var (
		parts = strings.Split(suffix, "|")
		ret   []*rsSpec
	)

	for _, part := range parts {
		size := len(part)
		if size == 0 {
			continue
		}

		spec := &rsSpec{
			patch: part[size-1] == '!',
		}

		if spec.patch {
			part = part[:size-1]
			// size-- // size not used any more
		}

		if idx := strings.LastIndexByte(part, '?'); idx >= 0 {
			spec.typeHint = part[idx+1:]
			part = part[:idx]
		}

		spec.renderer = part
		ret = append(ret, spec)
	}

	return ret
}

type Node struct {
	elemKey string
	parent  *Node

	rsSpecs []*rsSpec
	raw     *yaml.Node

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

		pairs, err := unmarshalMap(yn)
		if err != nil {
			return err
		}

		for _, pair := range pairs {
			k := pair[0].Value

			var rsSpecs []*rsSpec
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

			child.rsSpecs = rsSpecs
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

//go:linkname unmarshalMap arhat.dev/rs.unmarshalYamlMap
func unmarshalMap(n *yaml.Node) ([][]*yaml.Node, error)
