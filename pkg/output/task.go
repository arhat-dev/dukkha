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
	kindParts := []string{toolKind}
	if len(toolName) != 0 {
		kindParts = append(kindParts, toolName)
	}

	kindParts = append(kindParts, taskKind)

	fmt.Println(
		"---",
		strings.Join(kindParts, ":"),
		"[", taskName, "]",
		"{", matrixSpec, "}",
	)
}
