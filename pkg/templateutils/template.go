package templateutils

import (
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
)

type TemplateFuncFactory func(rc dukkha.RenderingContext) interface{}

var (
	toolSpecificTemplateFuncs = make(map[string]TemplateFuncFactory)
)

func RegisterTemplateFuncs(fm map[string]TemplateFuncFactory) {
	for k, f := range fm {
		if _, ok := toolSpecificTemplateFuncs[k]; ok {
			panic(fmt.Sprintf("func %q already registered", k))
		}

		toolSpecificTemplateFuncs[k] = f
	}
}

func CreateTemplate(rc dukkha.RenderingContext) *template.Template {
	fm := make(map[string]interface{})
	for k, createTemplateFunc := range toolSpecificTemplateFuncs {
		fm[k] = createTemplateFunc(rc)
	}

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
			"os_WriteFile": func(filename string, data []byte) error {
				return os.WriteFile(filename, data, 0640)
			},
			"appendFile": func(filename string, data []byte) error {
				f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
				if err != nil {
					return err
				}

				_, err = f.Write(data)
				return err
			},
			"toBytes": func(s string) []byte {
				return []byte(s)
			},

			"filepath_Join": filepath.Join,

			"jq":      textquery.JQ,
			"jqBytes": textquery.JQBytes,
			"yq":      textquery.YQ,
			"yqBytes": textquery.YQBytes,

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
		Funcs(fm)
}
