package file

import (
	"fmt"
	"strings"

	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
)

const (
	DefaultName = "file"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault(name string) dukkha.Renderer {
	return &Driver{
		name:        name,
		CacheConfig: renderer.CacheConfig{Enabled: false},
	}
}

var _ dukkha.Renderer = (*Driver)(nil)

type Driver struct {
	rs.BaseField `yaml:"-"`

	name string

	CacheConfig renderer.CacheConfig `yaml:"cache"`

	cache *cache.Cache
}

func (d *Driver) Init(ctx dukkha.ConfigResolvingContext) error {
	if d.CacheConfig.Enabled {
		d.cache = cache.NewCache(
			int64(d.CacheConfig.MaxItemSize),
			int64(d.CacheConfig.Size),
			int64(d.CacheConfig.Timeout.Seconds()),
		)
	}

	return nil
}

func (d *Driver) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{}, _ []dukkha.RendererAttribute,
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
				return rc.FS().ReadFile(path)
			},
		)
	} else {
		data, err = rc.FS().ReadFile(path)
	}

	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", d.name, err)
	}

	return data, err
}
