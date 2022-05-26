package templateutils

import (
	"bytes"
	"fmt"
	"io"

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
func (ns evalNS) Env(args ...String) (_ string, err error) {
	var (
		textData String
		flags    []string
	)
	switch n := len(args); n {
	case 0:
		err = errAtLeastOneArgGotZero
		return
	case 1:
		textData = args[0]
	default:
		flags, err = toStrings(args[:n-1])
		if err != nil {
			return
		}

		textData = args[n-1]
	}

	enableExec := false
	for _, opt := range flags {
		switch opt {
		case "disable_exec":
			enableExec = false
		case "enable_exec":
			enableExec = true
		}
	}

	return ExpandEnv(ns.rc, must(toString(textData)), enableExec)
}

// Shell runs script as bash script, optional inputs are data for stdin
func (ns evalNS) Shell(script String, inputs ...Bytes) (_ string, err error) {
	var stdin io.Reader
	if len(inputs) != 0 {
		var readers []io.Reader
		for _, in := range inputs {
			inData, inReader, isReader, err2 := toBytesOrReader(in)
			if err2 != nil {
				err = err2
				return
			}

			if isReader {
				readers = append(readers, inReader)
			} else {
				readers = append(readers, bytes.NewReader(inData))
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
		ns.rc.Stderr(),
	)
	if err != nil {
		return
	}

	spt, err := toString(script)
	if err != nil {
		return
	}

	err = RunScriptInEmbeddedShell(
		ns.rc,
		runner,
		syntax.NewParser(),
		spt,
	)

	return stringhelper.Convert[string, byte](stdout.Next(stdout.Len())), err
}
