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
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/tools/buildah"
)

const DefaultName = "template"

func New() dukkha.Renderer {
	return &driver{}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct{}

func (d *driver) Name() string { return DefaultName }

func (d *driver) RenderYaml(rc dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	tplBytes, err := renderer.ToYamlBytes(rawData)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: unsupported input type %T: %w", DefaultName, rawData, err)
	}

	tplStr := string(tplBytes)
	tpl, err := newTemplate(rc).Parse(tplStr)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: failed to parse template \n\n%s\n\n %w", DefaultName, tplStr, err)
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, rc)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: failed to execute template \n\n%s\n\n %w", DefaultName, tplStr, err)
	}

	return buf.Bytes(), nil
}

func newTemplate(rc dukkha.RenderingContext) *template.Template {
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

			"filepath_Join": filepath.Join,

			"getBuildahImageIDFile": func(imageName string) string {
				return buildah.GetImageIDFileForImageName(
					rc.CacheDir(),
					buildah.SetDefaultImageTagIfNoTagSet(rc, imageName),
				)
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
