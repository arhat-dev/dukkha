package transform

import (
	"bytes"
	"fmt"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/templateutils"
)

const DefaultName = "transform"

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

func NewDefault(name string) dukkha.Renderer { return &driver{name: name} }

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	rs.BaseField `yaml:"-"`

	name string
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error { return nil }

func (d *driver) RenderYaml(rc dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	rawBytes, err := yamlhelper.ToYamlBytes(rawData)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: unsupported input type %T: %w",
			d.name, rawData, err,
		)
	}

	spec := rs.Init(&Spec{}, rc).(*Spec)
	err = yaml.Unmarshal(rawBytes, spec)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: failed to resolve input as transform spec: %w",
			d.name, err,
		)
	}

	err = spec.ResolveFields(rc, -1)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: failed to resolve transform spec fields: %w",
			d.name, err,
		)
	}

	var data interface{} = spec.Value
	for i, op := range spec.Ops {
		data, err = op.Do(rc, data)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to do #%d transform operation: %w",
				d.name, i, err,
			)
		}
	}

	ret, err := yamlhelper.ToYamlBytes(data)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: failed to marshal result as yaml: %w",
			d.name, err,
		)
	}

	return ret, nil
}

// Spec for yaml data transformation
type Spec struct {
	rs.BaseField `yaml:"-"`

	// Value always be string, so you can decide which format is converts to in operations
	Value string `yaml:"value"`

	// Ops the transform operations to run
	Ops []*Operation `yaml:"ops"`
}

type Operation struct {
	rs.BaseField `yaml:"-"`

	Template *string `yaml:"template,omitempty"`
}

type tplDataType struct {
	dukkha.RenderingContext
	Value interface{}
}

func (op *Operation) Do(rc dukkha.RenderingContext, data interface{}) (interface{}, error) {
	switch {
	case op.Template != nil:
		tplStr := *op.Template

		tpl, err := templateutils.CreateTemplate(rc).Parse(tplStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %q: %w", tplStr, err)
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, &tplDataType{
			RenderingContext: rc,
			Value:            data,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to execute template %q: %w", tplStr, err)
		}

		return string(buf.Next(buf.Len())), nil
	default:
		// TODO: shall we consider nop as an error?
		return data, nil
	}
}
