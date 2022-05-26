package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/muesli/termenv"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
)

func WriteTaskStart(
	stdout io.Writer,
	prefixColor termenv.Color,
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

	if prefixColor != nil {
		printlnWithColor(stdout, output, prefixColor)
	} else {
		_, _ = fmt.Fprintln(stdout, strings.Join(output, " "))
	}
}

func printlnWithColor(stdout io.Writer, parts []string, color termenv.Color) {
	style := termenv.String(parts...).Foreground(color)
	_, _ = fmt.Fprintln(stdout, style.String())
}

func AssembleTaskKindID(
	k dukkha.ToolKey,
	taskKind dukkha.TaskKind,
) string {
	var sb strings.Builder
	sb.WriteString(string(k.Kind))

	if len(k.Name) != 0 {
		sb.WriteString(":")
		sb.WriteString(string(k.Name))
	}

	sb.WriteString(":")
	sb.WriteString(string(taskKind))

	return sb.String()
}
