package shell

import (
	"bytes"
	"fmt"

	"arhat.dev/pkg/stringhelper"
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
) (_ []byte, err error) {
	rawData, err = rs.NormalizeRawData(rawData)
	if err != nil {
		return
	}

	var (
		bufScripts [1]string
		scripts    []string
	)
	switch t := rawData.(type) {
	case string:
		bufScripts[0] = t
		scripts = bufScripts[:]
	case []byte:
		bufScripts[0] = stringhelper.Convert[string, byte](t)
		scripts = bufScripts[:]
	case []interface{}:
		for _, v := range t {
			var scriptBytes []byte
			scriptBytes, err = yamlhelper.ToYamlBytes(v)
			if err != nil {
				return nil, fmt.Errorf(
					"renderer.%s: unexpected list item type %T: %w", d.name, v, err,
				)
			}

			scripts = append(scripts, stringhelper.Convert[string, byte](scriptBytes))
		}
	default:
		return nil, fmt.Errorf("renderer.%s: unsupported input type %T", d.name, rawData)
	}

	var buf bytes.Buffer
	runner, err := templateutils.CreateShellRunner(
		rc.WorkDir(), rc, nil, &buf, rc.Stderr(),
	)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: creating embedded shell: %w", d.name, err)
	}

	parser := syntax.NewParser(
		syntax.Variant(syntax.LangBash),
	)

	for _, script := range scripts {
		err = templateutils.RunScript(rc, runner, parser, script)

		if err != nil {
			return nil, fmt.Errorf("renderer.%s: %w", d.name, err)
		}
	}

	return buf.Next(buf.Len()), nil
}
