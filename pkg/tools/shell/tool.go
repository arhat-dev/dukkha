package tools_shell

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/sha256helper"
)

var _ dukkha.Tool = (*Tool)(nil)

type Shell string

func (s *Shell) DefaultExecutable() string { return *(*string)(s) }
func (s *Shell) Kind() dukkha.ToolKind     { return "shell" }

type Tool struct {
	tools.BaseTool[Shell, *Shell]
}

func (t *Tool) InitWithName(name string, cacheFS *fshelper.OSFS) error {
	t.Impl = Shell(name)
	return t.BaseTool.Init(cacheFS)
}

// GetExecSpec is a helper func for shells
func (t *Tool) GetExecSpec(
	toExec []string, isFilePath bool,
) (env dukkha.NameValueList, cmd []string, err error) {
	if len(toExec) == 0 {
		return nil, nil, fmt.Errorf("invalid empty exec spec")
	}

	scriptPath := ""
	if !isFilePath {
		scriptPath, err = GetScriptCache(t.CacheFS, strings.Join(toExec, " "))
		if err != nil {
			return nil, nil, fmt.Errorf("unable to ensure script cache: %w", err)
		}
	} else {
		scriptPath = toExec[0]
	}

	cmd = sliceutils.NewStrings(t.Cmd)
	if len(cmd) == 0 {
		cmd = append(cmd, t.Impl.DefaultExecutable())
	}

	return t.Env, append(cmd, scriptPath), nil
}

func GetScriptCache(cacheFS *fshelper.OSFS, script string) (string, error) {
	scriptName := hex.EncodeToString(sha256helper.Sum([]byte(script)))

	_, err := cacheFS.Stat(scriptName)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("check existence of script cache: %w", err)
		}

		err = cacheFS.WriteFile(scriptName, []byte(script), 0600)
		if err != nil {
			return "", fmt.Errorf("writing script cache: %w", err)
		}
	}

	return cacheFS.Abs(scriptName)
}
