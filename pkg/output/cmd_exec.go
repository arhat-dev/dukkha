package output

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func WriteExecStart(ctx context.Context, name string, cmd []string, scriptName string) {
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
	toolKind, toolName, taskKind, taskName string,
	matrixSpec string,
	err error,
) {
	resultKind := "DONE"
	if err != nil {
		resultKind = "ERROR"
	}

	output := []interface{}{
		resultKind,
		assembleTaskKindID(toolKind, toolName, taskKind),
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
