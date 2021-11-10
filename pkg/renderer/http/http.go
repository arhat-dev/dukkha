package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"arhat.dev/pkg/rshelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
)

// nolint:revive
const (
	DefaultName = "http"
)

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

func NewDefault(name string) dukkha.Renderer {
	return &Driver{
		name: name,
		Config: cache.Config{
			EnableCache: true,
		},
	}
}

var _ dukkha.Renderer = (*Driver)(nil)

type Driver struct {
	rs.BaseField `yaml:"-"`
	name         string

	cache.Config `yaml:",inline"`

	DefaultConfig rendererHTTPConfig `yaml:",inline"`

	defaultClient *http.Client
	cache         *cache.TwoTierCache
}

func (d *Driver) Init(ctx dukkha.ConfigResolvingContext) error {
	dir := ctx.RendererCacheDir(d.name)
	if d.EnableCache {
		d.cache = cache.NewTwoTierCache(
			dir,
			int64(d.CacheItemSizeLimit),
			int64(d.CacheSizeLimit),
			int64(d.CacheMaxAge.Seconds()),
		)
	} else {
		d.cache = cache.NewTwoTierCache(dir, 0, 0, 0)
	}

	var err error
	d.defaultClient, err = d.DefaultConfig.createClient()
	return err
}

func (d *Driver) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{}, attributes []dukkha.RendererAttribute,
) ([]byte, error) {
	var (
		reqURL string
		client *http.Client
		config *rendererHTTPConfig
	)

	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	switch t := rawData.(type) {
	case string:
		reqURL = t
		client = d.defaultClient
		config = &d.DefaultConfig
	case []byte:
		reqURL = string(t)
		client = d.defaultClient
		config = &d.DefaultConfig
	default:
		var rawBytes []byte
		rawBytes, err = yamlhelper.ToYamlBytes(rawData)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: unexpected non yaml input: %w",
				d.name, err,
			)
		}

		spec := rshelper.InitAll(&inputHTTPSpec{}, &rs.Options{
			InterfaceTypeHandler: rc,
		}).(*inputHTTPSpec)
		err = yaml.Unmarshal(rawBytes, spec)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to unmarshal input as config: %w",
				d.name, err,
			)
		}

		err = spec.ResolveFields(rc, -1)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to resolve input config: %w",
				d.name, err,
			)
		}

		// config resolved

		reqURL = spec.URL
		client, err = spec.Config.createClient()
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to create http client for spec: %w",
				d.name, err,
			)
		}

		config = &spec.Config
	}

	data, err := renderer.HandleRenderingRequestWithRemoteFetch(
		d.cache,
		cache.IdentifiableString(reqURL),
		func(_ cache.IdentifiableObject) (io.ReadCloser, error) {
			return d.fetchRemote(client, reqURL, config)
		},
		attributes,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s failed to fetch http content: %w",
			d.name, err,
		)
	}

	return data, err
}

func (d *Driver) fetchRemote(
	client *http.Client,
	url string,
	config *rendererHTTPConfig,
) (io.ReadCloser, error) {
	var (
		req *http.Request
		err error
	)

	var body io.Reader
	if config.Body != nil {
		body = strings.NewReader(*config.Body)
	}

	method := strings.ToUpper(config.Method)
	if len(method) == 0 {
		method = http.MethodGet
	}

	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	if len(config.User) != 0 {
		req.SetBasicAuth(config.User, config.Password)
	}

	seen := make(map[string]struct{})
	for _, h := range config.Headers {
		_, ok := seen[h.Name]
		if !ok {
			seen[h.Name] = struct{}{}
			req.Header.Set(h.Name, h.Value)
		} else {
			req.Header.Add(h.Name, h.Value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
