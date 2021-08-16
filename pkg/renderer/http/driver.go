package http

import (
	"fmt"
	"io"
	"net/http"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
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
	rs.BaseField

	renderer.CacheConfig `yaml:",inline"`

	DefaultConfig rendererHTTPConfig `yaml:",inline"`

	defaultClient *http.Client
	cache         *renderer.Cache
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	if d.EnableCache {
		d.cache = renderer.NewCache(int64(d.CacheSizeLimit), d.CacheMaxAge)
	}

	d.defaultClient = d.DefaultConfig.createClient()
	return nil
}

func (d *driver) RenderYaml(rc dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	var (
		url    string
		client *http.Client
	)

	switch t := rawData.(type) {
	case string:
		url = t
		client = d.defaultClient
	case []byte:
		url = string(t)
		client = d.defaultClient
	default:
		rawBytes, err := yamlhelper.ToYamlBytes(rawData)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: unexpected non yaml input: %w", DefaultName, err,
			)
		}

		cfg := rs.Init(&inputHTTPConfig{}, rc).(*inputHTTPConfig)
		err = yaml.Unmarshal(rawBytes, cfg)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to unmarshal input as config: %w", DefaultName, err,
			)
		}

		err = cfg.ResolveFields(rc, -1)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to resolve input config: %w",
				DefaultName, err,
			)
		}

		// config resolved

		url = cfg.URL
		client = cfg.Config.createClient()
	}

	var (
		data []byte
		err  error
	)

	if d.cache != nil {
		data, err = d.cache.Get(url,
			renderer.CreateRefreshFuncForRemote(
				renderer.FormatCacheDir(rc.CacheDir(), DefaultName),
				d.CacheMaxAge,
				func(key string) ([]byte, error) {
					return d.fetchRemote(client, key)
				},
			),
		)
	} else {
		data, err = d.fetchRemote(client, url)
	}

	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", DefaultName, err)
	}

	return data, err
}

func (d *driver) fetchRemote(client *http.Client, url string) ([]byte, error) {
	// TODO: support more http features
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
