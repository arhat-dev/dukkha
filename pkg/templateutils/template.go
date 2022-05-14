package templateutils

import (
	"encoding/hex"
	"fmt"

	"arhat.dev/pkg/md5helper"
	"arhat.dev/pkg/stringhelper"
	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/third_party/golang/text/template"
	"arhat.dev/dukkha/third_party/gomplate/funcs"
)

type TemplateFuncFactory func(rc dukkha.RenderingContext) any

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
	fm := make(map[string]any)
	for k, createTemplateFunc := range toolSpecificTemplateFuncs {
		fm[k] = createTemplateFunc(rc)
	}

	var (
		strNS stringsNS
	)

	return template.New("tpl").
		// template func from sprig
		Funcs(template.FuncMap(sprig.TxtFuncMap())).
		// template func from gomplate
		Funcs(funcs.CreateNetFuncs(rc)).
		Funcs(funcs.CreateReFuncs(rc)).
		Funcs(funcs.CreateConvFuncs(rc)).
		Funcs(funcs.CreateTimeFuncs(rc)).
		Funcs(funcs.CreateMathFuncs(rc)).
		Funcs(funcs.CreateCryptoFuncs(rc)).
		Funcs(funcs.CreateFileFuncs(rc)).
		Funcs(funcs.CreatePathFuncs(rc)).
		Funcs(funcs.CreateSockaddrFuncs(rc)).
		Funcs(funcs.CreateCollFuncs(rc)).
		Funcs(funcs.CreateUUIDFuncs(rc)).
		Funcs(funcs.CreateRandomFuncs(rc)).
		Funcs(map[string]any{
			"strings": func() stringsNS { return strNS },

			"replaceAll": strNS.ReplaceAll,
			"title":      strNS.Title,
			"toUpper":    strNS.ToUpper,
			"toLower":    strNS.ToLower,
			"trimSpace":  strNS.TrimSpace,
			"indent":     strNS.Indent,
			"quote":      strNS.Quote,
			"shellQuote": strNS.ShellQuote,
			"squote":     strNS.Squote,

			"contains":  strNS.Contains,
			"hasPrefix": strNS.HasPrefix,
			"hasSuffix": strNS.HasSuffix,
			"split":     strNS.Split,
			"splitN":    strNS.SplitN,
			"trim":      strNS.Trim,

			"kebabcase": strNS.KebabCase,
			"snakecase": strNS.SnakeCase,
			"camelcase": strNS.CamelCase,

			"jq":    strNS.JQ,
			"yq":    strNS.YQ,
			"jqObj": strNS.JQObj,

			// multi-line string

			"addPrefix": func(args ...String) string {
				sep := "\n"
				if len(args) == 3 {
					sep = toString(args[0])
				}

				return strNS.AddPrefix(toString(args[len(args)-1]), toString(args[len(args)-2]), sep)
			},
			"removePrefix": func(args ...String) string {
				sep := "\n"
				if len(args) == 3 {
					sep = toString(args[0])
				}

				return strNS.RemovePrefix(toString(args[len(args)-1]), toString(args[len(args)-2]), sep)
			},
			"addSuffix": func(args ...String) string {
				sep := "\n"
				if len(args) == 3 {
					sep = toString(args[0])
				}

				return strNS.AddSuffix(toString(args[len(args)-1]), toString(args[len(args)-2]), sep)
			},
			"removeSuffix": func(args ...String) string {
				sep := "\n"
				if len(args) == 3 {
					sep = toString(args[0])
				}

				return strNS.RemoveSuffix(toString(args[len(args)-1]), toString(args[len(args)-2]), sep)
			},
		}).
		Funcs(map[string]any{
			"filepath": func() filepathNS { return createFilePathNS(rc) },
			"strconv":  func() strconvNS { return strconvNS{} },
			"dukkha":   func() dukkhaNS { return createDukkhaNS(rc) },
			"os":       func() osNS { return createOSNS(rc) },
			"archconv": func() archconvNS { return archconvNS{} },
			"git":      rc.GitValues,  // git.{tag, branch ...}
			"host":     rc.HostValues, // host.{arch, arch_simple, kernel ...}
			// eval shell and template
			"eval":   func() evalNS { return createEvalNS(rc) },
			"env":    rc.Env,
			"values": rc.Values,
			"matrix": func() map[string]string {
				mf := rc.MatrixFilter()
				return mf.AsEntry()
			},
			// state task execution
			"state": func() stateNS { return createStateNS(rc) },
			// for transform renderer
			"VALUE": func() any {
				vg, ok := rc.(di.VALUEGetter)
				if ok {
					return vg.VALUE()
				}

				return nil
			},
		}).
		// yaml processing
		Funcs(map[string]any{
			"fromYaml": func(v String) any {
				ret, err := fromYaml(rc, toString(v))
				if err != nil {
					panic(err)
				}
				return ret
			},
			"toYaml": func(v any) string {
				data, _ := yaml.Marshal(v)
				return stringhelper.Convert[string, byte](data)
			},
		}).
		// dukkha specific template func
		Funcs(map[string]any{
			"md5sum": func(s Bytes) string {
				return hex.EncodeToString(md5helper.Sum(toBytes(s)))
			},

			"totp":    totpTemplateFunc,
			"toBytes": func(s any) []byte { return toBytes(s) },

			"setDefaultImageTag": func(imageName String, flags ...String) string {
				keepKernelInfo := false
				for _, f := range flags {
					if toString(f) == "keepKernelInfo" {
						keepKernelInfo = true
					}
				}
				return SetDefaultImageTagIfNoTagSet(rc, toString(imageName), keepKernelInfo)
			},
			"setDefaultManifestTag": func(imageName String, flags ...String) string {
				return SetDefaultManifestTagIfNoTagSet(rc, toString(imageName))
			},

			"getDefaultImageTag": func(imageName String, flags ...String) string {
				keepKernelInfo := false
				for _, f := range flags {
					if toString(f) == "keepKernelInfo" {
						keepKernelInfo = true
					}
				}
				return GetDefaultImageTag(rc, toString(imageName), keepKernelInfo)
			},
			"getDefaultManifestTag": func(imageName String, flags ...String) string {
				return GetDefaultManifestTag(rc, toString(imageName))
			},
		}).
		Funcs(fm).
		// placeholder functions to be overridden before template.Execute
		Funcs(map[string]any{
			"var": func() map[string]any { return nil },
			// include like helm include
			"include": func(name string, data any) (string, error) {
				return "", fmt.Errorf("no implementation")
			},
		})
}
