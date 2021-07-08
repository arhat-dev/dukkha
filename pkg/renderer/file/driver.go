package file

import (
	"fmt"
	"os"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

// nolint:revive
const (
	DefaultName = "file"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault() dukkha.Renderer {
	return &driver{CacheConfig: renderer.CacheConfig{EnableCache: false}}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	field.BaseField

	renderer.CacheConfig `yaml:",inline"`

	cache *renderer.Cache
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	if d.EnableCache {
		d.cache = renderer.NewCache(int64(d.CacheSizeLimit), d.CacheMaxAge)
	}

	ctx.AddRenderer(DefaultName, d)
	return nil
}

func (d *driver) RenderYaml(_ dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	path, ok := rawData.(string)
	if !ok {
		return nil, fmt.Errorf("renderer.%s: unexpected non-string input %T", DefaultName, rawData)
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
		return nil, fmt.Errorf("renderer.%s: %w", DefaultName, err)
	}

	return data, err
}
