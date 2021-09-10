//go:build !real
// +build !real

package tools

import (
	"arhat.dev/dukkha/pkg/dukkha"
)

func (act *Action) genTaskActionSpecs(
	ctx dukkha.TaskExecContext, hookID string,
) (dukkha.RunTaskOrRunCmd, error) {
	return nil, nil
}

func (act *Action) genCmdActionSpecs(
	ctx dukkha.TaskExecContext, hookID string,
) (dukkha.RunTaskOrRunCmd, error) {
	return nil, nil
}

func (act *Action) genEmbeddedShellActionSpecs(
	ctx dukkha.TaskExecContext, hookID string,
) (dukkha.RunTaskOrRunCmd, error) {
	return nil, nil
}

func (act *Action) genExternalShellActionSpecs(
	ctx dukkha.TaskExecContext, hookID string,
) (dukkha.RunTaskOrRunCmd, error) {
	return nil, nil
}
