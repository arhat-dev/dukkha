package shell

import (
	"bytes"
	"fmt"
	"os"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/templateutils"
)

const DefaultName = "shell"

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

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

	var scripts []string
	switch t := rawData.(type) {
	case string:
		scripts = append(scripts, t)
	case []byte:
		scripts = append(scripts, string(t))
	case []interface{}:
		for _, v := range t {
			var scriptBytes []byte
			scriptBytes, err = yamlhelper.ToYamlBytes(v)
			if err != nil {
				return nil, fmt.Errorf(
					"renderer.%s: unexpected list item type %T: %w",
					d.name, v, err,
				)
			}

			scripts = append(scripts, string(scriptBytes))
		}
	default:
		return nil, fmt.Errorf(
			"renderer.%s: unsupported input type %T", d.name, rawData,
		)
	}

	buf := &bytes.Buffer{}
	runner, err := templateutils.CreateEmbeddedShellRunner(
		rc.WorkDir(), rc, nil, buf, os.Stderr,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: creating embedded shell: %w", d.name, err,
		)
	}

	parser := syntax.NewParser(
		syntax.Variant(syntax.LangBash),
	)

	for _, script := range scripts {
		err = templateutils.RunScriptInEmbeddedShell(rc, runner, parser, script)

		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: %w",
				d.name, err,
			)
		}
	}

	return buf.Next(buf.Len()), nil
}
