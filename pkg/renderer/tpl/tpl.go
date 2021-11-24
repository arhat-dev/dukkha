package tpl

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/third_party/golang/text/template"
)

const (
	DefaultName = "tpl"
)

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

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
		case renderer.AttrUseSpec:
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

	tpl := templateutils.CreateTemplate(rc)

	// arrange included template files in order
	// so we can include templates without {{ define "name" }} block
	// by filename and index in include

	definedTemplates := make(map[string]struct{})
	var tplList []*template.Template

	for _, inc := range include {
		// TODO: cache template files in memory
		// 	     maybe also parsed templates if we are sure rendering context
		// 	     is handled correctly

		tplBytes, err := os.ReadFile(inc)
		if err != nil {
			return nil, fmt.Errorf("failed to load template file: %q", err)
		}

		name := filepath.Base(inc)
		incTpl, err := tpl.New(name).Parse(string(tplBytes))
		if err != nil {
			return nil, fmt.Errorf("invalid template %q: %w", inc, err)
		}

		tplList = append(tplList, incTpl)

		definedTemplates[name] = struct{}{}
	}

	tplListSize := int64(len(tplList))

	for _, v := range tpl.Templates() {
		definedTemplates[v.Name()] = struct{}{}
	}

	// prevent infinite loop in template include
	const maxIncludeCount = 1000
	includedCount := make(map[string]int)

	// parse template entrypoint

	tpl, err := tpl.Parse(tplStr)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse template: %w", err,
		)
	}

	// Override placeholder funcs immediately before execute template
	tpl.Funcs(map[string]interface{}{
		"var": func() map[string]interface{} {
			return variables
		},

		// include like helm include
		"include": func(name string, data interface{}) (string, error) {
			count, ok := includedCount[name]
			if ok {
				if count >= maxIncludeCount {
					return "", fmt.Errorf("too many include of %q", name)
				}

				includedCount[name] = count + 1
			} else {
				includedCount[name] = 1
			}

			var (
				buf  strings.Builder
				err2 error
			)
			if _, defined := definedTemplates[name]; defined {
				err2 = tpl.ExecuteTemplate(&buf, name, data)
			} else {
				var idx int64
				idx, err2 = strconv.ParseInt(name, 10, 64)
				if err2 != nil {
					return "", fmt.Errorf("template %q undefined", name)
				}

				if idx < 0 {
					idx = tplListSize + idx
				}

				if idx < 0 || idx >= tplListSize {
					return "", fmt.Errorf(
						"invalid index out of range: %d not in [0,%d)",
						idx, tplListSize,
					)
				}

				err2 = tplList[idx].Execute(&buf, data)
			}

			includedCount[name]--

			if err2 != nil {
				return "", err2
			}

			return buf.String(), nil
		},
	})

	var buf bytes.Buffer
	err = tpl.Execute(&buf, rc)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: %s",
			err, tplStr,
		)
	}

	return buf.Next(buf.Len()), nil
}