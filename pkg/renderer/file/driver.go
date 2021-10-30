package file

import (
	"fmt"
	"os"

	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

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

	cache *renderer.Cache
}

func (d *Driver) Init(ctx dukkha.ConfigResolvingContext) error {
	if d.EnableCache {
		d.cache = renderer.NewCache(int64(d.CacheSizeLimit), d.CacheMaxAge)
	}

	return nil
}

func (d *Driver) RenderYaml(
	_ dukkha.RenderingContext, rawData interface{},
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

	var data []byte
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
