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
			"jq":      textquery.JQ,
			"jqBytes": textquery.JQBytes,

			"getAlpineArch":       constant.GetAlpineArch,
			"getAlpineTripleName": constant.GetAlpineTripleName,

			"getDebianArch": constant.GetDebianArch,
			"getDebianTripleName": func(mArch string, other ...string) string {
				targetKernel, targetLibc := "", ""
				if len(other) > 0 {
					targetKernel = other[0]
				}
				if len(other) > 1 {
					targetLibc = other[1]
				}

				return constant.GetDebianTripleName(mArch, targetKernel, targetLibc)
			},

			"getQemuArch": constant.GetQemuArch,

			"getOciOS":          constant.GetOciOS,
			"getOciArch":        constant.GetOciArch,
			"getOciArchVariant": constant.GetOciArchVariant,

			"getDockerOS":          constant.GetDockerOS,
			"getDockerArch":        constant.GetDockerArch,
			"getDockerArchVariant": constant.GetDockerArchVariant,

			"getDockerHubArch": func(mArch string, other ...string) string {
				mKernel := ""
				if len(other) > 0 {
					mKernel = other[0]
				}

				return constant.GetDockerHubArch(mArch, mKernel)
			},
			"getDockerPlatformArch": func(mArch string) string {
				arch := constant.GetDockerArch(mArch)
				variant := constant.GetDockerArchVariant(mArch)
				if len(variant) != 0 {
					return arch + "/" + variant
				}

				return arch
			},

			"getGolangArch": constant.GetGolangArch,
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
