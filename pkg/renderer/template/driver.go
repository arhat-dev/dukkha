package template

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "template"

var _ renderer.Interface = (*Driver)(nil)

type Driver struct{}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(ctx context.Context, rawValue string) (string, error) {
	tpl, err := template.New("").Parse(rawValue)
	if err != nil {
		return "", fmt.Errorf("renderer.template: failed to parse template: %w", err)
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, ctx.Value(constant.ContextKeyEnvironment))
	if err != nil {
		return "", fmt.Errorf("renderer.template: failed to execute template: %w", err)
	}

	return buf.String(), nil
}
