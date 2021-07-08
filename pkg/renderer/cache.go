package renderer

import (
	"time"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/utils"
	lru "github.com/die-net/lrucache"
)

type CacheConfig struct {
	field.BaseField

	EnableCache    bool          `yaml:"enable_cache"`
	CacheSizeLimit utils.Size    `yaml:"cache_size_limit"`
	CacheMaxAge    time.Duration `yaml:"cache_max_age"`
}

type CacheRefreshFunc func(key string) ([]byte, error)

func NewCache(
	limit int64,
	expiry time.Duration,
	refresh CacheRefreshFunc,
) *Cache {
	return &Cache{
		refresh: refresh,
		cache:   lru.New(int64(limit), int64(expiry.Seconds())),
	}
}

type Cache struct {
	refresh CacheRefreshFunc
	cache   *lru.LruCache
}

func (c *Cache) Get(key string) ([]byte, error) {
	data, ok := c.cache.Get(key)
	if ok {
		return data, nil
	}

	data, err := c.refresh(key)
	if err != nil {
		return nil, err
	}

	c.cache.Set(key, data)
	return data, nil
}
