package renderer

import (
	"time"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/utils"
)

// CacheConfig is the config for data caching
type CacheConfig struct {
	rs.BaseField `yaml:"-"`

	// Enabled activates data caching
	//
	// * for renderers reading data directly from local disk (e.g. file):
	//     will cache content in memory with size limit applied
	// * for renderers doing remote fetch (e.g. http, git, af):
	//     will cache data on local disk first (cache size limiting is not effective at this time)
	// 	   then cache data in memory with size limit applied
	//
	// Defaults to `false`
	Enabled bool `yaml:"enabled"`

	// MaxItemSize is the maximum size limit an item can be cached in memory
	//
	// Format: <number><unit>
	// 	where unit can be one of: [ , B, KB, MB, GB, TB, PB]
	//
	// Defaults to `0` (no size limit for single item)
	MaxItemSize utils.Size `yaml:"max_item_size"`

	// Size limits maximum in memory size of cached content
	//
	// Format: <number><unit>
	// 	where unit can be one of: [ , B, KB, MB, GB, TB, PB]
	//
	// Defaults to `0` (no size limit)
	Size utils.Size `yaml:"size"`

	// Timeout is the data caching duration
	//
	// if caching is enabled and this option is set to 0:
	//  in memory cache will never expire during runtime
	// 	file cache for remote content will expire immediately (probably that's not what you want)
	//
	// Defaults to `0`
	Timeout time.Duration `yaml:"timeout"`
}
