package output

import (
	"fmt"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
	"github.com/fatih/color"
)

func WriteTaskStart(
	prefixColor *color.Color,
	k dukkha.ToolKey,
	tk dukkha.TaskKey,
	matrixSpec matrix.Entry,
) {
	output := []interface{}{
		"---",
		AssembleTaskKindID(k, tk.Kind),
		"[", tk.Name, "]",
		"{", matrixSpec.String(), "}",
	}

	if prefixColor != nil {
		_, _ = prefixColor.Println(output...)
	} else {
		_, _ = fmt.Println(output...)
	}
}

func AssembleTaskKindID(
	k dukkha.ToolKey,
	taskKind dukkha.TaskKind,
) string {
	kindParts := []string{string(k.Kind)}
	if len(k.Name) != 0 {
		kindParts = append(kindParts, string(k.Name))
	}

	return strings.Join(append(kindParts, string(taskKind)), ":")
}
