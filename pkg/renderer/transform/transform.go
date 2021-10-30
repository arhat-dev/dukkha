package transform

import (
	"bytes"
	"fmt"
	"os"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/templateutils"
)

const DefaultName = "transform"

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

func NewDefault(name string) dukkha.Renderer { return &Driver{name: name} }

var _ dukkha.Renderer = (*Driver)(nil)

type Driver struct {
	rs.BaseField `yaml:"-"`

	name string
}

func (d *Driver) Init(ctx dukkha.ConfigResolvingContext) error { return nil }

func (d *Driver) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{},
) ([]byte, error) {
	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	rawBytes, err := yamlhelper.ToYamlBytes(rawData)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: unsupported input type %T: %w",
			d.name, rawData, err,
		)
	}

	spec := rs.Init(&Spec{}, &rs.Options{
		InterfaceTypeHandler: rc,
	}).(*Spec)
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
				"renderer.%s: failed to do #%d transformation: %w",
				d.name, i, err,
			)
		}
	}

	switch dt := data.(type) {
	case string:
		return []byte(dt), nil
	case []byte:
		return dt, nil
	default:
		return yaml.Marshal(data)
	}
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
	Shell    *string `yaml:"shell,omitempty"`
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
	case op.Shell != nil:
		script := *op.Shell

		valueBytes, err := yamlhelper.ToYamlBytes(data)
		if err != nil {
			return nil, err
		}

		rc.AddEnv(true, &dukkha.EnvEntry{
			Name:  "VALUE",
			Value: string(valueBytes),
		})

		buf := &bytes.Buffer{}
		runner, err := templateutils.CreateEmbeddedShellRunner(
			rc.WorkingDir(), rc, nil, buf, os.Stderr,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to create embedded shell: %w",
				err,
			)
		}

		parser := syntax.NewParser(
			syntax.Variant(syntax.LangBash),
		)

		err = templateutils.RunScriptInEmbeddedShell(rc, runner, parser, script)
		if err != nil {
			return nil, fmt.Errorf("failed to run shell script: %w", err)
		}

		return string(buf.Next(buf.Len())), nil
	default:
		// TODO: shall we consider nop as an error?
		return data, nil
	}
}
