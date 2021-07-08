package shell

import (
	"bytes"
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "shell"

func init() {
	dukkha.RegisterRenderer(
		DefaultName,
		func() dukkha.Renderer {
			return NewDefault(nil)
		},
	)
}

func NewDefault(getExecSpec dukkha.ExecSpecGetFunc) dukkha.Renderer {
	return &driver{getExecSpec: getExecSpec}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	field.BaseField

	getExecSpec dukkha.ExecSpecGetFunc
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	allShells := ctx.AllShells()
	for shellName := range allShells {
		name := DefaultName
		if len(shellName) == 0 {
			name += ":" + shellName
		}

		ctx.AddRenderer(
			name, &driver{
				getExecSpec: allShells[shellName].GetExecSpec,
			},
		)
	}

	return nil
}

func (d *driver) RenderYaml(rc dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	var scripts []string
	switch t := rawData.(type) {
	case string:
		scripts = append(scripts, t)
	case []byte:
		scripts = append(scripts, string(t))
	case []interface{}:
		for _, v := range t {
			scriptBytes, err := renderer.ToYamlBytes(v)
			if err != nil {
				return nil, fmt.Errorf("renderer.%s: unexpected list item type %T: %w", DefaultName, v, err)
			}

			scripts = append(scripts, string(scriptBytes))
		}
	default:
		return nil, fmt.Errorf("renderer.%s: unsupported input type %T", DefaultName, rawData)
	}

	buf := &bytes.Buffer{}
	for _, script := range scripts {
		err := renderer.RunShellScript(rc, script, false, buf, d.getExecSpec)
		if err != nil {
			return nil, fmt.Errorf("renderer.%s: %w", DefaultName, err)
		}
	}

	return buf.Bytes(), nil
}
