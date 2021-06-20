package template

import (
	"bytes"
	"fmt"
	"text/template"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "template"

var _ renderer.Interface = (*Driver)(nil)

type Driver struct{}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(ctx *field.RenderingContext, tplStr string) (string, error) {
	tpl, err := template.New("").Parse(tplStr)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to parse template: %w", DefaultName, err)
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, ctx.Values())
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to execute template: %w", DefaultName, err)
	}

	return buf.String(), nil
}
