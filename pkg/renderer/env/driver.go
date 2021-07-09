package env

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"arhat.dev/pkg/envhelper"
	"go.uber.org/multierr"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/utils"
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
		}

		ctx.AddRenderer(
			rendererName, &driver{
				getExecSpec: allShells[shellName].GetExecSpec,
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

	var err error

	result := &bytes.Buffer{}
	endAt := 0
	_ = envhelper.Expand(toExpand,
		createEnvExpandFunc(
			toExpand, result,
			func(at int) { endAt = at },
			func(name, origin string, at int) string {
				endAt = at
				v, ok := rc.Env()[name]
				if ok {
					return v
				}

				// env not found

				err = multierr.Append(err,
					fmt.Errorf("env %q not found", origin),
				)

				return origin
			},
			func(script string, err2 error, at int) string {
				endAt = at
				if err2 != nil {
					err = multierr.Append(err, err2)
					return ""
				}

				buf := &bytes.Buffer{}
				err = multierr.Append(err,
					renderer.RunShellScript(
						rc, script, false, buf, d.getExecSpec,
					),
				)

				return strings.TrimRightFunc(buf.String(), unicode.IsSpace)
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", DefaultName, err)
	}

	result.WriteString(toExpand[endAt:])

	return result.Bytes(), nil
}

func createEnvExpandFunc(
	toExpand string,
	result *bytes.Buffer,
	handleUpdate func(at int),
	handleEnv func(name, origin string, at int) string,
	handleExec func(script string, err error, at int) string,
) func(varName, origin string) string {
	lastAt := 0
	return func(varName, origin string) string {
		thisIdx := strings.Index(toExpand[lastAt:], origin)
		if thisIdx < 0 {
			// shell evaluation overrides the default behavior
			// of name matching
			return origin
		}

		if lastAt+thisIdx > 0 && toExpand[lastAt+thisIdx-1] == '\\' {
			// do not expand escaped env `\${SOME_ENV}`
			result.WriteString(toExpand[lastAt : lastAt+thisIdx-1])
			result.WriteString(origin)
			lastAt = lastAt + thisIdx + len(origin)

			handleUpdate(lastAt)

			return origin
		}

		result.WriteString(toExpand[lastAt : lastAt+thisIdx])

		lastAt += thisIdx

		if strings.HasPrefix(origin, "$(") {
			shellEval, err := utils.ParseBrackets(toExpand[lastAt+2:])
			if err != nil {
				lastAt += len(origin)
			} else {
				// 3 -> '$', '(' , ')'
				lastAt += len(shellEval) + 3
			}

			result.WriteString(handleExec(shellEval, err, lastAt))
			return origin
		}

		// no special handling
		lastAt += len(origin)

		result.WriteString(handleEnv(varName, origin, lastAt))
		return origin
	}
}
