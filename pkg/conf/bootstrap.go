package conf

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"arhat.dev/pkg/envhelper"
	"arhat.dev/pkg/exechelper"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"
)

type bootstrapPureConfig struct {
	Env       []string `yaml:"env"`
	Shell     string   `yaml:"shell"`
	ShellArgs []string `yaml:"shell_args"`
}

type BootstrapConfig struct {
	bootstrapPureConfig `yaml:"-"`
}

func (c *BootstrapConfig) UnmarshalYAML(n *yaml.Node) error {
	configBytes, err := yaml.Marshal(n)
	if err != nil {
		return fmt.Errorf("bootstrap: marshal back failed: %w", err)
	}

	if len(configBytes) != 0 {
		preparedDataStr := envhelper.Expand(string(configBytes), func(varName, origin string) string {
			if strings.HasPrefix(origin, "$(") {
				err = multierr.Append(
					err,
					fmt.Errorf("conf.bootstrap: shell evaluation %q is not allowed", origin),
				)
				return ""
			}

			val, ok := os.LookupEnv(varName)
			if !ok {
				err = multierr.Append(
					err,
					fmt.Errorf("conf.bootstrap: environment variable %q not found", val),
				)
				return ""
			}

			return val
		})
		if err != nil {
			return fmt.Errorf("conf.bootstrap: config expansion failed: %w", err)
		}

		err = yaml.Unmarshal([]byte(preparedDataStr), &c.bootstrapPureConfig)
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
