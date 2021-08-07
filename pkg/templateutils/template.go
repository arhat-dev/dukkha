package templateutils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"arhat.dev/pkg/hashhelper"
	"arhat.dev/pkg/textquery"
	"github.com/Masterminds/sprig/v3"
	"github.com/hairyhenderson/gomplate/v3/funcs"
	"gopkg.in/yaml.v3"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
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
		Funcs(map[string]interface{}{
			"fromYaml": func(v string) interface{} {
				out := new(field.AnyObject)
				err := yaml.Unmarshal([]byte(v), out)
				if err != nil {
					panic(fmt.Errorf("failed to unmarshal yaml data\n\n%s\n\nerr: %w", v, err))
				}

				err = out.ResolveFields(rc, -1)
				if err != nil {
					panic(fmt.Errorf("failed to resolve yaml data\n\n%s\n\nerr: %w", v, err))
				}

				return out
			},
			"toYaml": func(v interface{}) string {
				data, _ := yaml.Marshal(v)
				return string(data)
			},
		}).
		// text functions
		Funcs(map[string]interface{}{
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

			"getAlpineArch": func(mArch string) string {
				v, _ := constant.GetAlpineArch(mArch)
				return v
			},
			"getAlpineTripleName": func(mArch string) string {
				v, _ := constant.GetAlpineTripleName(mArch)
				return v
			},

			"getDebianArch": func(mArch string) string {
				v, _ := constant.GetDebianArch(mArch)
				return v
			},
			"getDebianTripleName": func(mArch string, other ...string) string {
				targetKernel, targetLibc := "", ""
				if len(other) > 0 {
					targetKernel = other[0]
				}
				if len(other) > 1 {
					targetLibc = other[1]
				}

				v, _ := constant.GetDebianTripleName(mArch, targetKernel, targetLibc)
				return v
			},

			"getQemuArch": func(mArch string) string {
				v, _ := constant.GetQemuArch(mArch)
				return v
			},

			"getOciOS": func(mKernel string) string {
				v, _ := constant.GetOciOS(mKernel)
				return v
			},
			"getOciArch": func(mArch string) string {
				v, _ := constant.GetOciArch(mArch)
				return v
			},
			"getOciArchVariant": func(mArch string) string {
				v, _ := constant.GetOciArchVariant(mArch)
				return v
			},

			"getDockerOS": func(mKernel string) string {
				v, _ := constant.GetDockerOS(mKernel)
				return v
			},
			"getDockerArch": func(mArch string) string {
				v, _ := constant.GetDockerArch(mArch)
				return v
			},
			"getDockerArchVariant": func(mArch string) string {
				v, _ := constant.GetDockerArchVariant(mArch)
				return v
			},

			"getDockerHubArch": func(mArch string, other ...string) string {
				mKernel := ""
				if len(other) > 0 {
					mKernel = other[0]
				}

				v, _ := constant.GetDockerHubArch(mArch, mKernel)
				return v
			},
			"getDockerPlatformArch": func(mArch string) string {
				arch, ok := constant.GetDockerArch(mArch)
				if !ok {
					return ""
				}

				variant, _ := constant.GetDockerArchVariant(mArch)
				if len(variant) != 0 {
					return arch + "/" + variant
				}

				return arch
			},

			"getGolangOS": func(mKernel string) string {
				v, _ := constant.GetGolangOS(mKernel)
				return v
			},
			"getGolangArch": func(mArch string) string {
				v, _ := constant.GetGolangArch(mArch)
				return v
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
		}).
		Funcs(fm)
}
