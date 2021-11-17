package transform

import (
	"bytes"
	"fmt"
	"os"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/third_party/golang/text/template"
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
	rc dukkha.RenderingContext, rawData interface{}, _ []dukkha.RendererAttribute,
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

	spec := rs.Init(&Spec{}, nil).(*Spec)
	err = yaml.Unmarshal(rawBytes, spec)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: failed to resolve input as transform spec: %w",
			d.name, err,
		)
	}

	// only resolve value and ops list, we need to resolve
	// each operation step by step with `value` injected
	err = spec.ResolveFields(rc, 2)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: failed to resolve transform spec fields: %w",
			d.name, err,
		)
	}

	data := []byte(spec.Value)
	for i, op := range spec.Ops {
		data, err = op.Do(rc, data)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to do #%d transformation: %w",
				d.name, i, err,
			)
		}
	}

	return data, nil
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

	Template *string   `yaml:"template,omitempty"`
	Shell    *string   `yaml:"shell,omitempty"`
	Checksum *Checksum `yaml:"checksum,omitempty"`
}

func (op *Operation) Do(_rc dukkha.RenderingContext, valueBytes []byte) ([]byte, error) {
	rc2 := _rc.(interface {
		dukkha.RenderingContext
		SetVALUE(s string)
		VALUE() string
	})

	rc2.SetVALUE(string(valueBytes))
	rc2.AddEnv(true, &dukkha.EnvEntry{
		Name:  "VALUE",
		Value: rc2.VALUE(),
	})

	// do not expose SetVALUE to template operation
	rc := rc2.(interface {
		dukkha.RenderingContext
		VALUE() string
	})

	err := op.ResolveFields(rc, -1)
	if err != nil {
		return nil, err
	}

	switch {
	case op.Template != nil:
		tplStr := *op.Template

		var tpl *template.Template
		tpl, err = templateutils.CreateTemplate(rc).Parse(tplStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %q: %w", tplStr, err)
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, rc)
		if err != nil {
			return nil, fmt.Errorf("failed to execute template %q: %w", tplStr, err)
		}

		return buf.Next(buf.Len()), nil
	case op.Shell != nil:
		script := *op.Shell

		buf := &bytes.Buffer{}
		var runner *interp.Runner
		runner, err = templateutils.CreateEmbeddedShellRunner(
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
			return nil, err
		}

		return buf.Next(buf.Len()), nil
	case op.Checksum != nil:
		return valueBytes, op.Checksum.Verify()
	default:
		return valueBytes, nil
	}
}
