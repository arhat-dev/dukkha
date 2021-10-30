package template

import (
	"bytes"
	"fmt"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/templateutils"
)

const DefaultName = "template"

func init() {
	dukkha.RegisterRenderer(
		DefaultName,
		NewDefault,
	)
}

func NewDefault(name string) dukkha.Renderer {
	return &Driver{
		name: name,
	}
}

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

	tplBytes, err := yamlhelper.ToYamlBytes(rawData)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: unsupported input type %T: %w",
			d.name, rawData, err,
		)
	}

	tplStr := string(tplBytes)
	tpl, err := templateutils.CreateTemplate(rc).Parse(tplStr)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: failed to parse template \n\n%s\n\n %w",
			d.name, tplStr, err,
		)
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, rc)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: failed to execute template \n\n%s\n\n %w",
			d.name, tplStr, err,
		)
	}

	return buf.Next(buf.Len()), nil
}
