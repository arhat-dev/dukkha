// Package archivefile provides a renderer generating value by extracting
// file content from archive directly
package archivefile

import (
	"fmt"
	"strings"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
)

// nolint:revive
const (
	DefaultName = "archivefile"
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
	} else {
		d.cache = cache.NewCache(0, 0, -1)
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

	var (
		onlineSpec string
		spec       *inputSpec
	)

	switch t := rawData.(type) {
	case string:
		onlineSpec = t
	case []byte:
		onlineSpec = string(t)
	default:
		var rawDataBytes []byte
		rawDataBytes, err = yamlhelper.ToYamlBytes(t)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", d.name, err)
		}

		spec := rs.Init(&inputSpec{}, nil).(*inputSpec)
		err = yaml.Unmarshal(rawDataBytes, spec)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid input spec: %w", d.name, err)
		}

		err = spec.ResolveFields(rc, -1)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to resolve input spec: %w", d.name, err)
		}

		return nil, fmt.Errorf(
			"%s: unexpected non-string input %T", d.name, t,
		)
	}

	if spec == nil {
		spec, err = convertPathToSpec(onlineSpec)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid oneline spec %w", d.name, err)
		}
	}

	data, err := d.cache.Get(spec, extractFileFromArchive)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", d.name, err)
	}

	return data, err
}

// nolint:unparam
func convertPathToSpec(onelineSpec string) (*inputSpec, error) {
	onelineSpec = strings.TrimSpace(onelineSpec)
	_ = onelineSpec
	return nil, nil
}

func extractFileFromArchive(obj cache.IdentifiableObject) ([]byte, error) {
	spec := obj.(*inputSpec)
	_ = spec
	return nil, nil
}
