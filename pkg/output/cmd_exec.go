package output

import (
	"fmt"
	"os"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"github.com/fatih/color"
)

func WriteExecStart(
	prefixColor *color.Color,
	k dukkha.ToolKey,
	cmd []string,
	scriptName string,
) {
	output := []interface{}{
		">>>",
		// task name
		k.Name,
		// commands
		"[", strings.Join(cmd, " "), "]",
	}

	if len(scriptName) != 0 {
		output = append(output, "@", scriptName)
	}

	if prefixColor != nil {
		_, _ = prefixColor.Println(output...)
	} else {
		_, _ = fmt.Println(output...)
	}
}

func WriteExecResult(
	prefixColor *color.Color,
	k dukkha.ToolKey,
	tk dukkha.TaskKey,
	matrixSpec string,
	err error,
) {
	resultKind := "DONE"
	if err != nil {
		resultKind = "ERROR"
	}

	output := []interface{}{
		resultKind,
		AssembleTaskKindID(k, tk.Kind),
		"[", tk.Name, "]",
		"{", matrixSpec,
	}

	if err != nil {
		output = append(output, "}:", err.Error())
	} else {
		output = append(output, "}")
	}

	if prefixColor != nil {
		_, _ = prefixColor.Fprintln(os.Stderr, output...)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, output...)
	}
}
