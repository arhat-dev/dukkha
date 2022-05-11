package templateutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/pkg/stringhelper"
)

func createEvalNS(rc dukkha.RenderingContext) evalNS { return evalNS{rc: rc} }

type evalNS struct{ rc dukkha.RenderingContext }

func (ns evalNS) Template(tplData interface{}) (string, error) {
	var tplStr string
	switch tt := tplData.(type) {
	case string:
		tplStr = tt
	case []byte:
		tplStr = string(tt)
	default:
		return "", fmt.Errorf(
			"invalid template data, want string or bytes, got %T",
			tt,
		)
	}

	tpl, err := CreateTemplate(ns.rc).Parse(tplStr)
	if err != nil {
		return "", fmt.Errorf(
			"parsing template\n\n%s\n\n %w",
			tplStr, err,
		)
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, ns.rc)
	if err != nil {
		return "", fmt.Errorf(
			"evaluating template\n\n%s\n\n %w",
			tplStr, err,
		)
	}

	return stringhelper.Convert[string, byte](buf.Next(buf.Len())), nil
}

// Env expands environment variables in last argument, which can be string or bytes
//
// valid options before last argument are
// `disable_exec` / `enable_exec` to deny (default behvior) / allow shell script evaluation during expansion
func (ns evalNS) Env(args ...String) (string, error) {
	var textData String
	switch len(args) {
	case 0:
		return "", nil
	case 1:
		textData = args[0]
	default:
		textData = args[len(args)-1]
	}

	enableExec := false
	for _, opt := range args[:len(args)-1] {
		switch toString(opt) {
		case "disable_exec":
			enableExec = false
		case "enable_exec":
			enableExec = true
		}
	}

	return ExpandEnv(ns.rc, toString(textData), enableExec)
}

// Shell evaluates scriptData as bash script, inputs as stdin streams (if any)
//
// inputs can be string, bytes, or io.Reader
func (ns evalNS) Shell(script String, inputs ...interface{}) (string, error) {
	var stdin io.Reader
	if len(inputs) != 0 {
		var readers []io.Reader
		for _, in := range inputs {
			switch it := in.(type) {
			case io.Reader:
				readers = append(readers, it)
			case string:
				readers = append(readers, strings.NewReader(it))
			case []byte:
				readers = append(readers, bytes.NewReader(it))
			default:
				return "", fmt.Errorf(
					"invalid input type want reader, string or bytes, got %T",
					it,
				)
			}

		}

		stdin = io.MultiReader(readers...)
	}

	var stdout bytes.Buffer
	runner, err := CreateEmbeddedShellRunner(
		ns.rc.WorkDir(),
		ns.rc,
		stdin,
		&stdout,
		os.Stderr,
	)
	if err != nil {
		return "", err
	}

	err = RunScriptInEmbeddedShell(
		ns.rc,
		runner,
		syntax.NewParser(),
		toString(script),
	)
	return stringhelper.Convert[string, byte](stdout.Next(stdout.Len())), err
}
