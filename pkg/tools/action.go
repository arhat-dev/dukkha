package tools

import (
	"fmt"
	"strconv"
	"sync"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
)

type Action struct {
	rs.BaseField

	// Name of this action, optional
	Name string `yaml:"name"`

	// Env specific to this action
	Env dukkha.Env `yaml:"env"`

	// Task reference of this action
	//
	// Task, Cmd, EmbeddedShell, ExternalShell are mutually exclusive
	Task string `yaml:"task"`

	// EmbeddedShell using embedded shell
	//
	// Task, Cmd, EmbeddedShell, ExternalShell are mutually exclusive
	EmbeddedShell string `yaml:"shell"`

	// EmbeddedShell script for this action
	//
	// Task, Cmd, EmbeddedShell, ExternalShell are mutually exclusive
	ExternalShell map[string]string `rs:"other"`

	// Cmd execution, not in any shell
	//
	// Task, Cmd, EmbeddedShell, ExternalShell are mutually exclusive
	Cmd []string `yaml:"cmd"`

	// Chdir change working directory before executing command
	// this option only applies to Cmd, EmbeddedShell, ExternalShell action
	Chdir string `yaml:"chdir"`

	// ContuineOnError ignores error occurred in this action and continue
	// following actions in list (if any)
	ContinueOnError bool `yaml:"continue_on_error"`

	mu sync.Mutex
}

func (act *Action) DoAfterFieldResolved(mCtx dukkha.TaskExecContext, do func(h *Action) error) error {
	act.mu.Lock()
	defer act.mu.Unlock()

	err := dukkha.ResolveEnv(act, mCtx, "Env")
	if err != nil {
		return fmt.Errorf("failed to resolve action specific env: %w", err)
	}

	err = act.ResolveFields(mCtx, -1)
	if err != nil {
		return fmt.Errorf("failed to resolve fields: %w", err)
	}

	return do(act)
}

func (act *Action) GenSpecs(
	ctx dukkha.TaskExecContext, index int,
) (dukkha.RunTaskOrRunCmd, error) {
	hookID := "#" + strconv.FormatInt(int64(index), 10)
	if len(act.Name) != 0 {
		hookID = fmt.Sprintf("%s (%s)", act.Name, hookID)
	}

	switch {
	case len(act.Task) != 0:
		return act.genTaskActionSpecs(ctx, hookID)
	case len(act.Cmd) != 0:
		return act.genCmdActionSpecs(ctx, hookID)
	case len(act.EmbeddedShell) != 0:
		return act.genEmbeddedShellActionSpecs(ctx, hookID)
	default:
		return act.genExternalShellActionSpecs(ctx, hookID)
	}
}
