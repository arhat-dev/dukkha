package tpl

import (
	"bytes"
	"fmt"
	"os"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/templateutils"
)

const (
	DefaultName = "tpl"

	AttrUseSpec = "use-spec"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
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

	Options configSpec `yaml:",inline"`

	variables map[string]interface{}
}

func (d *Driver) Init(ctx dukkha.ConfigResolvingContext) error {
	d.variables = d.Options.Variables.NormalizedValue()
	return nil
}

func (d *Driver) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{},
	attributes []dukkha.RendererAttribute,
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

	var (
		useSpec bool
	)
	for _, attr := range attributes {
		switch attr {
		case AttrUseSpec:
			useSpec = true
		default:
		}
	}

	var (
		include   []string
		variables map[string]interface{}
		tplStr    string
	)

	if useSpec {
		var spec *inputSpec
		spec, err = resolveInputSpec(rc, tplBytes)
		if err != nil {
			return nil, fmt.Errorf("renderer.%s: %s", d.name, err)
		}

		tplStr = spec.Template
		include = spec.Config.Include
		variables = spec.Config.Variables.NormalizedValue()
	} else {
		tplStr = string(tplBytes)
		include = d.Options.Include
		variables = d.variables
	}

	data, err := renderTemplate(rc, include, variables, tplStr)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", d.name, err)
	}

	return data, nil
}

func resolveInputSpec(rc dukkha.RenderingContext, tplBytes []byte) (*inputSpec, error) {
	spec := rs.Init(&inputSpec{}, nil).(*inputSpec)

	err := yaml.Unmarshal(tplBytes, spec)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid template input spec: %w", err,
		)
	}

	err = spec.ResolveFields(rc, -1)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve template input spec: %w", err,
		)
	}

	return spec, nil
}

func renderTemplate(
	rc dukkha.RenderingContext,
	inc []string,
	variables map[string]interface{},

	tplStr string,
) ([]byte, error) {
	_fs := afero.NewIOFS(afero.NewOsFs())

	var include []string
	for _, inc := range inc {
		matches, err := doublestar.Glob(_fs, inc)
		if err != nil {
			_, err2 := os.Stat(inc)
			if err2 != nil {
				return nil, err
			}

			include = append(include, inc)
		} else {
			include = append(include, matches...)
		}
	}

	tpl := templateutils.CreateTemplate(rc).
		Funcs(map[string]interface{}{
			"var": func() map[string]interface{} {
				return variables
			},
		})

	if len(include) != 0 {
		var err error
		tpl, err = tpl.ParseFiles(include...)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to load included files: %w", err,
			)
		}
	}

	tpl, err := tpl.Parse(tplStr)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse template: %w", err,
		)
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, rc)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to execute template \n\n%s\n\n %w",
			tplStr, err,
		)
	}

	return buf.Next(buf.Len()), nil
}
