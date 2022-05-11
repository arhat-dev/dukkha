package cache

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/log"
	"arhat.dev/pkg/md5helper"
	lru "github.com/die-net/lrucache"
)

// NewTwoTierCache
//
// itemMaxBytes < 0, no limit to item size
// 				> 0, only items with size below can be cached
// 				== 0, in memory caching disabled
//
// maxBytes < 0, no limit to total cache size
// 			> 0, limit cache size to maxBytes
// 			== 0, in memory caching disabled
//
// maxAgeSeconds <= 0, once cached in memory, always valid during runtime,
// 					   but will always fetch from remote if in memory cache lost
// 				 > 0, limit both in memory and local file cache to this long.
func NewTwoTierCache(
	cacheFS *fshelper.OSFS,
	itemMaxBytes, maxBytes, maxAgeSeconds int64,
) *TwoTierCache {
	if maxBytes < 0 {
		maxBytes = math.MaxInt64
	}

	if itemMaxBytes < 0 {
		itemMaxBytes = math.MaxInt64
	}

	return &TwoTierCache{
		itemMaxBytes: itemMaxBytes,

		cacheFS:  cacheFS,
		memcache: lru.New(maxBytes, maxAgeSeconds),
	}
}

type TwoTierCache struct {
	itemMaxBytes int64

	cacheFS  *fshelper.OSFS
	memcache *lru.LruCache
}

// Get cached content
//
// now is the unix timestamp of the time being
func (c *TwoTierCache) Get(
	obj IdentifiableObject,
	now int64,
	allowExpired bool,
	refresh RemoteCacheRefreshFunc,
) ([]byte, bool, error) {
	_, content, expired, err := c.get(obj, now, true, refresh)
	return content, expired, err
}

// GetPath find local file path to cached data
//
// now is the unix timestamp of the time being
func (c *TwoTierCache) GetPath(
	obj IdentifiableObject,
	now int64,
	allowExpired bool,
	refresh RemoteCacheRefreshFunc,
) (string, bool, error) {
	f, _, expired, err := c.get(obj, now, false, refresh)
	return f, expired, err
}

func (c *TwoTierCache) get(
	obj IdentifiableObject,
	now int64,
	retConent bool,
	refresh RemoteCacheRefreshFunc,
) (file string, content []byte, isExpired bool, err error) {
	if retConent {
		var ok bool
		content, ok = c.memcache.Get(obj.ScopeUniqueID())
		if ok {
			return "", content, false, nil
		}
	}

	cacheFilenamePrefix := formatCacheFilenamePrefix(obj.ScopeUniqueID())
	suffix := obj.Ext()
	active, expired, _, err := lookupLocalCache(
		c.cacheFS, cacheFilenamePrefix, suffix, now-c.memcache.MaxAge,
	)
	if err != nil {
		return "", nil, false, err
	}

	// actively remove all but last expired cache
	if len(expired) > 1 {
		for _, v := range expired[:len(expired)-1] {
			// best effort
			_ = c.cacheFS.Chmod(v, 0600)
			err = c.cacheFS.Remove(v)
			if err != nil {
				log.Log.I("removing expired cache",
					log.String("file", v), log.Error(err),
				)
			}
		}
	}

	if len(active) != 0 {
		// use latest active cache
		file = active[len(active)-1]
		isExpired = false
		if retConent {
			content, err = c.cacheFS.ReadFile(file)
		}

		file, err = c.cacheFS.Abs(file)
		return
	}

	// no active cache, fetch from remote

	// first ensure cache dir exists
	_, err = c.cacheFS.Stat(".")
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return
		}

		err = c.cacheFS.MkdirAll(".", 0755)
		if err != nil {
			return
		}
	}

	r, err := refresh(obj)
	if err != nil {
		// failed fetching from remote, fallback to last expired
		if len(expired) == 0 {
			// no expired cache, fail
			return
		}

		file = expired[len(expired)-1]
		isExpired = true

		var err2 error
		if retConent {
			content, err2 = fs.ReadFile(c.cacheFS, file)
			if err2 != nil {
				err = fmt.Errorf("%v: %w", err, err2)
			}
		}

		file, err2 = c.cacheFS.Abs(file)
		if err2 != nil {
			err = fmt.Errorf("%v: %w", err, err2)
		}

		return
	}
	defer func() { _ = r.Close() }()

	isExpired = false
	_file := formatLocalCacheFilename(cacheFilenamePrefix, suffix, now)
	size, content, err := storeLocalCache(c.cacheFS, _file, r, retConent)
	if err != nil {
		return
	}

	file, err = c.cacheFS.Abs(_file)
	if err != nil {
		return
	}

	// no error, handle in memory cache

	if size > c.itemMaxBytes || size > c.memcache.MaxSize {
		// do not cache this item since it's too large
		return
	}

	// do not actively read from cached file
	if retConent {
		c.memcache.Set(obj.ScopeUniqueID(), content)
	}

	return
}

func formatCacheFilenamePrefix(id string) string {
	return hex.EncodeToString(md5helper.Sum([]byte(id)))
}

func formatLocalCacheFilename(prefix, suffix string, now int64) string {
	timestamp := strconv.FormatInt(now, 10)
	// int64 can have at most 20 digits
	timestamp = strings.Repeat("0", 20-len(timestamp)) + timestamp
	return prefix + "-" + timestamp + suffix
}

func storeLocalCache(
	ofs *fshelper.OSFS,
	dest string,
	r io.Reader,
	returnContent bool,
) (int64, []byte, error) {
	if ofs.Chmod(dest, 0600) == nil {
		defer func() { _ = ofs.Chmod(dest, 0400) }()
	}

	f, err := ofs.OpenFile(dest, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0400)
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = f.Close() }()

	var dst io.Writer = f.(*os.File)
	var buf bytes.Buffer
	if returnContent {
		dst = io.MultiWriter(dst, &buf)
	}

	n, err := io.Copy(dst, r)
	if err != nil {
		return 0, nil, err
	}

	if returnContent {
		return n, buf.Next(buf.Len()), nil
	}

	return n, nil, nil
}

// lookupLocalCache to find latest cache file in cacheDir for object
// it will also delete all but last expired cache file
func lookupLocalCache(
	cacheDir fs.FS,
	prefix string,
	// optional suffix to cached file (e.g. ".yaml")
	suffix string,
	// notBefore is the unix timestamp, all cache created before this timetamp is marked as expired
	notBefore int64,
) (active, expired, invalid []string, err error) {
	// find from local cache
	// ${DUKKHA_CACHE_DIR}/renderer-<rendererName>/<md5sum(key)>-<unix-timestamp>

	entries, err := fs.ReadDir(cacheDir, ".")
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			err = fmt.Errorf("checking local cache dir: %w", err)
			return
		}

		// no cache exists
		err = nil
		return
	}

	// ensure not working with entries
	// helps normalizing entry index rules for later processing
	if len(entries) == 0 {
		return
	}

	// entries are sorted per fs.ReadDirFS.ReadDir requirement
	// so we can do binary search directly
	start := sort.Search(len(entries), func(i int) bool {
		return prefix <= entries[i].Name()
	})

	if start == len(entries) {
		// not found
		return
	}

	// find last entry with same prefix
	// then we have a full list of cached data
	end := start
	for ; end < len(entries); end++ {
		if !strings.HasPrefix(entries[end].Name(), prefix) {
			break
		}
	}

	for _, info := range entries[start:end] {
		filename := info.Name()

		parts := strings.SplitN(strings.TrimSuffix(filename, suffix), "-", 2)
		if len(parts) != 2 || parts[0] != prefix {
			// invalid cache file
			invalid = append(invalid, filename)
			continue
		}

		timestamp, err2 := strconv.ParseInt(
			// trim padding
			strings.TrimLeft(parts[1], "0"),
			10, 64,
		)
		if err2 != nil {
			invalid = append(invalid, filename)
			continue
		}

		if timestamp < notBefore {
			expired = append(expired, filename)
			continue
		}

		active = append(active, filename)
	}

	return
}
