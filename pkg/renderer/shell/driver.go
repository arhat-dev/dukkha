package shell

import (
	"bytes"
	"fmt"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/shell_file"
)

const DefaultName = "shell"

func init() {
	renderer.Register(&Config{}, NewDriver)
}

func NewDriver(config interface{}) (renderer.Interface, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("unexpected non %s renderer config: %T", DefaultName, config)
	}

	if cfg.GetExecSpec == nil {
		return nil, fmt.Errorf("required GetExecSpec func not set")
	}

	return &Driver{getExecSpec: cfg.GetExecSpec}, nil
}

var _ renderer.Config = (*Config)(nil)

type Config shell_file.Config

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
