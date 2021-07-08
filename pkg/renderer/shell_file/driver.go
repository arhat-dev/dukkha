package shell_file

import (
	"bytes"
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/types"
)

const DefaultName = "shell_file"

func New(getExecSpec dukkha.ExecSpecGetFunc) dukkha.Renderer {
	return &driver{getExecSpec: getExecSpec}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	getExecSpec dukkha.ExecSpecGetFunc
}

func (d *driver) Name() string { return DefaultName }

func (d *driver) RenderYaml(rc types.RenderingContext, rawData interface{}) ([]byte, error) {
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
