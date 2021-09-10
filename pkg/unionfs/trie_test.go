package unionfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathElementKeyFunc(t *testing.T) {
	assert.EqualValues(t, []string{"foo", "foo", "test"}, PathReverseKeysFunc("/test/foo/foo"))
}

func TestTrie_New(t *testing.T) {

	t.Run("Defaul", func(t *testing.T) {
		tr := NewTrie(nil)
		assert.NotNil(t, tr.GetReversedKeys)
		assert.NotNil(t, tr.Root)
	})

	t.Run("Custom ElementKeyFunc", func(t *testing.T) {
		testElemFunc := ReverseKeysFunc(func(key string) []string {
			return nil
		})

		tr := NewTrie(testElemFunc)
		// assert.EqualValues(t,
		// 	tr.GetElementKey,
		// 	testElemFunc,
		// )

		assert.NotNil(t, tr.Root)
	})
}

func TestTrie_Add(t *testing.T) {
	tr := NewTrie(nil)
	tr.Add("/test/foo/foo", "foo")
	assert.Len(t, tr.Root.children, 1)
	assert.Len(t, tr.Root.get("test").children, 1)
	assert.Len(t, tr.Root.get("test").get("foo").children, 1)
	assert.Len(t, tr.Root.get("test").get("foo").get("foo").children, 0)
	assert.Equal(t, tr.Root.get("test").get("foo").get("foo").Value(), "foo")
	assert.Equal(t, tr.Root.get("test").get("foo").get("foo").value, "foo")

	tr.Add("/test/foo/bar", "bar")
	assert.Len(t, tr.Root.children, 1)
	assert.Len(t, tr.Root.get("test").children, 1)
	assert.Len(t, tr.Root.get("test").get("foo").children, 2)
	assert.Len(t, tr.Root.get("test").get("foo").get("bar").children, 0)
	assert.Equal(t, tr.Root.get("test").get("foo").get("bar").Value(), "bar")
	assert.Equal(t, tr.Root.get("test").get("foo").get("bar").value, "bar")

	tr.Add("test/bar", "foo-bar")
	assert.Len(t, tr.Root.children, 1)
	assert.Len(t, tr.Root.get("test").children, 2)
	assert.Len(t, tr.Root.get("test").get("bar").children, 0)
	assert.Equal(t, tr.Root.get("test").get("bar").Value(), "foo-bar")
	assert.Equal(t, tr.Root.get("test").get("bar").value, "foo-bar")
}

func TestTrie_Get(t *testing.T) {
	tr := NewTrie(nil)
	tr.Add("test/foo/bar/woo", "ok")
	exact, nearest, ok := tr.Get("test/foo/bar/woo")
	assert.True(t, ok)
	assert.NotNil(t, exact)
	assert.Nil(t, nearest)
	assert.Equal(t, exact.elemKey, "woo")

	exact, nearest, ok = tr.Get("test/foo/bar/nop")
	assert.False(t, ok)
	assert.Nil(t, exact)
	assert.NotNil(t, nearest)
	assert.Equal(t, nearest.elemKey, "bar")
}
