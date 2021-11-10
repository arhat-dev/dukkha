package cache

import (
	"time"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/utils"
)

// Config is the config for rendered data caching
type Config struct {
	rs.BaseField `yaml:"-"`

	// EnableCache activates rendered data caching
	//
	// * for renderers reading data directly from local disk (e.g. file, archivefile):
	//     will cache content in memory with size limit applied
	// * for renderers doing remote fetch (e.g. http, git):
	//     will cache data on local disk first (cache size limiting is not effective at this time)
	// 	   then cache data in memory with size limit applied
	//
	// Defaults to `false`
	EnableCache bool `yaml:"enable_cache"`

	// CacheItemSizeLimit is the maximum size limit an item can be cached in memory
	//
	// Format: <number><unit>
	// 	where unit can be one of: [ , B, KB, MB, GB, TB, PB]
	//
	// Defaults to `0` (no size limit for single item)
	CacheItemSizeLimit utils.Size `yaml:"cache_item_size_limit"`

	// CacheSizeLimit limits maximum in memory size of cache of the renderer
	//
	// Format: <number><unit>
	// 	where unit can be one of: [ , B, KB, MB, GB, TB, PB]
	//
	// Defaults to `0` (no size limit)
	CacheSizeLimit utils.Size `yaml:"cache_size_limit"`

	// CacheMaxAge limits maximum data caching time
	//
	// if caching is enabled and this option is set to 0:
	//  in memory cache will never expire during runtime
	// 	file cache for remote content will expire immediately (probably that's not what you want)
	//
	// Defaults to `0`
	CacheMaxAge time.Duration `yaml:"cache_max_age"`
}
