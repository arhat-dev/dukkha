package output

import (
	"fmt"
	"strings"

	"github.com/muesli/termenv"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
)

func WriteTaskStart(
	prefixColor dukkha.TermColor,
	k dukkha.ToolKey,
	tk dukkha.TaskKey,
	matrixSpec matrix.Entry,
) {
	output := []string{
		"---",
		AssembleTaskKindID(k, tk.Kind),
		"[", string(tk.Name), "]",
		"{", matrixSpec.String(), "}",
	}

	if prefixColor != 0 {
		printlnWithColor(output, termenv.ANSIColor(prefixColor))
	} else {
		_, _ = fmt.Println(strings.Join(output, " "))
	}
}

func printlnWithColor(parts []string, color termenv.Color) {
	style := termenv.String(parts...).Foreground(color)
	_, _ = fmt.Println(style.String())
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
