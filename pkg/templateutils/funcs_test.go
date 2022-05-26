package templateutils

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"text/template"

	di "arhat.dev/dukkha/internal"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"github.com/stretchr/testify/assert"

	_ "embed" // for go:embed
)

var (
	//go:embed gen/funcs.tpl
	funcs_template string
)

func TestGenerateFuncs(t *testing.T) {
	t.SkipNow()

	type Values struct {
		StaticFuncs    []TemplateFuncInfo
		LastStaticFunc TemplateFuncInfo

		ContextualFuncs    []TemplateFuncInfo
		LastContextualFunc TemplateFuncInfo

		PlaceholderFuncs    []TemplateFuncInfo
		LastPlaceholderFunc TemplateFuncInfo
	}

	var val Values

	val.StaticFuncs = collectTemplateFuncs(staticFuncMaps[:])
	val.LastStaticFunc = val.StaticFuncs[len(val.StaticFuncs)-1]

	val.ContextualFuncs = collectTemplateFuncs(contextualFuncs[:])
	val.LastContextualFunc = val.ContextualFuncs[len(val.ContextualFuncs)-1]

	val.PlaceholderFuncs = collectTemplateFuncs(placeholderFuncMaps[:])
	val.LastPlaceholderFunc = val.PlaceholderFuncs[len(val.PlaceholderFuncs)-1]

	println(len(val.PlaceholderFuncs))

	tpl, err := template.New("").Parse(funcs_template)
	assert.NoError(t, err)

	var out bytes.Buffer
	err = tpl.Execute(&out, val)
	assert.NoError(t, err)

	data, err := format.Source(out.Next(out.Len()))
	assert.NoError(t, err)

	err = os.WriteFile("funcs.go", data, 0644)
	assert.NoError(t, err)
}

type TemplateFuncInfo struct {
	namespace string
	name      string
	funcType  string
}

func (tf TemplateFuncInfo) Name() string {
	if len(tf.namespace) != 0 {
		return tf.namespace + "." + tf.name
	}

	return tf.name
}

func (tf TemplateFuncInfo) Ident() string {
	if len(tf.namespace) != 0 {
		return tf.namespace + "_" + tf.name
	}

	return tf.name
}

func collectTemplateFuncs(fms []map[string]any) []TemplateFuncInfo {
	visited := make(map[string]struct{})

	var ret []TemplateFuncInfo
	for _, fm := range fms {
		for k, fn := range fm {
			tfs := parseTemplateFunc(k, fn)

			for _, tf := range tfs {
				if _, ok := visited[tf.Ident()]; ok {
					continue
				}

				visited[tf.Ident()] = struct{}{}
				ret = append(ret, tf)
			}
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name() < ret[j].Name()
	})

	return ret
}

var replacer = strings.NewReplacer(
	"templateutils.String", "String",
	"templateutils.Bytes", "Bytes",
	"interface {}", "any",
)

func parseTemplateFunc(k string, fn any) []TemplateFuncInfo {
	v := reflect.ValueOf(fn)
	vt := reflect.TypeOf(fn)

	if vt.NumIn() != 0 {
		// is a func
		return []TemplateFuncInfo{{
			name:     k,
			funcType: replacer.Replace(vt.String()),
		}}
	}

	// namespaced func
	ns := v.Call(nil)[0].Type()

	switch ns.Kind() {
	case reflect.Map:
		return []TemplateFuncInfo{{
			name:     k,
			funcType: replacer.Replace(vt.String()),
		}}
	default:
		// is a namespace, see below
	}

	nMethods := ns.NumMethod()
	ret := make([]TemplateFuncInfo, 0, nMethods)
	for i := 0; i < nMethods; i++ {
		m := ns.Method(i)
		if len(m.PkgPath) != 0 {
			// unexported, ignore
			continue
		}

		ft := m.Func.Type()
		var (
			fin  []reflect.Type
			fout []reflect.Type
		)
		// skip first (receiver)
		for i := 1; i < ft.NumIn(); i++ {
			fin = append(fin, ft.In(i))
		}

		for i := 0; i < ft.NumOut(); i++ {
			fout = append(fout, ft.Out(i))
		}

		ret = append(ret, TemplateFuncInfo{
			namespace: k,
			name:      m.Name,
			funcType:  replacer.Replace(reflect.FuncOf(fin, fout, ft.IsVariadic()).String()),
		})
	}

	return ret
}

var (
	testRC = dukkha_test.NewTestContext(context.TODO())
)

var contextualFuncs = [...]map[string]any{
	{
		"os": func() osNS { return osNS{} },

		"fs":    func() fsNS { return createFSNS(testRC) },
		"touch": func(file String) (struct{}, error) { return fsNS{}.WriteFile(file) },
		"read":  fsNS{}.ReadFile,
		"write": fsNS{}.WriteFile,
		"mkdir": fsNS{}.Mkdir,
		"find":  fsNS{}.Find,

		"dukkha": func() dukkhaNS { return createDukkhaNS(testRC) },

		"git":  testRC.GitValues,  // git.{tag, branch ...}
		"host": testRC.HostValues, // host.{arch, arch_simple, kernel ...}
		// eval shell and template
		"env":    testRC.Env,
		"values": testRC.Values,
		"matrix": func() map[string]string {
			mf := testRC.MatrixFilter()
			return mf.AsEntry()
		},
		// state task execution
		"state": func() stateNS { return createStateNS(testRC) },
		// for transform renderer
		"VALUE": func() any {
			vg, ok := testRC.(di.VALUEGetter)
			if ok {
				return vg.VALUE()
			}

			return nil
		},

		"fromJson": func(v String) (any, error) { return nil, nil },
		"fromYaml": func(v String) (any, error) { return nil, nil },
	},

	{
		"setDefaultImageTag": func(imageName String, flags ...String) string {
			return ""
		},
		"setDefaultManifestTag": func(imageName String, flags ...String) string {
			return ""
		},

		"getDefaultImageTag": func(imageName String, flags ...String) string {
			return ""
		},
		"getDefaultManifestTag": func(imageName String, flags ...String) string {
			return ""
		},
	},
}

var staticFuncMaps = [...]map[string]any{
	{
		"close": close,

		"path": func() pathNS { return pathNS{} },
		"dns":  func() dnsNS { return dnsNS{} },
		"uuid": func() uuidNS { return uuidNS{} },
		"re":   func() regexpNS { return regexpNS{} },

		"hash": func() hashNS { return hashNS{} },

		"md5":    hashNS{}.MD5,
		"sha256": hashNS{}.SHA256,
		"sha512": hashNS{}.SHA512,

		"cred": func() credentialNS { return credentialNS{} },
		"totp": credentialNS{}.Totp,

		"time": func() timeNS { return timeNS{} },
		"now":  timeNS{}.Now,

		"enc":    func() encNS { return encNS{} },
		"base64": encNS{}.Base64,
		"hex":    encNS{}.Hex,
		"toJson": encNS{}.JSON,
		"toYaml": encNS{}.YAML,
	},
	{
		"text": func() stringsNS { return stringsNS{} },

		"replaceAll": stringsNS{}.ReplaceAll,
		"title":      stringsNS{}.Title,
		"trimSpace":  stringsNS{}.TrimSpace,
		"indent":     stringsNS{}.Indent,
		"quote":      stringsNS{}.DoubleQuote,
		"shellQuote": stringsNS{}.ShellQuote,
		"squote":     stringsNS{}.SingleQuote,

		"contains":  stringsNS{}.Contains,
		"hasPrefix": stringsNS{}.HasPrefix,
		"hasSuffix": stringsNS{}.HasSuffix,
		"split":     stringsNS{}.Split,
		"splitN":    stringsNS{}.SplitN,
		"trim":      stringsNS{}.Trim,

		"kebabcase": stringsNS{}.KebabCase,
		"snakecase": stringsNS{}.SnakeCase,
		"camelcase": stringsNS{}.CamelCase,

		// "jq": textNS{}.JQ,
		// "yq": textNS{}.YQ,

		// multi-line string

		"addPrefix": func(args ...String) string {
			return ""
		},
		"removePrefix": func(args ...String) string {
			return ""
		},
		"addSuffix": func(args ...String) string {
			return ""
		},
		"removeSuffix": func(args ...String) string {
			return ""
		},
	},

	{
		"archconv": func() archconvNS { return archconvNS{} },
		"toBytes":  func(s any) []byte { return nil },
	},
}

// placeholder functions to be overridden before template.Execute
var placeholderFuncMaps = [...]map[string]any{
	{
		// var as template variable
		"var": func() map[string]any { return nil },

		// include is like helm include
		"include": func(name string, data any) (string, error) {
			return "", fmt.Errorf("no implementation")
		},
	},
}
