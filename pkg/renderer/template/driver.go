package template

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"arhat.dev/pkg/hashhelper"
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

func (d *Driver) Render(ctx *field.RenderingContext, rawData interface{}) (string, error) {
	tplBytes, err := renderer.ToYamlBytes(rawData)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: unsupported input type %T: %w", DefaultName, rawData, err)
	}

	tplStr := string(tplBytes)
	tpl, err := newTemplate().Parse(tplStr)
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

func newTemplate() *template.Template {
	return template.New("template").
		Funcs(sprig.TxtFuncMap()).
		Funcs(map[string]interface{}{
			"md5sum": func(s string) string {
				return hex.EncodeToString(hashhelper.MD5Sum([]byte(s)))
			},
			"os_ReadFile": func(filename string) (string, error) {
				data, err := os.ReadFile(filename)
				if err != nil {
					return "", err
				}

				return string(data), nil
			},

			"filepath_Join": func(parts ...string) string {
				return filepath.Join(parts...)
			},

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
		})
}
