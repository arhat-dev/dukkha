package env

import (
	"fmt"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/templateutils"
)

// nolint:revive
const (
	DefaultName = "env"
)

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

func NewDefault(name string) dukkha.Renderer { return &Driver{name: name} }

var _ dukkha.Renderer = (*Driver)(nil)

type Driver struct {
	rs.BaseField `yaml:"-"`

	name string
}

func (d *Driver) Init(ctx dukkha.ConfigResolvingContext) error {
	return nil
}

func (d *Driver) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{},
) ([]byte, error) {
	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	bytesToExpand, err := yamlhelper.ToYamlBytes(rawData)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: unsupported input type %T: %w",
			d.name, rawData, err,
		)
	}

	ret, err := templateutils.ExpandEnv(rc, string(bytesToExpand))
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", d.name, err)
	}

	return []byte(ret), nil
}
