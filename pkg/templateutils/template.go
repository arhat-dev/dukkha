package templateutils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"arhat.dev/pkg/md5helper"
	"arhat.dev/pkg/textquery"
	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
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

	return template.New("template").
		// template func from sprig
		Funcs(sprig.TxtFuncMap()).
		// template func from gomplate
		Funcs(funcs.CreateNetFuncs(rc)).
		Funcs(funcs.CreateReFuncs(rc)).
		Funcs(funcs.CreateStringFuncs(rc)).
		Funcs(funcs.CreateConvFuncs(rc)).
		Funcs(funcs.CreateTimeFuncs(rc)).
		Funcs(funcs.CreateMathFuncs(rc)).
		Funcs(funcs.CreateCryptoFuncs(rc)).
		Funcs(funcs.CreateFileFuncs(rc)).
		Funcs(funcs.CreateFilePathFuncs(rc)).
		Funcs(funcs.CreatePathFuncs(rc)).
		Funcs(funcs.CreateSockaddrFuncs(rc)).
		Funcs(funcs.CreateCollFuncs(rc)).
		Funcs(funcs.CreateUUIDFuncs(rc)).
		Funcs(funcs.CreateRandomFuncs(rc)).
		Funcs(map[string]interface{}{
			"strconv": func() *_strconvNS {
				return strconvNS
			},
			"dukkha": func() *dukkhaNS {
				return createDukkhaNS(rc)
			},
			"os": func() *_osNS {
				return osNS
			},
			"archconv": func() *_archconvNS {
				return archconvNS
			},
		}).
		// run shell commands in template
		Funcs(map[string]interface{}{
			"shell": func(script string, inputs ...string) (string, error) {
				var stdin io.Reader
				if len(inputs) != 0 {
					var readers []io.Reader
					for _, in := range inputs {
						readers = append(readers, strings.NewReader(in))
					}

					stdin = io.MultiReader(readers...)
				} else {
					stdin = os.Stdin
				}

				stdout := &bytes.Buffer{}
				runner, err := CreateEmbeddedShellRunner(
					rc.WorkingDir(), rc, stdin, stdout, os.Stderr,
				)
				if err != nil {
					return "", err
				}

				err = RunScriptInEmbeddedShell(rc, runner, syntax.NewParser(), script)
				return stdout.String(), err
			},
		}).
		// text processing
		Funcs(map[string]interface{}{
			"jq":      textquery.JQ,
			"jqBytes": textquery.JQBytes,
			"yq":      textquery.YQ,
			"yqBytes": textquery.YQBytes,

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

			"appendFile": func(filename string, data []byte) error {
				f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
				if err != nil {
					return err
				}

				_, err = f.Write(data)
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
		Funcs(fm)
}
