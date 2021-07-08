package shell

import (
	"bytes"
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/types"
)

const DefaultName = "shell"

func New(getExecSpec dukkha.ExecSpecGetFunc) dukkha.Renderer {
	return &driver{getExecSpec: getExecSpec}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	getExecSpec dukkha.ExecSpecGetFunc
}

func (d *driver) Name() string { return DefaultName }

func (d *driver) RenderYaml(rc types.RenderingContext, rawData interface{}) ([]byte, error) {
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
