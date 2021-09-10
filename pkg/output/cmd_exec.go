package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/muesli/termenv"

	"arhat.dev/dukkha/pkg/dukkha"
)

func WriteExecStart(
	prefixColor dukkha.TermColor,
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

	if prefixColor != 0 {
		printlnWithColor(output, termenv.ANSIColor(prefixColor))
	} else {
		_, _ = fmt.Println(strings.Join(output, " "))
	}
}

func WriteExecResult(
	prefixColor dukkha.TermColor,
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

	if prefixColor != 0 {
		printlnWithColor(output, termenv.ANSIColor(prefixColor))
	} else {
		_, _ = fmt.Fprintln(os.Stderr, strings.Join(output, " "))
	}
}
