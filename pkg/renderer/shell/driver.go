package shell

import (
	"bytes"
	"fmt"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "shell"

var _ renderer.Interface = (*Driver)(nil)

type Driver struct {
	getExecSpec field.ExecSpecGetFunc
}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(ctx *field.RenderingContext, rawData interface{}) (string, error) {
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
				return "", fmt.Errorf("renderer.%s: unexpected list item type %T: %w", DefaultName, v, err)
			}

			scripts = append(scripts, string(scriptBytes))
		}
	default:
		return "", fmt.Errorf("renderer.%s: unsupported input type %T", DefaultName, rawData)
	}

	buf := &bytes.Buffer{}
	for _, script := range scripts {
		err := renderer.RunShellScript(ctx, script, false, buf, d.getExecSpec)
		if err != nil {
			return "", fmt.Errorf("renderer.%s: %w", DefaultName, err)
		}
	}

	return buf.String(), nil
}
