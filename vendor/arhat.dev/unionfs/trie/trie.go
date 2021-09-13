package trie

// New creates a new trie tree with optional keyFunc
func New(keyFunc KeyFunc) *Trie {
	if keyFunc == nil {
		keyFunc = RuneKeyFunc
	}

	return &Trie{
		nextKey: keyFunc,
		root:    newNode(nil, ""),
	}
}

type Trie struct {
	nextKey KeyFunc

	root *Node
}

// Add node
func (t *Trie) Add(key string, value interface{}) *Node {
	last := t.root

	t.doEachElement(key, func(k string) bool {
		current := last.get(k)
		if current == nil {
			current = newNode(last, k)
			last.add(current)
		}
		last = current
		return true
	})

	if last == t.root {
		return nil
	}

	last.value = value

	return last
}

func (t *Trie) Get(key string) (_ *Node, exact bool) {
	ret := t.root
	exact = true

	t.doEachElement(key, func(k string) bool {
		node := ret.get(k)
		if node == nil {
			exact = false
			return false
		}
		ret = node
		return true
	})

	if ret == t.root {
		ret = nil
		exact = false
	}

	return ret, exact
}

func (t *Trie) Delete(key string) (toDel *Node, deleted bool) {
	toDel = t.find(key)
	deleted = toDel != nil

	if deleted {
		toDel.deleteSelf()
	}

	return
}

func (t *Trie) find(key string) *Node {
	n := t.root
	t.doEachElement(key, func(k string) bool {
		n = n.get(k)
		return n != nil
	})

	if n == t.root {
		return nil
	}

	return n
}

func (t *Trie) doEachElement(key string, do func(k string) bool) {
	for k, i := t.nextKey(key, 0); i > 0; k, i = t.nextKey(key, i) {
		if !do(k) {
			return
		}
	}
}
