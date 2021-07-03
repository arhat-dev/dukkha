package env

import (
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

type Driver struct{}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(ctx *field.RenderingContext, rawData interface{}) (string, error) {
	dataBytes, err := renderer.ToYamlBytes(rawData)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: unsupported input type %T: %w", DefaultName, rawData, err)
	}

	ret := envhelper.Expand(string(dataBytes), func(varName, origin string) string {
		if strings.HasPrefix(origin, "$(") {
			err = multierr.Append(
				err,
				fmt.Errorf("renderer:%s: shell evaluation is not supported", DefaultName),
			)
			return origin
		}

		v, ok := ctx.Values().Env[varName]
		if ok {
			return v
		}

		// env not found

		err = multierr.Append(
			err,
			fmt.Errorf("renderer:%s: env %q not found", DefaultName, origin),
		)

		return origin
	})

	return ret, nil
}
