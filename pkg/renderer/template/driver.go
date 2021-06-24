package template

import (
	"bytes"
	"fmt"
	"text/template"

	"arhat.dev/pkg/textquery"
	"github.com/Masterminds/sprig/v3"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "template"

var _ renderer.Interface = (*Driver)(nil)

type Driver struct{}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(ctx *field.RenderingContext, tplStr string) (string, error) {
	tpl, err := template.New("template").
		Funcs(sprig.TxtFuncMap()).
		Funcs(map[string]interface{}{
			"jq":              textquery.JQ,
			"jqBytes":         textquery.JQBytes,
			"getAlpineArch":   constant.GetAlpineArch,
			"getAlpineTriple": constant.GetAlpineTripleName,
			"getDebianArch":   constant.GetDebianArch,
			"getDebianTriple": constant.GetDebianTripleName,
		}).
		Parse(tplStr)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to parse template \n\n%s\n\n %w", DefaultName, tplStr, err)
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, ctx.Values())
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to execute template \n\n%s\n\n %w", DefaultName, tplStr, err)
	}

	return buf.String(), nil
}
