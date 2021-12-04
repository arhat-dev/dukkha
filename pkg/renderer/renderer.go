package renderer

import (
	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/pkg/fshelper"
	"arhat.dev/rs"
)

type BaseRenderer struct {
	rs.BaseField `yaml:"-"`

	RendererAlias string `yaml:"alias"`

	DefaultAttributes []dukkha.RendererAttribute `yaml:"attributes"`
}

func (r *BaseRenderer) Init(cacheFS *fshelper.OSFS) error {
	return nil
}

func (r *BaseRenderer) Alias() string {
	return r.RendererAlias
}

func (r *BaseRenderer) Attributes(override []dukkha.RendererAttribute) []dukkha.RendererAttribute {
	if len(override) == 0 {
		return r.DefaultAttributes
	}

	return override
}

type BaseInMemCachedRenderer struct {
	rs.BaseField `yaml:"-"`

	BaseRenderer `yaml:",inline"`

	CacheConfig CacheConfig `yaml:"cache"`

	// CacheFS provided when calling Init
	CacheFS *fshelper.OSFS `yaml:"-"`

	// Cache is the in memory cache, nil if not enabled in CacheConfig
	Cache *cache.Cache `yaml:"-"`
}

func (d *BaseInMemCachedRenderer) Init(cacheFS *fshelper.OSFS) error {
	d.CacheFS = cacheFS
	if d.CacheConfig.Enabled {
		d.Cache = cache.NewCache(
			int64(d.CacheConfig.MaxItemSize),
			int64(d.CacheConfig.Size),
			int64(d.CacheConfig.Timeout.Seconds()),
		)
	}

	return nil
}

type BaseTwoTierCachedRenderer struct {
	rs.BaseField `yaml:"-"`

	BaseRenderer `yaml:",inline"`

	CacheConfig CacheConfig `yaml:"cache"`

	// CacheFS provided when calling Init
	CacheFS *fshelper.OSFS `yaml:"-"`

	// Cache is always not nil after Init
	Cache *cache.TwoTierCache `yaml:"-"`
}

func (d *BaseTwoTierCachedRenderer) Init(cacheFS *fshelper.OSFS) error {
	d.CacheFS = cacheFS
	if d.CacheConfig.Enabled {
		d.Cache = cache.NewTwoTierCache(
			cacheFS,
			int64(d.CacheConfig.MaxItemSize),
			int64(d.CacheConfig.Size),
			int64(d.CacheConfig.Timeout.Seconds()),
		)
	} else {
		d.Cache = cache.NewTwoTierCache(cacheFS, 0, 0, -1)
	}

	return nil
}
