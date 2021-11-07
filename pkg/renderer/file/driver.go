package file

import (
	"fmt"
	"os"
	"strings"

	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/cache"
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
	return &Driver{
		name:        name,
		CacheConfig: renderer.CacheConfig{EnableCache: false},
	}
}

var _ dukkha.Renderer = (*Driver)(nil)

type Driver struct {
	rs.BaseField `yaml:"-"`
	name         string

	renderer.CacheConfig `yaml:",inline"`

	cache *cache.Cache
}

func (d *Driver) Init(ctx dukkha.ConfigResolvingContext) error {
	if d.EnableCache {
		d.cache = cache.NewCache(
			int64(d.CacheItemSizeLimit),
			int64(d.CacheSizeLimit),
			int64(d.CacheMaxAge.Seconds()),
		)
	}

	return nil
}

func (d *Driver) RenderYaml(
	_ dukkha.RenderingContext, rawData interface{}, _ []dukkha.RendererAttribute,
) ([]byte, error) {
	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	var path string
	switch t := rawData.(type) {
	case string:
		path = t
	case []byte:
		path = string(t)
	case *yaml.Node:
		path = t.Value
	default:
		return nil, fmt.Errorf(
			"renderer.%s: unexpected non-string input %T",
			d.name, rawData,
		)
	}

	path = strings.TrimSpace(path)

	var data []byte
	if d.cache != nil {
		data, err = d.cache.Get(
			cache.IdentifiableString(path),
			func(_ cache.IdentifiableObject) ([]byte, error) {
				return os.ReadFile(path)
			},
		)
	} else {
		data, err = os.ReadFile(path)
	}

	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", d.name, err)
	}

	return data, err
}
