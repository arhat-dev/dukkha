package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/muesli/termenv"

	"arhat.dev/dukkha/pkg/dukkha"
)

func WriteExecStart(
	stdout io.Writer,
	prefixColor termenv.Color,
	k dukkha.ToolKey,
	cmd []string,
	scriptName string,
) {
	output := []string{
		">>>",
		// task name
		string(k.Name),
		// commands
		"[", strings.Join(cmd, " "), "]",
	}

	if len(scriptName) != 0 {
		output = append(output, "@", scriptName)
	}

	if prefixColor != nil {
		printlnWithColor(stdout, output, prefixColor)
	} else {
		_, _ = fmt.Fprintln(stdout, strings.Join(output, " "))
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
	resultKind := "DONE"
	if err != nil {
		resultKind = "ERROR"
	}

	output := []string{
		resultKind,
		AssembleTaskKindID(k, tk.Kind),
		"[", string(tk.Name), "]",
		"{", matrixSpec,
	}

	if err != nil {
		output = append(output, "}:", err.Error())
	} else {
		output = append(output, "}")
	}

	if prefixColor != nil {
		printlnWithColor(stderr, output, prefixColor)
	} else {
		_, _ = fmt.Fprintln(stderr, strings.Join(output, " "))
	}
}
