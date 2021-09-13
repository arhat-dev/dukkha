package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrie_New(t *testing.T) {

	t.Run("Defaul", func(t *testing.T) {
		tr := New(nil)
		assert.NotNil(t, tr.nextKey)
		assert.NotNil(t, tr.root)
	})

	t.Run("Custom ElementKeyFunc", func(t *testing.T) {
		testElemFunc := KeyFunc(func(key string, pos int) (string, int) {
			return "", -1
		})

		tr := New(testElemFunc)
		// assert.EqualValues(t,
		// 	tr.GetElementKey,
		// 	testElemFunc,
		// )

		assert.NotNil(t, tr.root)
	})
}

func TestTrie_Add(t *testing.T) {
	tr := New(nil)

	assert.Nil(t, tr.Add("", "anything"))

	assert.NotNil(t, tr.Add("ABB", "foo"))
	assert.Len(t, tr.root.children, 1)
	assert.Len(t, tr.root.get("A").Children(), 1)
	assert.Len(t, tr.root.get("A").get("B").children, 1)
	assert.Len(t, tr.root.get("A").get("B").get("B").children, 0)
	assert.Equal(t, "B", tr.root.get("A").get("B").get("B").elemKey)
	assert.Equal(t, "B", tr.root.get("A").get("B").get("B").ElementKey())
	assert.Equal(t, "foo", tr.root.get("A").get("B").get("B").Value())
	assert.Equal(t, "foo", tr.root.get("A").get("B").get("B").value)
	assert.Equal(t, "ABB", tr.root.get("A").get("B").get("B").Key())

	assert.NotNil(t, tr.Add("ABC", "bar"))
	assert.Len(t, tr.root.children, 1)
	assert.Len(t, tr.root.get("A").Children(), 1)
	assert.Len(t, tr.root.get("A").get("B").children, 2)
	assert.Len(t, tr.root.get("A").get("B").get("C").children, 0)
	assert.Equal(t, "C", tr.root.get("A").get("B").get("C").elemKey)
	assert.Equal(t, "C", tr.root.get("A").get("B").get("C").ElementKey())
	assert.Equal(t, "bar", tr.root.get("A").get("B").get("C").Value())
	assert.Equal(t, "bar", tr.root.get("A").get("B").get("C").value)
	assert.Equal(t, "ABC", tr.root.get("A").get("B").get("C").Key())

	assert.NotNil(t, tr.Add("AC", "foo-bar"))
	assert.Len(t, tr.root.children, 1)
	assert.Len(t, tr.root.get("A").Children(), 2)
	assert.Len(t, tr.root.get("A").get("C").children, 0)
	assert.Equal(t, "foo-bar", tr.root.get("A").get("C").Value())
	assert.Equal(t, "foo-bar", tr.root.get("A").get("C").value)
	assert.Equal(t, "AC", tr.root.get("A").get("C").Key())
}

func TestTrie_Get(t *testing.T) {
	tr := New(nil)
	ret, exact := tr.Get("")
	assert.Nil(t, ret)
	assert.False(t, exact)

	assert.NotNil(t, tr.Add("ABC", "nearest-match"))
	assert.NotNil(t, tr.Add("ABCD", "exact-match"))

	ret, exact = tr.Get("ABCD")
	assert.True(t, exact)
	assert.NotNil(t, ret)
	assert.Equal(t, "D", ret.elemKey)
	assert.Equal(t, "exact-match", ret.value)

	ret, exact = tr.Get("ABCE")
	assert.False(t, exact)
	assert.NotNil(t, ret)
	assert.Equal(t, "C", ret.elemKey)
	assert.Equal(t, "nearest-match", ret.value)
}

func TestTrie_DeletePrefix(t *testing.T) {
	tr := New(nil)
	del, ok := tr.Delete("")
	assert.False(t, ok)
	assert.Nil(t, del)

	assert.NotNil(t, tr.Add("ABCD", "G"))
	del, ok = tr.Delete("ABC")
	assert.True(t, ok)
	assert.Equal(t, del.elemKey, "C")
	assert.Len(t, tr.root.get("A").get("B").children, 0)

	del, ok = tr.Delete("ABCDEF")
	assert.False(t, ok)
	assert.Nil(t, del)
}
