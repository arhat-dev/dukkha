package renderer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"arhat.dev/pkg/exechelper"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/field"
)

func ToYamlBytes(in interface{}) ([]byte, error) {
	switch t := in.(type) {
	case string:
		return []byte(t), nil
	case []byte:
		return t, nil
	default:
	}

	ret, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func RunShellScript(
	rc *field.RenderingContext,
	script string,
	isFilePath bool,
	stdout io.Writer,
	getExecSpec field.ExecSpecGetFunc,
) error {
	env, cmd, err := getExecSpec([]string{script}, false)
	if err != nil {
		return fmt.Errorf("failed to get exec spec: %w", err)
	}

	execCtx := rc.Clone()
	execCtx.AddEnv(env...)

	p, err := exechelper.Do(exechelper.Spec{
		Context: execCtx.Context(),
		Command: cmd,
		Env:     execCtx.Values().Env,

		Stdout: stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		return fmt.Errorf("failed to run script [%s]: %w",
			strings.Join(cmd, " "), err,
		)
	}

	_, err = p.Wait()
	if err != nil {
		return fmt.Errorf("cmd failed: %w", err)
	}

	return nil
}
