package templateutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
)

func createEvalNS(rc dukkha.RenderingContext) evalNS {
	return evalNS{rc: rc}
}

type evalNS struct {
	rc dukkha.RenderingContext
}

func (ens evalNS) Template(tplData interface{}) (string, error) {
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

	tpl, err := CreateTemplate(ens.rc).Parse(tplStr)
	if err != nil {
		return "", fmt.Errorf(
			"parsing template\n\n%s\n\n %w",
			tplStr, err,
		)
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, ens.rc)
	if err != nil {
		return "", fmt.Errorf(
			"evaluating template\n\n%s\n\n %w",
			tplStr, err,
		)
	}

	return string(buf.Next(buf.Len())), nil
}

// Env expands environment variables in last argument, which can be string or bytes
// valid options before last argument are
// `disable_exec` to deny shell evaluation (default behvior)
// `enable_exec` to allow shell evaluation during expansing
func (ens evalNS) Env(args ...interface{}) (string, error) {
	var textData interface{}
	switch len(args) {
	case 0:
		return "", nil
	case 1:
		textData = args[0]
	default:
		textData = args[len(args)-1]
	}

	var toExpand string
	switch tt := textData.(type) {
	case string:
		toExpand = tt
	case []byte:
		toExpand = string(tt)
	default:
		return "", fmt.Errorf("invalid non text data for expansion, got %T", tt)
	}

	enableExec := false
	for _, opt := range args[:len(args)-1] {
		switch opt {
		case "disable_exec":
			enableExec = false
		case "enable_exec":
			enableExec = true
		}
	}

	return ExpandEnv(ens.rc, toExpand, enableExec)
}

// Shell evaluates scriptData as bash script, inputs as stdin streams (if any)
// if no input was provided, it will use os.Stdin as shell stdin
// scriptData can be string or bytes
// inputs can be string, bytes, or io.Reader
func (ens evalNS) Shell(scriptData interface{}, inputs ...interface{}) (string, error) {
	var script string
	switch tt := scriptData.(type) {
	case string:
		script = tt
	case []byte:
		script = string(tt)
	default:
		return "", fmt.Errorf(
			"invalid template data, want string or bytes, got %T",
			tt,
		)
	}

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
	} else {
		stdin = os.Stdin
	}

	stdout := &bytes.Buffer{}
	runner, err := CreateEmbeddedShellRunner(
		ens.rc.WorkDir(), ens.rc, stdin, stdout, os.Stderr,
	)
	if err != nil {
		return "", err
	}

	err = RunScriptInEmbeddedShell(ens.rc, runner, syntax.NewParser(), script)
	return stdout.String(), err
}
