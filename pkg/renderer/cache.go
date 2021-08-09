package renderer

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"arhat.dev/pkg/hashhelper"
	lru "github.com/die-net/lrucache"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/output"
	"arhat.dev/dukkha/pkg/utils"
)

type CacheConfig struct {
	rs.BaseField

	EnableCache    bool          `yaml:"enable_cache"`
	CacheSizeLimit utils.Size    `yaml:"cache_size_limit"`
	CacheMaxAge    time.Duration `yaml:"cache_max_age"`
}

type CacheRefreshFunc func(key string) ([]byte, error)

func NewCache(limit int64, expiry time.Duration) *Cache {
	return &Cache{
		cache: lru.New(limit, int64(expiry.Seconds())),
	}
}

type Cache struct {
	cache *lru.LruCache
}

func (c *Cache) Get(key string, refresh CacheRefreshFunc) ([]byte, error) {
	data, ok := c.cache.Get(key)
	if ok {
		return data, nil
	}

	data, err := refresh(key)
	if err != nil {
		return nil, err
	}

	c.cache.Set(key, data)
	return data, nil
}

func CreateRefreshFuncForRemote(
	cacheDir string,
	maxCacheAge time.Duration,
	doRemoteFetch CacheRefreshFunc,
) CacheRefreshFunc {
	return func(key string) ([]byte, error) {
		localCacheFilePrefix := hex.EncodeToString(hashhelper.MD5Sum([]byte(key)))

		// find from local cache
		// ${DUKKHA_CACHE_DIR}/renderer-<rendererName>/<md5sum(key)>-<unix-timestamp>

		var expiredLatestLocalCache string

		entries, err := os.ReadDir(cacheDir)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf(
					"failed to check local cache dir: %w", err,
				)
			}

			// no cache exists

			err = os.MkdirAll(cacheDir, 0750)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to ensure local cache dir: %w", err,
				)
			}
			// fetch from remote
		} else if len(entries) > 0 {
			start := sort.Search(len(entries), func(i int) bool {
				return strings.HasPrefix(entries[i].Name(), localCacheFilePrefix)
			})

			switch start {
			case len(entries):
				// (not found) do nothing
			default:
				latestAt := start
				for ; latestAt+1 < len(entries); latestAt++ {
					if !strings.HasPrefix(entries[latestAt+1].Name(), localCacheFilePrefix) {
						break
					}
				}

				for _, info := range entries[start:latestAt] {
					_ = os.Remove(filepath.Join(cacheDir, info.Name()))
				}

				targetFile := entries[latestAt].Name()
				targetPath := filepath.Join(cacheDir, targetFile)

				parts := strings.SplitN(targetFile, "-", 2)
				if len(parts) != 2 || parts[0] != localCacheFilePrefix {
					// invalid cache file
					return nil, fmt.Errorf(
						"invalid cache file, please remove %q", targetPath,
					)
				}

				timestamp, err2 := strconv.ParseInt(
					// trim padding
					strings.TrimLeft(parts[1], "0"),
					10, 64,
				)
				if err2 != nil {
					return nil, fmt.Errorf(
						"invalid timestamp, please remove local cache file %q: %w",
						targetPath, err2,
					)
				}

				if time.Since(time.Unix(timestamp, 0)) < maxCacheAge {
					return os.ReadFile(targetPath)
				}

				// cache expired, but do not remove unless remote fetch is successful
				// _ = os.Remove(targetPath)
				expiredLatestLocalCache = targetPath
			}
		}

		// pad timestamp to get future results sorted by os.ReadDir
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		// int64 can have at most 20 digits
		timestamp = strings.Repeat("0", 20-len(timestamp)) + timestamp

		localCacheFile := filepath.Join(
			cacheDir, localCacheFilePrefix+"-"+timestamp,
		)

		data, err := doRemoteFetch(key)
		if err != nil {
			if len(expiredLatestLocalCache) != 0 {
				output.WriteUsingExpiredCacheWarning(key)
				return os.ReadFile(expiredLatestLocalCache)
			}

			return nil, err
		}

		err = os.WriteFile(localCacheFile, data, 0640)
		if err != nil {
			// remove incomplete local cache
			_ = os.Remove(localCacheFile)

			if len(expiredLatestLocalCache) != 0 {
				output.WriteUsingExpiredCacheWarning(key)
				return os.ReadFile(expiredLatestLocalCache)
			}

			return nil, err
		}

		// remove expired local cache since refresh succeeded
		if len(expiredLatestLocalCache) != 0 {
			_ = os.Remove(expiredLatestLocalCache)
		}

		return data, nil
	}
}
