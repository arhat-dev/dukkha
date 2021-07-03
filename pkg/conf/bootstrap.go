package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"arhat.dev/pkg/envhelper"
	"arhat.dev/pkg/log"
	"go.uber.org/multierr"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/tools"
)

type BootstrapConfig struct {
	// CacheDir to store script file and temporary task execution data
	// it is used when
	// 		- resolving fields using shell renderer
	// 		- executing tasks need to run commands
	CacheDir string `yaml:"cache_dir"`

	// Env
	Env []string `yaml:"env"`

	// The command to run scripts (files) directly
	ScriptCmd []string `yaml:"script_cmd"`
}

// Resolve bootstrap config
//
// 1. resolve env and set these environment variables as global env
// 2. resolve cache_dir, set global env DUKKHA_CACHE_DIR to its absolute path
// 3. resolve (expand) script_cmd with global env
func (c *BootstrapConfig) Resolve() error {
	logger := log.Log.WithName("bootstrap.resovle")

	var err error
	expandEnvFunc := func(varName, origin string) string {
		logger.V("expanding env", log.String("origin", origin))

		if strings.HasPrefix(origin, "$(") {
			err = multierr.Append(
				err,
				fmt.Errorf("shell evaluation %q is not allowed", origin),
			)
			return ""
		}

		val, ok := os.LookupEnv(varName)
		if !ok {
			err = multierr.Append(
				err,
				fmt.Errorf("environment variable %q not found", val),
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

		logger.V("setting global env",
			log.String("name", key),
			log.String("vale", value),
		)
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
		c.CacheDir = constant.DefaultCacheDir
	}

	c.CacheDir, err = filepath.Abs(c.CacheDir)
	if err != nil {
		return fmt.Errorf("bootstrap: failed to get absolute path of cache dir: %w", err)
	}

	logger.V("resolved dukkha cache dir", log.String("path", c.CacheDir))
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
				"bootstrap: unable to expand part %q in script_cmd: %w",
				cmdPart, err,
			)
		}
	}

	if len(c.ScriptCmd) == 0 {
		// to make it consistent among all platforms, always defaults to `sh`
		c.ScriptCmd = []string{"sh"}
	}

	logger.V("resolved script cmd", log.Strings("cmd", c.ScriptCmd))

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
