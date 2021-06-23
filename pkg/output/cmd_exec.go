package output

import (
	"context"
	"fmt"
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

func WriteExecFailure() {
	// TODO
}

func WriteExecSuccess() {
	// TODO
}
