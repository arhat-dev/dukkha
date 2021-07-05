package env

import (
	"bytes"
	"fmt"
	"strings"

	"arhat.dev/pkg/envhelper"
	"go.uber.org/multierr"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

// nolint:revive
const (
	DefaultName = "env"
)

var _ renderer.Interface = (*Driver)(nil)

type Driver struct {
	getExecSpec field.ExecSpecGetFunc
}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(ctx *field.RenderingContext, rawData interface{}) (string, error) {
	var toExpand string

	switch t := rawData.(type) {
	case string:
		toExpand = t
	case []byte:
		toExpand = string(t)
	default:
		dataBytes, err := renderer.ToYamlBytes(rawData)
		if err != nil {
			return "", fmt.Errorf("renderer.%s: unsupported input type %T: %w", DefaultName, rawData, err)
		}
		toExpand = string(dataBytes)
	}

	var err error

	buf := &bytes.Buffer{}
	result := &bytes.Buffer{}
	endAt := 0
	_ = envhelper.Expand(toExpand,
		createEnvExpandFunc(
			toExpand, result,
			func(name, origin string, at int) string {
				endAt = at
				v, ok := ctx.Values().Env[name]
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

				buf.Reset()

				err = multierr.Append(err,
					renderer.RunShellScript(ctx, script, false, buf, d.getExecSpec),
				)

				return buf.String()
			},
		),
	)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: %w", DefaultName, err)
	}

	result.WriteString(toExpand[endAt:])

	return result.String(), nil
}

func createEnvExpandFunc(
	toExpand string,
	result *bytes.Buffer,
	handleEnv func(name, origin string, at int) string,
	handleExec func(script string, err error, at int) string,
) func(varName, origin string) string {
	lastAt := 0
	return func(varName, origin string) string {
		thisIdx := strings.Index(toExpand[lastAt:], origin)
		if thisIdx < 0 {
			return ""
		}

		result.WriteString(toExpand[lastAt : lastAt+thisIdx])

		lastAt += thisIdx

		if strings.HasPrefix(origin, "$(") {
			shellEval, err := parseShellEval(toExpand[lastAt+2:])
			if err != nil {
				lastAt += len(origin)
			} else {
				lastAt += len(shellEval) + 3
			}

			result.WriteString(handleExec(shellEval, err, lastAt))
			return ""
		}

		// no special handling
		lastAt += len(origin)

		result.WriteString(handleEnv(varName, origin, lastAt))
		return ""
	}
}

func parseShellEval(toEpand string) (string, error) {
	leftBrackets := 0
	for i := range toEpand {
		switch toEpand[i] {
		case '(':
			leftBrackets++
		case ')':
			if leftBrackets == 0 {
				return toEpand[:i], nil
			}
			leftBrackets--
		}
	}

	// invalid data
	return "", fmt.Errorf("unexpected non-terminated brackets")
}
