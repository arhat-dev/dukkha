package tpl

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
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
		include   []*includeSpec
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
	inc []*includeSpec,
	variables map[string]interface{},

	tplStr string,
) ([]byte, error) {
	var (
		includeFiles []string
		includeText  []string
	)
	for _, inc := range inc {
		switch {
		case len(inc.Path) != 0:
			matches, err := rc.FS().Glob(inc.Path)
			if err != nil {
				_, err2 := rc.FS().Stat(inc.Path)
				if err2 != nil {
					return nil, err
				}

				includeFiles = append(includeFiles, inc.Path)
			} else {
				includeFiles = append(includeFiles, matches...)
			}
		case len(inc.Text) != 0:
			includeText = append(includeText, inc.Text)
		}
	}

	tpl := templateutils.CreateTemplate(rc)

	// arrange included template files in order
	// so we can include templates without {{ define "name" }} block
	// by filename and index in include

	definedTemplates := make(map[string]struct{})
	var tplList []*template.Template

	for _, inc := range includeFiles {
		// TODO: cache template files in memory
		// 	     maybe also parsed templates if we are sure rendering context
		// 	     is handled correctly

		tplBytes, err := rc.FS().ReadFile(inc)
		if err != nil {
			return nil, fmt.Errorf("failed to load template file: %q", err)
		}

		name := filepath.Base(inc)
		incTpl, err := tpl.New(name).Parse(string(tplBytes))
		if err != nil {
			return nil, fmt.Errorf("invalid template file %q: %w", inc, err)
		}

		tplList = append(tplList, incTpl)

		definedTemplates[name] = struct{}{}
	}

	tplListSize := int64(len(tplList))

	for i, inc := range includeText {
		name := "#" + strconv.FormatInt(int64(i), 10)
		_, err := tpl.New(name).Parse(inc)
		if err != nil {
			return nil, fmt.Errorf("invalid template text %s: %w", inc, err)
		}
	}

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
