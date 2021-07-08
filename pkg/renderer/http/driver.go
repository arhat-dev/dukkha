package http

import (
	"fmt"
	"io"
	"net/http"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

// nolint:revive
const (
	DefaultName = "http"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault() dukkha.Renderer {
	return &driver{CacheConfig: renderer.CacheConfig{EnableCache: true}}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	field.BaseField

	renderer.CacheConfig `yaml:",inline"`

	fetch renderer.CacheRefreshFunc
	cache *renderer.Cache
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	d.fetch = renderer.CreateFetchFunc(
		ctx.CacheDir(), DefaultName, d.CacheMaxAge,
		func(url string) ([]byte, error) {
			// TODO: support more http features
			resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}

			respBody, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				return nil, err
			}

			return respBody, nil
		},
	)

	if d.EnableCache {
		d.cache = renderer.NewCache(
			int64(d.CacheSizeLimit), d.CacheMaxAge, d.fetch,
		)
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
		data, err = d.cache.Get(path)
	} else {
		data, err = d.fetch(path)
	}

	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", DefaultName, err)
	}

	return data, err
}
