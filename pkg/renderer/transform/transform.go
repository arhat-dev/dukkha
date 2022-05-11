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

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/third_party/golang/text/template"
)

const DefaultName = "T"

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

func NewDefault(name string) dukkha.Renderer { return &Driver{name: name} }

type Driver struct {
	rs.BaseField `yaml:"-"`

	renderer.BaseRenderer `yaml:",inline"`

	name string
}

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
			"renderer.%s: unmarshal transform spec: %w",
			d.name, err,
		)
	}

	// only resolve value and ops list, we need to resolve
	// each operation step by step with `value` injected
	err = spec.ResolveFields(rc, 2)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: resolving transform spec: %w",
			d.name, err,
		)
	}

	data := []byte(spec.Value)
	for i, op := range spec.Ops {
		data, err = op.Do(rc, data)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: step #%d: %w",
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
		di.VALUEGetter
		di.VALUESetter
	})

	rc2.SetVALUE(string(valueBytes))
	rc2.AddEnv(true, &dukkha.EnvEntry{
		Name:  "VALUE",
		Value: rc2.VALUE().(string),
	})

	// do not expose SetVALUE to template operation
	rc := rc2.(interface {
		dukkha.RenderingContext
		di.VALUEGetter
	})

	err := op.ResolveFields(rc, -1)
	if err != nil {
		return nil, err
	}

	switch {
	case op.Template != nil:
		var (
			tplStr = *op.Template
			tpl    *template.Template
			buf    bytes.Buffer
		)

		tpl, err = templateutils.CreateTemplate(rc).Parse(tplStr)
		if err != nil {
			return nil, fmt.Errorf("parsing template %q: %w", tplStr, err)
		}

		err = tpl.Execute(&buf, rc)
		if err != nil {
			return nil, fmt.Errorf("executing template %q: %w", tplStr, err)
		}

		return buf.Next(buf.Len()), nil
	case op.Shell != nil:
		var (
			buf    bytes.Buffer
			runner *interp.Runner
		)
		runner, err = templateutils.CreateEmbeddedShellRunner(
			rc.WorkDir(), rc, nil, &buf, os.Stderr,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"creating embedded shell: %w",
				err,
			)
		}

		parser := syntax.NewParser(
			syntax.Variant(syntax.LangBash),
		)

		err = templateutils.RunScriptInEmbeddedShell(rc, runner, parser, *op.Shell)
		if err != nil {
			return nil, err
		}

		return buf.Next(buf.Len()), nil
	case op.Checksum != nil:
		return valueBytes, op.Checksum.Verify(rc.FS())
	default:
		return valueBytes, nil
	}
}
