package conf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"arhat.dev/pkg/envhelper"
	"arhat.dev/pkg/exechelper"
	"gopkg.in/yaml.v3"
)

type BootstrapConfig struct {
	Env       []string `yaml:"env"`
	Shell     string   `yaml:"shell"`
	ShellArgs []string `yaml:"shell_args"`
}

func (c *BootstrapConfig) Exec(script string, spec *exechelper.Spec) (int, error) {
	spec.Command = append(append([]string{c.Shell}, c.ShellArgs...), script)

	cmd, err := exechelper.Do(*spec)
	if err != nil {
		return 128, fmt.Errorf("conf.Bootstrap.exec: %w", err)
	}

	exitCode, err := cmd.Wait()
	if err != nil {
		return exitCode, fmt.Errorf("conf.Bootstrap.exec: %w", err)
	}

	return 0, nil
}

func (c *BootstrapConfig) Resolve(ctx context.Context, data interface{}) error {
	var (
		configBytes []byte
		err         error
	)

	switch t := data.(type) {
	case nil: // use default config
	case string:
		configBytes = []byte(t)
	case []byte:
		configBytes = t
	default:
		configBytes, err = yaml.Marshal(data)
		if err != nil {
			return fmt.Errorf("conf.bootstrap: unable to marshal: %w", err)
		}
	}

	if len(configBytes) != 0 {
		preparedDataStr := envhelper.Expand(string(configBytes), func(varName, origin string) string {
			val, ok := os.LookupEnv(varName)
			if !ok {
				err = fmt.Errorf("conf.bootstrap: environment variable %q not found", val)
				return ""
			}

			return val
		})
		if err != nil {
			return fmt.Errorf("conf.bootstrap: expansion failed: %w", err)
		}

		err = Unmarshal(strings.NewReader(preparedDataStr), c)
		if err != nil {
			return fmt.Errorf("conf.bootstrap: failed to unmarshal config: %w", err)
		}
	}

	if len(c.Shell) == 0 {
		c.Shell = os.Getenv("SHELL")
		c.ShellArgs = getDefaultShellArgs(c.Shell)
	}

	if len(c.Shell) == 0 {
		switch runtime.GOOS {
		case "windows":
			_, err := exec.LookPath("sh")
			if err == nil {
				c.Shell = "sh"
				break
			}

			_, err = exec.LookPath("pwsh")
			if err == nil {
				c.Shell = "pwsh"
				break
			}

			_, err = exec.LookPath("powershell")
			if err == nil {
				c.Shell = "powershell"
				break
			}

			c.Shell = "cmd"
		default:
			c.Shell = "sh"
		}

		c.ShellArgs = getDefaultShellArgs(c.Shell)
	}

	return nil
}

func getDefaultShellArgs(shell string) []string {
	switch {
	case strings.HasSuffix(shell, "sh"),
		shell == "powershell":
		// sh, bash, zsh, pwsh
		return []string{"-c"}
	case shell == "cmd":
		return []string{"/c"}
	}

	return nil
}
