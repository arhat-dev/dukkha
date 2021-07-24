package shell

import (
	"bytes"
	"fmt"
	"os"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/shell"
)

const DefaultName = "shell"

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault() dukkha.Renderer {
	return &driver{}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	field.BaseField

	getExecSpec dukkha.ExecSpecGetFunc
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	allShells := ctx.AllShells()
	for shellName := range allShells {
		ctx.AddRenderer(
			DefaultName+":"+shellName,
			&driver{
				getExecSpec: allShells[shellName].GetExecSpec,
			},
		)
	}

	ctx.AddRenderer(DefaultName, NewDefault())

	return nil
}

func (d *driver) RenderYaml(rc dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	var scripts []string
	switch t := rawData.(type) {
	case string:
		scripts = append(scripts, t)
	case []byte:
		scripts = append(scripts, string(t))
	case []interface{}:
		for _, v := range t {
			scriptBytes, err := renderer.ToYamlBytes(v)
			if err != nil {
				return nil, fmt.Errorf("renderer.%s: unexpected list item type %T: %w", DefaultName, v, err)
			}

			scripts = append(scripts, string(scriptBytes))
		}
	default:
		return nil, fmt.Errorf("renderer.%s: unsupported input type %T", DefaultName, rawData)
	}

	buf := &bytes.Buffer{}

	var (
		parser *syntax.Parser
		runner *interp.Runner
		err    error
	)
	if d.getExecSpec == nil {
		runner, err = shell.CreateEmbeddedShellRunner(
			rc.WorkingDir(), rc, nil, buf, os.Stderr,
		)
		if err != nil {
			return nil, fmt.Errorf("renderer.%s: failed to create embedded shell: %w", DefaultName, err)
		}

		parser = syntax.NewParser(
			syntax.Variant(syntax.LangBash),
		)
	}

	for _, script := range scripts {
		if d.getExecSpec == nil {
			err = shell.RunScriptInEmbeddedShell(rc, runner, parser, script)
		} else {
			err = shell.RunScript(rc, script, false, buf, d.getExecSpec)
		}

		if err != nil {
			return nil, fmt.Errorf("renderer.%s: %w", DefaultName, err)
		}
	}

	return buf.Bytes(), nil
}
