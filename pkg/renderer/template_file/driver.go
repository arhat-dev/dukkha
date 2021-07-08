package template_file

import (
	"fmt"
	"os"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/template"
)

const DefaultName = "template_file"

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault() dukkha.Renderer {
	return &driver{
		impl:        template.NewDefault(),
		CacheConfig: renderer.CacheConfig{EnableCache: false},
	}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	field.BaseField

	impl dukkha.Renderer

	renderer.CacheConfig `yaml:",inline"`

	cache *renderer.Cache
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	if d.EnableCache {
		d.cache = renderer.NewCache(
			int64(d.CacheSizeLimit), d.CacheMaxAge, os.ReadFile,
		)
	}

	ctx.AddRenderer(DefaultName, d)
	return nil
}

func (d *driver) RenderYaml(rc dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	path, ok := rawData.(string)
	if !ok {
		return nil, fmt.Errorf("renderer.%s: unexpected non string input %T", DefaultName, rawData)
	}

	var (
		tplBytes []byte
		err      error
	)

	if d.cache != nil {
		tplBytes, err = d.cache.Get(path)
	} else {
		tplBytes, err = os.ReadFile(path)
	}

	if err != nil {
		return nil, fmt.Errorf("renderer.%s: failed to read template file: %w", DefaultName, err)
	}

	result, err := d.impl.RenderYaml(rc, tplBytes)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: failed to render file %q: %w", DefaultName, path, err)
	}

	return result, nil
}
