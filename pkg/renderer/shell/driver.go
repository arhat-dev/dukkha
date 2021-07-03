package shell

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "shell"

var _ renderer.Interface = (*Driver)(nil)

type Driver struct {
	getExecSpec field.ExecSpecGetFunc
}

func (d *Driver) Name() string {
	return DefaultName
}

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
		env, cmd, err := d.getExecSpec([]string{script}, false)
		if err != nil {
			return "", fmt.Errorf(
				"renderer.%s: failed to get exec spec: %w",
				DefaultName, err,
			)
		}

		execCtx := ctx.Clone()
		execCtx.AddEnv(env...)

		p, err := exechelper.Do(exechelper.Spec{
			Context: execCtx.Context(),
			Command: cmd,
			Env:     execCtx.Values().Env,

			Stdout: buf,
			Stderr: os.Stderr,
		})

		if err != nil {
			return "", fmt.Errorf(
				"renderer.%s: failed to run script [%s]: %w",
				DefaultName, strings.Join(cmd, " "), err,
			)
		}

		_, err = p.Wait()
		if err != nil {
			return "", fmt.Errorf(
				"renderer.%s: cmd failed: %w",
				DefaultName, err,
			)
		}
	}

	return buf.String(), nil
}
