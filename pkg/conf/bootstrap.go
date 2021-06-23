package conf

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"arhat.dev/pkg/envhelper"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/tools"
)

type bootstrapPureConfig struct {
	// Directory to store command cache
	CacheDir string `yaml:"cacheDir"`

	Env []string `yaml:"env"`

	// The command to run scripts (files) directly
	ScriptCmd []string `yaml:"script_cmd"`
}

type BootstrapConfig struct {
	bootstrapPureConfig `yaml:"-"`
}

func (c *BootstrapConfig) Resolve() error {
	if len(c.CacheDir) == 0 {
		c.CacheDir = ".dukkha/cache"
	}

	_ = os.MkdirAll(c.CacheDir, 0750)

	if len(c.ScriptCmd) != 0 {
		return nil
	}

	// to make it consistent among all platforms, always try to find `sh` first
	switch runtime.GOOS {
	case "windows":
		_, err := exec.LookPath("sh")
		if err == nil {
			c.ScriptCmd = []string{"sh"}
			break
		}

		_, err = exec.LookPath("pwsh")
		if err == nil {
			c.ScriptCmd = []string{"pwsh"}
			break
		}

		_, err = exec.LookPath("powershell")
		if err == nil {
			c.ScriptCmd = []string{"powershell"}
			break
		}

		c.ScriptCmd = []string{"cmd"}
	default:
		c.ScriptCmd = []string{"sh"}
	}

	return nil
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

		dec := yaml.NewDecoder(strings.NewReader(preparedDataStr))
		dec.KnownFields(true)

		err = dec.Decode(&c.bootstrapPureConfig)
		if err != nil {
			return fmt.Errorf("conf.bootstrap: failed to unmarshal config: %w", err)
		}
	}

	return nil
}

func (c *BootstrapConfig) GetExecSpec(script string, isFilePath bool) (env, cmd []string, err error) {
	scriptPath := script
	if !isFilePath {
		scriptPath, err = tools.GetScriptCache(c.CacheDir, script)
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: failed to ensure script cache: %w", err)
		}
	}

	cmd = append(cmd, c.ScriptCmd...)
	return c.Env, append(cmd, scriptPath), nil
}
