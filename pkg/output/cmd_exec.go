package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/muesli/termenv"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
)

func WriteExecStart(
	stdout io.Writer,
	prefixColor termenv.Color,
	k dukkha.ToolKey,
	cmds []string,
	scriptName string,
) {
	var sb strings.Builder
	sb.WriteString(">>> ")
	sb.WriteString(string(k.Name))
	sb.WriteString(" [")
	for _, c := range cmds {
		sb.WriteByte(' ')

		s, err := syntax.Quote(c, syntax.LangBash)
		if err != nil {
			sb.WriteString(c)
		} else {
			sb.WriteString(s)
		}

		sb.WriteByte(' ')
	}
	sb.WriteString("]")

	if len(scriptName) != 0 {
		sb.WriteString(" @ ")
		sb.WriteString(scriptName)
	}

	if prefixColor != nil {
		printlnWithColor(stdout, sb.String(), prefixColor)
	} else {
		_, _ = fmt.Fprintln(stdout, sb.String())
	}
}

func WriteExecResult(
	stderr io.Writer,
	prefixColor termenv.Color,
	k dukkha.ToolKey,
	tk dukkha.TaskKey,
	matrixSpec string,
	err error,
) {
	var sb strings.Builder
	if err != nil {
		sb.WriteString("ERROR ")
	} else {
		sb.WriteString("DONE ")
	}

	sb.WriteString(AssembleTaskKindID(k, tk.Kind))
	sb.WriteString(" [ ")
	sb.WriteString(string(tk.Name))
	sb.WriteString(" ] { ")
	sb.WriteString(matrixSpec)

	if err != nil {
		sb.WriteString(" }: ")
		sb.WriteString(err.Error())
	} else {
		sb.WriteString(" }")
	}

	if prefixColor != nil {
		printlnWithColor(stderr, sb.String(), prefixColor)
	} else {
		_, _ = fmt.Fprintln(stderr, sb.String())
	}
}
