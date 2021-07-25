package renderer

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		const (
			cacheSize = 10
			dataSize  = cacheSize / 2
		)

		data := make([]byte, dataSize)
		cache := NewCache(cacheSize, time.Second)
		dataFromCache, err := cache.Get("foo", func(key string) ([]byte, error) {
			return data, nil
		})

		assert.NoError(t, err)
		assert.EqualValues(t, data, dataFromCache)
	})

	t.Run("Excessive Data", func(t *testing.T) {
		const (
			cacheSize = 10
			dataSize  = cacheSize * 10
		)

		data := make([]byte, dataSize)
		cache := NewCache(cacheSize, time.Second)
		dataFromCache, err := cache.Get("foo", func(key string) ([]byte, error) {
			return data, nil
		})

		assert.NoError(t, err)
		assert.EqualValues(t, data, dataFromCache)
	})

	t.Run("Cache With Remote Fetch", func(t *testing.T) {
		const (
			cacheSize   = 10
			dataSize    = cacheSize / 2
			cachePeriod = time.Second
		)

		successData := bytes.Repeat([]byte("1"), dataSize)
		cache := NewCache(cacheSize, cachePeriod)

		cacheDir, err := os.MkdirTemp(os.TempDir(), "dukkha-test-cache-*")
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			_ = os.RemoveAll(cacheDir)
		}()

		refreshSuccessFunc := CreateRefreshFuncForRemote(
			cacheDir,
			cachePeriod,
			func(key string) ([]byte, error) {
				return successData, nil
			},
		)

		// initial fresh cache
		dataFromCache, err := cache.Get("foo", refreshSuccessFunc)
		assert.NoError(t, err)
		assert.EqualValues(t, successData, dataFromCache)

		entries, err := os.ReadDir(cacheDir)
		assert.NoError(t, err)
		assert.Len(t, entries, 1)

		// immediate retrieve cache
		dataFromCache, err = cache.Get("foo", refreshSuccessFunc)
		assert.NoError(t, err)
		assert.EqualValues(t, successData, dataFromCache)

		newEntries, err := os.ReadDir(cacheDir)
		assert.NoError(t, err)
		assert.Len(t, newEntries, 1)
		assert.Equal(t, entries[0].Name(), newEntries[0].Name())

		// wait to expire local cache to test local cache refreshing
		//
		// local cache file should have different name
		time.Sleep(cachePeriod * 2)

		dataFromCache, err = cache.Get("foo", refreshSuccessFunc)
		assert.NoError(t, err)
		assert.EqualValues(t, successData, dataFromCache)

		newEntries, err = os.ReadDir(cacheDir)
		assert.NoError(t, err)
		assert.Len(t, newEntries, 1)
		assert.NotEqual(t, entries[0].Name(), newEntries[0].Name())

		// always keep last cached data, no matter expired or not
		time.Sleep(cachePeriod * 2)

		refreshErrorFunc := CreateRefreshFuncForRemote(
			cacheDir, cachePeriod,
			func(key string) ([]byte, error) {
				return nil, fmt.Errorf("always error")
			},
		)

		// should print warning and get previously cached but expired data
		dataFromCache, err = cache.Get("foo", refreshErrorFunc)
		assert.NoError(t, err)
		assert.EqualValues(t, successData, dataFromCache)

		expiredEntries, err := os.ReadDir(cacheDir)
		assert.NoError(t, err)
		assert.Len(t, expiredEntries, 1)
		assert.Equal(t, newEntries[0].Name(), expiredEntries[0].Name())
	})
}
