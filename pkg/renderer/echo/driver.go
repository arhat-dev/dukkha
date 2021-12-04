package echo

import (
	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
)

const (
	DefaultName = "echo"
)

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

func NewDefault(name string) dukkha.Renderer { return &Driver{name: name} }

var _ dukkha.Renderer = (*Driver)(nil)

type Driver struct {
	rs.BaseField `yaml:"-"`

	RendererAlias string `yaml:"alias"`

	name string
}

func (d *Driver) Alias() string { return d.RendererAlias }

func (d *Driver) Init(cacheFS *fshelper.OSFS) error { return nil }

func (d *Driver) RenderYaml(
	_ dukkha.RenderingContext, rawData interface{}, _ []dukkha.RendererAttribute,
) ([]byte, error) {
	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	return yamlhelper.ToYamlBytes(rawData)
}
