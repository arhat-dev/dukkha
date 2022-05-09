package templateutils

import (
	"encoding/hex"
	"fmt"
	"os"

	"arhat.dev/pkg/md5helper"
	"arhat.dev/pkg/textquery"
	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/third_party/golang/text/template"
	"arhat.dev/dukkha/third_party/gomplate/funcs"
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
		Funcs(map[string]interface{}{
			"strings": func() stringsNS { return stringsNS{} },

			"replaceAll": stringsNS{}.ReplaceAll,
			"title":      stringsNS{}.Title,
			"toUpper":    stringsNS{}.ToUpper,
			"toLower":    stringsNS{}.ToLower,
			"trimSpace":  stringsNS{}.TrimSpace,
			"indent":     stringsNS{}.Indent,
			"quote":      stringsNS{}.Quote,
			"shellQuote": stringsNS{}.ShellQuote,
			"squote":     stringsNS{}.Squote,

			"contains":  stringsNS{}.Contains,
			"hasPrefix": stringsNS{}.HasPrefix,
			"hasSuffix": stringsNS{}.HasSuffix,
			"split":     stringsNS{}.Split,
			"splitN":    stringsNS{}.SplitN,
			"trim":      stringsNS{}.Trim,

			"kebabcase": stringsNS{}.KebabCase,
			"snakecase": stringsNS{}.SnakeCase,
			"camelcase": stringsNS{}.CamelCase,
		}).
		Funcs(map[string]interface{}{
			"filepath": func() filepathNS { return createFilePathNS(rc) },
			"strconv":  func() strconvNS { return strconvNS{} },
			"dukkha":   func() dukkhaNS { return createDukkhaNS(rc) },
			"os":       func() osNS { return createOSNS(rc) },
			"archconv": func() archconvNS { return archconvNS{} },
			"git":      rc.GitValues,
			"host":     rc.HostValues,
			// eval shell and template
			"eval":   func() evalNS { return createEvalNS(rc) },
			"env":    rc.Env,
			"values": rc.Values,
			"matrix": func() map[string]string { return rc.MatrixFilter().AsEntry() },
			// state task execution
			"state": func() stateNS { return createStateNS(rc) },
			// for transform renderer
			"VALUE": func() interface{} {
				vg, ok := rc.(di.VALUEGetter)
				if ok {
					return vg.VALUE()
				}

				return nil
			},
		}).
		// text processing
		Funcs(map[string]interface{}{
			"jq":       textquery.JQ,
			"jqBytes":  textquery.JQBytes,
			"jqObject": jqObject,
			"yq":       textquery.YQ,
			"yqBytes":  textquery.YQBytes,

			"fromYaml": func(v string) interface{} {
				ret, err := fromYaml(rc, v)
				if err != nil {
					panic(err)
				}
				return ret
			},
			"toYaml": func(v interface{}) string {
				data, _ := yaml.Marshal(v)
				return string(data)
			},

			"addPrefix": func(args ...string) string {
				sep := "\n"
				if len(args) == 3 {
					sep = args[0]
				}

				return AddPrefix(args[len(args)-1], args[len(args)-2], sep)
			},
			"removePrefix": func(args ...string) string {
				sep := "\n"
				if len(args) == 3 {
					sep = args[0]
				}

				return RemovePrefix(args[len(args)-1], args[len(args)-2], sep)
			},
			"addSuffix": func(args ...string) string {
				sep := "\n"
				if len(args) == 3 {
					sep = args[0]
				}

				return AddSuffix(args[len(args)-1], args[len(args)-2], sep)
			},
			"removeSuffix": func(args ...string) string {
				sep := "\n"
				if len(args) == 3 {
					sep = args[0]
				}

				return RemoveSuffix(args[len(args)-1], args[len(args)-2], sep)
			},
		}).
		// dukkha specific template func
		Funcs(map[string]interface{}{
			"md5sum": func(s string) string {
				return hex.EncodeToString(md5helper.Sum([]byte(s)))
			},

			"totp": totpTemplateFunc,

			"appendFile": func(filename string, data []byte) error {
				f, err := rc.FS().OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
				if err != nil {
					return err
				}

				_, err = f.(*os.File).Write(data)
				return err
			},
			"toBytes": func(s interface{}) ([]byte, error) {
				switch dt := s.(type) {
				case string:
					return []byte(dt), nil
				case []byte:
					return dt, nil
				case []rune:
					return []byte(string(dt)), nil
				default:
					return nil, fmt.Errorf(
						"invalid non string, bytes, nor runes: %T", s,
					)
				}
			},

			"setDefaultImageTag": func(imageName string, flags ...string) string {
				keepKernelInfo := false
				for _, f := range flags {
					if f == "keepKernelInfo" {
						keepKernelInfo = true
					}
				}
				return SetDefaultImageTagIfNoTagSet(rc, imageName, keepKernelInfo)
			},
			"setDefaultManifestTag": func(imageName string, flags ...string) string {
				return SetDefaultManifestTagIfNoTagSet(rc, imageName)
			},

			"getDefaultImageTag": func(imageName string, flags ...string) string {
				keepKernelInfo := false
				for _, f := range flags {
					if f == "keepKernelInfo" {
						keepKernelInfo = true
					}
				}
				return GetDefaultImageTag(rc, imageName, keepKernelInfo)
			},
			"getDefaultManifestTag": func(imageName string, flags ...string) string {
				return GetDefaultManifestTag(rc, imageName)
			},
		}).
		Funcs(fm).
		// placeholder functions to be overridden before Execute
		Funcs(map[string]interface{}{
			"var": func() map[string]interface{} { return nil },
			// include like helm include
			"include": func(name string, data interface{}) (string, error) {
				return "", fmt.Errorf("no implementation")
			},
		})
}
