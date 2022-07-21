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
	var sb strings.Builder
	sb.WriteString("--- ")
	sb.WriteString(AssembleTaskKindID(k, tk.Kind))
	sb.WriteString(" [ ")
	sb.WriteString(string(tk.Name))
	sb.WriteString(" ] { ")
	sb.WriteString(matrixSpec.String())
	sb.WriteString(" }")

	if prefixColor != nil {
		printlnWithColor(stdout, sb.String(), prefixColor)
	} else {
		_, _ = fmt.Fprintln(stdout, sb.String())
	}
}

func printlnWithColor(stdout io.Writer, str string, color termenv.Color) {
	style := termenv.String(str).Foreground(color)
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
