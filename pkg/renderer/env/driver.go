package env

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/templateutils"
)

// nolint:revive
const (
	DefaultName = "env"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault() dukkha.Renderer {
	return &driver{}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	rs.BaseField
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	return nil
}

func (d *driver) RenderYaml(rc dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	var toExpand string

	switch t := rawData.(type) {
	case string:
		toExpand = t
	case []byte:
		toExpand = string(t)
	default:
		dataBytes, err := renderer.ToYamlBytes(rawData)
		if err != nil {
			return nil, fmt.Errorf("renderer.%s: unsupported input type %T: %w", DefaultName, rawData, err)
		}
		toExpand = string(dataBytes)
	}

	parser := syntax.NewParser(
		syntax.Variant(syntax.LangBash),
	)

	word, err := parser.Document(strings.NewReader(toExpand))
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: invalid expansion text %q: %w",
			DefaultName, toExpand, err,
		)
	}

	embeddedShellOutput := &bytes.Buffer{}
	runner, err := templateutils.CreateEmbeddedShellRunner(
		rc.WorkingDir(), rc, nil, embeddedShellOutput, os.Stderr,
	)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: failed to create runner for env: %w", DefaultName, err)
	}

	printer := syntax.NewPrinter(
		syntax.FunctionNextLine(false),
		syntax.Indent(2),
	)

	result, err := expand.Document(&expand.Config{
		Env: rc,
		CmdSubst: func(w io.Writer, cs *syntax.CmdSubst) error {
			buf := &bytes.Buffer{}
			err2 := printer.Print(buf, cs)
			if err2 != nil {
				return fmt.Errorf("failed to get evaluation commands: %w", err2)
			}

			script := string(buf.Bytes()[2 : buf.Len()-1])

			embeddedShellOutput.Reset()
			err2 = templateutils.RunScriptInEmbeddedShell(rc, runner, parser, script)
			if err2 != nil {
				return err2
			}

			_, err2 = embeddedShellOutput.WriteTo(w)
			if err2 != nil {
				return fmt.Errorf(
					"failed to write embedded shell output to result value: %w", err,
				)
			}

			return nil
		},
		ProcSubst: nil,
		ReadDir:   ioutil.ReadDir,
		GlobStar:  true,
		NullGlob:  true,
		NoUnset:   true,
	},
		word,
	)

	if err != nil {
		return nil, fmt.Errorf("renderer.%s: env expansion failed: %w", DefaultName, err)
	}

	return []byte(result), nil
}
