package output

import (
	"context"
	"fmt"
	"os"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
)

func WriteExecStart(name string, cmd []string, scriptName string) {
	fmt.Println(
		">>>",
		// task name
		name,
		// commands
		"[", strings.Join(cmd, " "), "]",
		// script filename prefix
		"@", scriptName[:7],
	)
}

func WriteExecResult(
	ctx context.Context,
	toolKind dukkha.ToolKind,
	toolName dukkha.ToolName,
	taskKind dukkha.TaskKind,
	taskName dukkha.TaskName,
	matrixSpec string,
	err error,
) {
	resultKind := "DONE"
	if err != nil {
		resultKind = "ERROR"
	}

	output := []interface{}{
		resultKind,
		AssembleTaskKindID(toolKind, toolName, taskKind),
		"[", taskName, "]",
		"{", matrixSpec,
	}

	if err != nil {
		output = append(output, "}:", err.Error())
	} else {
		output = append(output, "}")
	}

	_, _ = fmt.Fprintln(os.Stderr, output...)
}
