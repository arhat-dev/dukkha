package env

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"go.uber.org/multierr"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

// nolint:revive
const (
	DefaultName = "env"
)

func init() {
	dukkha.RegisterRenderer(
		DefaultName,
		func() dukkha.Renderer {
			return NewDefault(nil)
		},
	)
}

func NewDefault(getExecSpec dukkha.ExecSpecGetFunc) dukkha.Renderer {
	return &driver{getExecSpec: getExecSpec}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	field.BaseField

	getExecSpec dukkha.ExecSpecGetFunc
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	allShells := ctx.AllShells()
	for shellName := range allShells {
		rendererName := DefaultName
		if len(shellName) != 0 {
			rendererName += ":" + shellName
			ctx.AddRenderer(
				rendererName, &driver{
					getExecSpec: allShells[shellName].GetExecSpec,
				},
			)
			continue
		}

		ctx.AddRenderer(
			rendererName, &driver{
				getExecSpec: nil,
			},
		)
	}

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
		return nil, fmt.Errorf("renderer.%s: invalid expansion text: %w", DefaultName, err)
	}

	var errEnvMissing error
	environ := expand.FuncEnviron(func(name string) string {
		v, ok := rc.Env()[name]
		if ok {
			return v
		}

		switch name {
		case "IFS":
			return " \t\n"
		case "OPTIND":
			return "1"
		case "PWD":
			return rc.WorkingDir()
		case "UID":
			// os.Getenv("UID") usually retruns empty value
			// so we have to call os.Getuid
			return strconv.FormatInt(int64(os.Getuid()), 10)
		// case "GID":
		default:
			errEnvMissing = multierr.Append(errEnvMissing,
				fmt.Errorf("env %q not found", name),
			)
		}

		return ""
	})

	embeddedShellOutput := &bytes.Buffer{}
	runner, err := interp.New(
		interp.Env(environ),
		interp.Dir(rc.WorkingDir()),
		interp.StdIO(nil, embeddedShellOutput, ioutil.Discard),
		interp.ExecHandler(interp.DefaultExecHandler(0)),
	)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: failed to create runner for env: %w", DefaultName, err)
	}

	printer := syntax.NewPrinter(
		syntax.FunctionNextLine(false),
		syntax.Indent(2),
	)

	result, err := expand.Document(&expand.Config{
		Env: environ,
		CmdSubst: func(w io.Writer, cs *syntax.CmdSubst) error {
			buf := &bytes.Buffer{}
			err := printer.Print(buf, cs)
			if err != nil {
				return err
			}

			script := string(buf.Bytes()[2 : buf.Len()-1])

			if d.getExecSpec == nil {
				embeddedShellOutput.Reset()
				f, err := parser.Parse(strings.NewReader(script), "")
				if err != nil {
					return fmt.Errorf(
						"failed to parse shell evaluation \n\n%s\n\nusing embedded shell: %w",
						buf.String(),
						err,
					)
				}

				err = runner.Run(rc, f)
				if err != nil {
					return fmt.Errorf(
						"failed to evaluate command \n\n%s\n\nusing embedded shell: %w",
						script, err,
					)
				}

				_, err = embeddedShellOutput.WriteTo(w)
				if err != nil {
					return fmt.Errorf(
						"failed to write embedded shell output to result value: %w", err,
					)
				}

				return nil
			}

			return renderer.RunShellScript(
				rc, script, false, w, d.getExecSpec,
			)
		},
		ProcSubst: nil,
		ReadDir:   ioutil.ReadDir,
		GlobStar:  true,
		NullGlob:  true,
		NoUnset:   false,
	},
		word,
	)

	if errEnvMissing != nil {
		return nil, fmt.Errorf("renderer.%s: some env not resolved: %w", DefaultName, errEnvMissing)
	}

	if err != nil {
		return nil, fmt.Errorf("renderer.%s: env expansion failed: %w", DefaultName, err)
	}

	return []byte(result), nil
}
