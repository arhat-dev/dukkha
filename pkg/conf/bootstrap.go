package conf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"arhat.dev/pkg/envhelper"
	"go.uber.org/multierr"

	"arhat.dev/dukkha/pkg/constant"
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
	bootstrapPureConfig `yaml:",inline"`
}

// Resolve bootstrap config
// 	- resolve env and set these environment vairables as global env
// 	- resolve cache dir, set global env DUKKHA_CACHE_DIR to its absolute path
// 	- resolve script cmd using global env
func (c *BootstrapConfig) Resolve() error {
	var err error
	expandEnvFunc := func(varName, origin string) string {
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
	}

	for i, env := range c.Env {
		realEnv := envhelper.Expand(env, expandEnvFunc)
		if err != nil {
			return fmt.Errorf("bootstrap: failed to expand env: %w", err)
		}

		c.Env[i] = realEnv

		parts := strings.SplitN(c.Env[i], "=", 2)
		key, value := parts[0], ""
		if len(parts) == 2 {
			value = parts[1]
		}

		err = os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("bootstrap: failed to set global env %q: %w", key, err)
		}
	}

	c.CacheDir = envhelper.Expand(c.CacheDir, expandEnvFunc)
	if err != nil {
		return fmt.Errorf("bootstrap: failed to resolve cache dir: %w", err)
	}

	if len(c.CacheDir) == 0 {
		c.CacheDir = ".dukkha/cache"
	}

	c.CacheDir, err = filepath.Abs(c.CacheDir)
	if err != nil {
		return fmt.Errorf("bootstrap: failed to get absolute path of cache dir: %w", err)
	}

	err = os.MkdirAll(c.CacheDir, 0750)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("bootstrap: failed to ensure cache dir: %w", err)
	}

	err = os.Setenv(constant.ENV_DUKKHA_CACHE_DIR, c.CacheDir)
	if err != nil {
		return fmt.Errorf("bootstrap: failed to set cache dir global env %q: %w",
			constant.ENV_DUKKHA_CACHE_DIR, err,
		)
	}

	for i, cmdPart := range c.ScriptCmd {
		c.ScriptCmd[i] = envhelper.Expand(cmdPart, expandEnvFunc)
		if err != nil {
			return fmt.Errorf(
				"bootstrap: unable to resolve script cmd %q: %w",
				cmdPart, err,
			)
		}
	}

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

func (c *BootstrapConfig) GetExecSpec(toExec []string, isFilePath bool) (env, cmd []string, err error) {
	if len(toExec) == 0 {
		return nil, nil, fmt.Errorf("bootstrap: invalid empty exec spec")
	}

	scriptPath := toExec[0]
	if !isFilePath {
		scriptPath, err = tools.GetScriptCache(c.CacheDir, strings.Join(toExec, " "))
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: failed to ensure script cache: %w", err)
		}
	}

	cmd = append(cmd, c.ScriptCmd...)
	return c.Env, append(cmd, scriptPath), nil
}
