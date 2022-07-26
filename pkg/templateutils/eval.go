package templateutils

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
)

func createEvalNS(rc dukkha.RenderingContext) evalNS { return evalNS{rc: rc} }

type evalNS struct{ rc dukkha.RenderingContext }

// TODO: support writer as the second last argument
func (ns evalNS) Template(tplData String) (_ string, err error) {
	tplStr, err := toString(tplData)
	if err != nil {
		return
	}

	tpl, err := CreateTextTemplate(ns.rc).Parse(tplStr)
	if err != nil {
		err = fmt.Errorf("parse template %q: %w", tplStr, err)
		return
	}

	var buf strings.Builder
	err = tpl.Execute(&buf, ns.rc)
	if err != nil {
		err = fmt.Errorf("execute template %q: %w", tplStr, err)
		return
	}

	return buf.String(), nil
}

// Env expands environment variables in last argument, which can be string or bytes
//
// valid options before last argument are
// `disable_exec` / `enable_exec` to deny (default behvior) / allow shell script evaluation during expansion
func (ns evalNS) Env(args ...String) (_ string, err error) {
	var flags []string
	n := len(args)
	switch n {
	case 0:
		err = errAtLeastOneArgGotZero
		return
	case 1:
	default:
		flags, err = toStrings(args[:n-1])
		if err != nil {
			return
		}
	}

	text, err := toString(args[n-1])
	if err != nil {
		return
	}

	// TODO: better control over exec on/off in dukkha config
	enableExec := false
	for _, opt := range flags {
		switch opt {
		case "disable_exec":
			enableExec = false
		case "enable_exec":
			enableExec = true
		}
	}

	return ExpandEnv(ns.rc, text, enableExec)
}

// Shell runs script as bash script, optional inputs are data for stdin
// TODO: support writer as the second last argument
func (ns evalNS) Shell(script String, inputs ...Bytes) (ret cmdOutput, err error) {
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

	var stdout, stderr strings.Builder
	runner, err := CreateShellRunner(
		ns.rc.WorkDir(),
		ns.rc,
		stdin,
		&stdout,
		&stderr,
	)
	if err != nil {
		return
	}

	spt, err := toString(script)
	if err != nil {
		return
	}

	err = RunScript(
		ns.rc,
		runner,
		syntax.NewParser(),
		spt,
	)

	ret.Stdout = stdout.String()
	ret.Stderr = stderr.String()
	return
}
