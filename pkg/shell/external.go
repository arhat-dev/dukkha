package shell

import (
	"fmt"
	"io"
	"os"
	"strings"

	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/dukkha"
)

func RunScript(
	rc dukkha.RenderingContext,
	script string,
	isFilePath bool,
	stdout io.Writer,
	getExecSpec dukkha.ExecSpecGetFunc,
) error {
	env, cmd, err := getExecSpec([]string{script}, false)
	if err != nil {
		return fmt.Errorf("failed to get exec spec: %w", err)
	}

	rc.AddEnv(env...)

	p, err := exechelper.Do(exechelper.Spec{
		Context: rc,
		Command: cmd,
		Env:     rc.Env(),

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
