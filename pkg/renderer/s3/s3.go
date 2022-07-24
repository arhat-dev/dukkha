package s3

import (
	"fmt"
	"io"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
)

const (
	DefaultName = "s3"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault(name string) dukkha.Renderer {
	return &Driver{name: name}
}

var _ dukkha.Renderer = (*Driver)(nil)

type Driver struct {
	rs.BaseField `yaml:"-"`

	renderer.BaseTwoTierCachedRenderer `yaml:",inline"`

	name string

	DefaultConfig rendererS3Config `yaml:",inline"`

	defaultClient *s3Client
}

func (d *Driver) Init(cacheFS *fshelper.OSFS) error {
	err := d.BaseTwoTierCachedRenderer.Init(cacheFS)
	if err != nil {
		return err
	}

	d.defaultClient, err = d.DefaultConfig.createClient()
	return err
}

func (d *Driver) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{}, attributes []dukkha.RendererAttribute,
) ([]byte, error) {
	var (
		path   string
		client *s3Client
	)

	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	switch t := rawData.(type) {
	case string:
		path = t
		client = d.defaultClient
	case []byte:
		path = string(t)
		client = d.defaultClient
	default:
		var rawBytes []byte
		rawBytes, err = yamlhelper.ToYamlBytes(rawData)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: unexpected non yaml input: %w",
				d.name, err,
			)
		}

		var (
			spec inputS3Sepc
			opts rs.Options
		)
		opts.InterfaceTypeHandler = rc
		rs.Init(&spec, &opts)

		err = yaml.Unmarshal(rawBytes, &spec)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: unmarshal input spec: %w",
				d.name, err,
			)
		}

		err = spec.ResolveFields(rc, -1)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: resolving input spec: %w",
				d.name, err,
			)
		}

		// config resolved

		path = spec.Path
		client, err = spec.Config.createClient()
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: creating s3 client for spec: %w",
				d.name, err,
			)
		}
	}

	data, err := renderer.HandleRenderingRequestWithRemoteFetch(
		d.Cache,
		cache.IdentifiableString(path),
		func(key cache.IdentifiableObject) (io.ReadCloser, error) {
			return client.download(rc, key.ScopeUniqueID())
		},
		d.Attributes(attributes),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: %w",
			d.name, err,
		)
	}

	return data, err
}
