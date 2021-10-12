package file

import (
	"fmt"
	"os"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
)

// nolint:revive
const (
	DefaultName = "file"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault(name string) dukkha.Renderer {
	return &driver{
		name:        name,
		CacheConfig: renderer.CacheConfig{EnableCache: false},
	}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	rs.BaseField `yaml:"-"`
	name         string

	renderer.CacheConfig `yaml:",inline"`

	cache *renderer.Cache
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	if d.EnableCache {
		d.cache = renderer.NewCache(int64(d.CacheSizeLimit), d.CacheMaxAge)
	}

	return nil
}

func (d *driver) RenderYaml(_ dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	var path string
	switch t := rawData.(type) {
	case string:
		path = t
	case []byte:
		path = string(t)
	default:
		return nil, fmt.Errorf(
			"renderer.%s: unexpected non-string input %T",
			d.name, rawData,
		)
	}

	var (
		data []byte
		err  error
	)

	if d.cache != nil {
		data, err = d.cache.Get(path, os.ReadFile)
	} else {
		data, err = os.ReadFile(path)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: %w",
			d.name, err,
		)
	}

	return data, err
}
