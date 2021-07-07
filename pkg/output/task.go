package output

import (
	"fmt"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/types"
)

func WriteTaskStart(
	toolKind dukkha.ToolKind,
	toolName dukkha.ToolName,
	taskKind dukkha.TaskKind,
	taskName dukkha.TaskName,
	matrixSpec types.MatrixSpec,
) {
	_, _ = fmt.Println(
		"---",
		AssembleTaskKindID(toolKind, toolName, taskKind),
		"[", taskName, "]",
		"{", matrixSpec.String(), "}",
	)
}

func AssembleTaskKindID(
	toolKind dukkha.ToolKind,
	toolName dukkha.ToolName,
	taskKind dukkha.TaskKind,
) string {
	kindParts := []string{string(toolKind)}
	if len(toolName) != 0 {
		kindParts = append(kindParts, string(toolName))
	}

	return strings.Join(append(kindParts, string(taskKind)), ":")
}
