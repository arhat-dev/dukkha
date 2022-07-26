package transform

import (
	"fmt"
	"strings"

	"arhat.dev/pkg/stringhelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"arhat.dev/tlang"
	"gopkg.in/yaml.v3"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/third_party/golang/text/template"
)

const (
	DefaultName = "T"
)

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

	data := spec.Value
	for i, op := range spec.Ops {
		data, err = op.Do(rc, data)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: step #%d: %w",
				d.name, i, err,
			)
		}
	}

	return stringhelper.ToBytes[byte, byte](data), nil
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

	TLang    *string   `yaml:"tlang,omitempty"`
	Template *string   `yaml:"template,omitempty"`
	Shell    *string   `yaml:"shell,omitempty"`
	Checksum *Checksum `yaml:"checksum,omitempty"`
}

func (op *Operation) Do(_rc dukkha.RenderingContext, value string) (_ string, err error) {
	rc2 := _rc.(interface {
		dukkha.RenderingContext
		di.VALUEGetter
		di.VALUESetter
	})

	rc2.SetVALUE(value)
	rc2.AddEnv(true, &dukkha.NameValueEntry{
		Name:  "VALUE",
		Value: rc2.VALUE().(string),
	})

	// do not expose SetVALUE to template operation
	rc := rc2.(interface {
		dukkha.RenderingContext
		di.VALUEGetter
	})

	err = op.ResolveFields(rc, -1)
	if err != nil {
		return
	}

	switch {
	case op.TLang != nil:
		var (
			script = *op.TLang
			parsed *tlang.Template
			buf    strings.Builder
		)

		parsed, err = templateutils.CreateTLangTemplate(rc).Parse(script)
		if err != nil {
			err = fmt.Errorf("parse tlang %q: %w", script, err)
			return
		}

		err = parsed.Execute(&buf, rc)
		if err != nil {
			err = fmt.Errorf("executing tlang %q: %w", script, err)
			return
		}

		return buf.String(), nil
	case op.Template != nil:
		var (
			tmpl   = *op.Template
			parsed *template.Template
			buf    strings.Builder
		)

		parsed, err = templateutils.CreateTextTemplate(rc).Parse(tmpl)
		if err != nil {
			err = fmt.Errorf("parsing template %q: %w", tmpl, err)
			return
		}

		err = parsed.Execute(&buf, rc)
		if err != nil {
			err = fmt.Errorf("executing template %q: %w", tmpl, err)
			return
		}

		return buf.String(), nil
	case op.Shell != nil:
		var (
			buf    strings.Builder
			runner *interp.Runner
		)
		runner, err = templateutils.CreateShellRunner(
			rc.WorkDir(), rc, nil, &buf, rc.Stderr(),
		)
		if err != nil {
			err = fmt.Errorf("create embedded shell: %w", err)
			return
		}

		parser := syntax.NewParser(
			syntax.Variant(syntax.LangBash),
		)

		err = templateutils.RunScript(rc, runner, parser, *op.Shell)
		if err != nil {
			err = fmt.Errorf("run script: %w", err)
			return
		}

		return buf.String(), nil
	case op.Checksum != nil:
		return value, op.Checksum.Verify(rc.FS())
	default:
		return value, nil
	}
}
