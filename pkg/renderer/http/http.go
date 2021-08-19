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

func NewDefault(name string) dukkha.Renderer {
	if len(name) != 0 {
		name = DefaultName + ":" + name
	} else {
		name = DefaultName
	}

	return &driver{
		name: name,
		CacheConfig: renderer.CacheConfig{
			EnableCache: true,
		},
	}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	rs.BaseField
	name string

	renderer.CacheConfig `yaml:",inline"`

	DefaultConfig rendererHTTPConfig `yaml:",inline"`

	defaultClient *http.Client
	cache         *renderer.Cache
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	if d.EnableCache {
		d.cache = renderer.NewCache(int64(d.CacheSizeLimit), d.CacheMaxAge)
	}

	var err error
	d.defaultClient, err = d.DefaultConfig.createClient()
	return err
}

func (d *driver) RenderYaml(
	rc dukkha.RenderingContext,
	rawData interface{},
) ([]byte, error) {
	var (
		reqURL string
		client *http.Client
		config *rendererHTTPConfig
	)

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
		rawBytes, err := yamlhelper.ToYamlBytes(rawData)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: unexpected non yaml input: %w",
				d.name, err,
			)
		}

		spec := rshelper.InitAll(&inputHTTPSpec{}, rc).(*inputHTTPSpec)
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

	var (
		data []byte
		err  error
	)

	if d.cache != nil {
		data, err = d.cache.Get(reqURL,
			renderer.CreateRefreshFuncForRemote(
				renderer.FormatCacheDir(rc.CacheDir(), DefaultName),
				d.CacheMaxAge,
				func(key string) ([]byte, error) {
					// key is the url we passed in
					return d.fetchRemote(client, key, config)
				},
			),
		)
	} else {
		data, err = d.fetchRemote(client, reqURL, config)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s failed to fetch http content: %w",
			d.name, err,
		)
	}

	return data, err
}

func (d *driver) fetchRemote(
	client *http.Client,
	url string,
	config *rendererHTTPConfig,
) ([]byte, error) {
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

	respBody, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
