package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache_Get(t *testing.T) {
	t.Parallel()

	const cachedData = "test-data"
	calledOK := 0
	fetchDataAlwaysOk := LocalCacheRefreshFunc(func(io IdentifiableObject) ([]byte, error) {
		calledOK++
		return []byte(cachedData), nil
	})

	fetchDataErr := fmt.Errorf("test error")
	calledFail := 0
	fetchDataAlwaysFail := LocalCacheRefreshFunc(func(io IdentifiableObject) ([]byte, error) {
		calledFail++
		return nil, fetchDataErr
	})

	obj := IdentifiableString("foo")
	t.Run("Never Cached", func(t *testing.T) {
		defer t.Cleanup(func() {
			calledFail = 0
			calledOK = 0
		})

		cache := NewCache(0, 0, -1)

		data, err := cache.Get(obj, fetchDataAlwaysFail)
		assert.Equal(t, 1, calledFail)
		assert.ErrorIs(t, err, fetchDataErr)
		assert.Nil(t, data)

		data, err = cache.Get(obj, fetchDataAlwaysOk)
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))
		assert.Equal(t, 1, calledOK)
		_, ok := cache.cache.Get(obj.ScopeUniqueID())
		assert.False(t, ok)

		data, err = cache.Get(obj, fetchDataAlwaysOk)
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))
		assert.Equal(t, 2, calledOK)
		_, ok = cache.cache.Get(obj.ScopeUniqueID())
		assert.False(t, ok)
	})

	t.Run("Cached", func(t *testing.T) {
		defer t.Cleanup(func() {
			calledFail = 0
			calledOK = 0
		})

		cache := NewCache(-1, -1, 0)

		data, err := cache.Get(obj, fetchDataAlwaysOk)
		assert.NoError(t, err)
		assert.Equal(t, cachedData, string(data))

		_, ok := cache.cache.Get(obj.ScopeUniqueID())
		assert.True(t, ok)
	})
}
