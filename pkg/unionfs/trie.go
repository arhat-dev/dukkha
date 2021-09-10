package unionfs

import (
	"path"
	"sort"
	"strings"
)

func PathReverseKeysFunc(key string) []string {
	key = strings.TrimSuffix(key, "/")

	var keys []string
	for dir, file := path.Split(key); len(file) != 0; {
		keys = append(keys, file)
		dir, file = path.Split(strings.TrimSuffix(dir, "/"))
	}

	return keys
}

// NewTrie creates
// reversedKeysGenFunc should generate keys in reversed order
// e.g. `/a/b/c` shoud generate `[]string{"c", "b", "a"}`
func NewTrie(reversedKeysGenFunc ReverseKeysFunc) *Trie {
	if reversedKeysGenFunc == nil {
		reversedKeysGenFunc = PathReverseKeysFunc
	}

	return &Trie{
		GetReversedKeys: reversedKeysGenFunc,
		Root:            newNode("", nil),
	}
}

type ReverseKeysFunc func(key string) []string

type Trie struct {
	GetReversedKeys ReverseKeysFunc

	Root *Node
}

// Add node
func (t *Trie) Add(key string, value interface{}) bool {
	keys := t.GetReversedKeys(key)
	if len(keys) == 0 {
		return false
	}

	node := newNode(keys[0], value)

	keys = keys[1:]
	for i, k := range keys {
		parentNode := t.getNode(keys[i:])
		if parentNode == nil {
			parentNode = newNode(k, nil)
		}
		parentNode.add(node)

		node = parentNode
	}

	t.Root.add(node)

	return true
}

func (t *Trie) Get(key string) (exact, nearest *Node, found bool) {
	keys := t.GetReversedKeys(key)
	exact = t.getNode(keys)
	found = exact != nil

	if !found {
		for nearest == nil && len(keys) != 0 {
			keys = keys[1:]
			nearest = t.getNode(keys)
		}
	}

	return
}

func (t *Trie) getNode(reversedKeys []string) *Node {
	node := t.Root

	for i := len(reversedKeys) - 1; i >= 0; i-- {
		node = node.get(reversedKeys[i])
		if node == nil {
			return nil
		}
	}

	return node
}

func newNode(elemKey string, val interface{}) *Node {
	return &Node{
		elemKey: elemKey,
		value:   val,

		children: make(map[string]*Node),
	}
}

type Node struct {
	elemKey string
	value   interface{}

	children map[string]*Node
}

func (n *Node) ElementKey() string {
	return n.elemKey
}

func (n *Node) Value() interface{} {
	return n.value
}

func (n *Node) Children() []Node {
	var ret []Node
	for k := range n.children {
		ret = append(ret, *n.children[k])
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
