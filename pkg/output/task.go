package output

import (
	"context"
	"fmt"
	"strings"
)

func WriteTaskStart(
	ctx context.Context,
	toolKind, toolName, taskKind, taskName string,
	matrixSpec string,
) {
	_, _ = fmt.Println(
		"---",
		AssembleTaskKindID(toolKind, toolName, taskKind),
		"[", taskName, "]",
		"{", matrixSpec, "}",
	)
}

func AssembleTaskKindID(toolKind, toolName, taskKind string) string {
	kindParts := []string{toolKind}
	if len(toolName) != 0 {
		kindParts = append(kindParts, toolName)
	}

	return strings.Join(append(kindParts, taskKind), ":")
}
