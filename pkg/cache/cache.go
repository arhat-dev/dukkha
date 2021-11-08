package cache

import (
	"math"

	lru "github.com/die-net/lrucache"
)

// NewCache is a in memory cache
//
// maxItemBytes < 0, no limit
// 				> 0, only data within this limit can be cached
// 				== 0, disable in memory caching
//
// maxBytes < 0, no limit
// 			> 0, limit total cached in memory data with this size
// 			== 0, disable in memory caching
//
// maxAgeSeconds < 0, always fetch data
// 				 > 0, data become invalid after this duration
// 				 == 0, defaults to 5
func NewCache(maxItemBytes, maxBytes, maxAgeSeconds int64) *Cache {
	if maxBytes < 0 {
		maxBytes = math.MaxInt64
	}

	if maxItemBytes < 0 {
		maxItemBytes = math.MaxInt64
	}

	if maxAgeSeconds == 0 {
		maxAgeSeconds = 5
	}

	return &Cache{
		maxItemSize: maxItemBytes,
		cache:       lru.New(maxBytes, maxAgeSeconds),
	}
}

type Cache struct {
	maxItemSize int64
	cache       *lru.LruCache
}

func (c *Cache) Get(obj IdentifiableObject, refresh LocalCacheRefreshFunc) ([]byte, error) {
	if c.cache.MaxAge < 0 {
		return refresh(obj)
	}

	key := obj.ScopeUniqueID()
	data, ok := c.cache.Get(key)
	if ok {
		return data, nil
	}

	data, err := refresh(obj)
	if err != nil {
		return nil, err
	}

	if int64(len(data)) <= c.maxItemSize {
		c.cache.Set(key, data)
	}

	return data, nil
}
