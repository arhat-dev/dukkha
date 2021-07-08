package shell_file

import (
	"bytes"
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "shell_file"

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
	scriptPath, ok := rawData.(string)
	if !ok {
		return nil, fmt.Errorf("renderer.%s: unexpected non-string input %T", DefaultName, rawData)
	}

	buf := &bytes.Buffer{}
	err := renderer.RunShellScript(rc, scriptPath, true, buf, d.getExecSpec)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", DefaultName, err)
	}

	return buf.Bytes(), nil
}
