package shell_file

import (
	"bytes"
	"fmt"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "shell_file"

func init() {
	renderer.Register(&Config{}, NewDriver)
}

func NewDriver(config interface{}) (renderer.Interface, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("unexpected non %s renderer config: %T", DefaultName, config)
	}

	if cfg.GetExecSpec == nil {
		return nil, fmt.Errorf("required exec func not set")
	}

	return &Driver{getExecSpec: cfg.GetExecSpec}, nil
}

var _ renderer.Config = (*Config)(nil)

type Config struct {
	GetExecSpec field.ExecSpecGetFunc
}

var _ renderer.Interface = (*Driver)(nil)

type Driver struct {
	getExecSpec field.ExecSpecGetFunc
}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(ctx *field.RenderingContext, rawData interface{}) (string, error) {
	scriptPath, ok := rawData.(string)
	if !ok {
		return "", fmt.Errorf("renderer.%s: unexpected non-string input %T", DefaultName, rawData)
	}

	buf := &bytes.Buffer{}
	err := renderer.RunShellScript(ctx, scriptPath, true, buf, d.getExecSpec)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: %w", DefaultName, err)
	}

	return buf.String(), nil
}
