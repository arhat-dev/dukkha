package templateutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"mvdan.cc/sh/v3/syntax"
)

func createEvalNS(rc dukkha.RenderingContext) *_evalNS {
	return &_evalNS{rc: rc}
}

type _evalNS struct {
	rc dukkha.RenderingContext
}

func (ens *_evalNS) Template(tplData interface{}) (string, error) {
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
			"failed to parse template to eval \n\n%s\n\n %w",
			tplStr, err,
		)
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, ens.rc)
	if err != nil {
		return "", fmt.Errorf(
			"failed to eval template \n\n%s\n\n %w",
			tplStr, err,
		)
	}

	return string(buf.Next(buf.Len())), nil
}

// Env expands environment variables in textData
// textData can be string or bytes
func (ens *_evalNS) Env(textData interface{}) (string, error) {
	var toExpand string
	switch tt := textData.(type) {
	case string:
		toExpand = tt
	case []byte:
		toExpand = string(tt)
	default:
		return "", fmt.Errorf("invalid non text data for expansion, got %T", tt)
	}

	return ExpandEnv(ens.rc, toExpand)
}

// Shell evaluates scriptData as bash script, inputs as stdin streams (if any)
// if no input was provided, it will use os.Stdin as shell stdin
// scriptData can be string or bytes
// inputs can be string, bytes, or io.Reader
func (ens *_evalNS) Shell(scriptData interface{}, inputs ...interface{}) (string, error) {
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
		ens.rc.WorkingDir(), ens.rc, stdin, stdout, os.Stderr,
	)
	if err != nil {
		return "", err
	}

	err = RunScriptInEmbeddedShell(ens.rc, runner, syntax.NewParser(), script)
	return stdout.String(), err
}
