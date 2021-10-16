package s3

import (
	"fmt"

	"arhat.dev/pkg/rshelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
)

// nolint:revive
const (
	DefaultName = "s3"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault(name string) dukkha.Renderer {
	return &driver{
		name:        name,
		CacheConfig: renderer.CacheConfig{EnableCache: false},
	}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	rs.BaseField `yaml:"-"`
	name         string

	renderer.CacheConfig `yaml:",inline"`

	DefaultConfig rendererS3Config `yaml:",inline"`

	defaultClient *s3Client

	cache *renderer.Cache
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
	rc dukkha.RenderingContext, rawData interface{},
) (interface{}, error) {
	var (
		path   string
		client *s3Client
	)

	switch t := rawData.(type) {
	case string:
		path = t
		client = d.defaultClient
	case []byte:
		path = string(t)
		client = d.defaultClient
	default:
		rawBytes, err := yamlhelper.ToYamlBytes(rawData)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: unexpected non yaml input: %w",
				d.name, err,
			)
		}

		spec := rshelper.InitAll(&inputS3Sepc{}, &rs.Options{
			InterfaceTypeHandler: rc,
		}).(*inputS3Sepc)
		err = yaml.Unmarshal(rawBytes, spec)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to unmarshal input as s3 spec: %w",
				d.name, err,
			)
		}

		err = spec.ResolveFields(rc, -1)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to resolve s3 spec: %w",
				d.name, err,
			)
		}

		// config resolved

		path = spec.Path
		client, err = spec.Config.createClient()
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to create s3 client for spec: %w",
				d.name, err,
			)
		}
	}

	var (
		data []byte
		err  error
	)

	if d.cache != nil {
		data, err = d.cache.Get(path,
			renderer.CreateRefreshFuncForRemote(
				renderer.FormatCacheDir(rc.CacheDir(), d.name),
				d.CacheMaxAge,
				func(key string) ([]byte, error) {
					return client.download(rc, path)
				},
			),
		)
	} else {
		data, err = client.download(rc, path)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: %w",
			d.name, err,
		)
	}

	return data, err
}
