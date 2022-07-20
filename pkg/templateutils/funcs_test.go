package templateutils

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"

	_ "embed" // for go:embed
)

var (
	//go:embed gen/funcs.tpl
	funcs_template string // nolint:revive

	//go:embed gen/docs.tpl
	funcs_docs string // nolint:revive

	//go:embed gen/tengo_symbols.tpl
	funcs_tengo_symbols string // nolint:revive
)

func TestGenerateFuncs(t *testing.T) {
	type Values struct {
		StaticFuncs    []TemplateFuncInfo
		LastStaticFunc TemplateFuncInfo

		ContextualFuncs    []TemplateFuncInfo
		LastContextualFunc TemplateFuncInfo

		PlaceholderFuncs    []TemplateFuncInfo
		LastPlaceholderFunc TemplateFuncInfo
	}

	var val Values

	val.StaticFuncs = collectTemplateFuncs(staticFuncMaps)
	val.LastStaticFunc = val.StaticFuncs[len(val.StaticFuncs)-1]

	val.ContextualFuncs = collectTemplateFuncs(contextualFuncs)
	val.LastContextualFunc = val.ContextualFuncs[len(val.ContextualFuncs)-1]

	val.PlaceholderFuncs = collectTemplateFuncs(placeholderFuncMaps)
	val.LastPlaceholderFunc = val.PlaceholderFuncs[len(val.PlaceholderFuncs)-1]

	t.Run("funcs", func(t *testing.T) {
		tpl, err := template.New("").Parse(funcs_template)
		assert.NoError(t, err)

		var out bytes.Buffer
		err = tpl.Execute(&out, val)
		if !assert.NoError(t, err) {
			return
		}

		data, err := format.Source(out.Bytes())
		if !assert.NoError(t, err) {
			return
		}

		err = os.WriteFile("funcs.go", data, 0644)
		assert.NoError(t, err)
	})

	t.Run("tengo_symbols", func(t *testing.T) {
		tpl, err := template.New("").Parse(funcs_tengo_symbols)
		assert.NoError(t, err)

		var out bytes.Buffer
		err = tpl.Execute(&out, val)
		if !assert.NoError(t, err) {
			return
		}

		data, err := format.Source(out.Bytes())
		if !assert.NoError(t, err) {
			return
		}

		err = os.WriteFile("../renderer/tengo/symbols.go", data, 0644)
		assert.NoError(t, err)
	})

	t.Run("docs", func(t *testing.T) {
		nStatic := len(val.StaticFuncs)
		nContextual := len(val.ContextualFuncs)
		nPlaceholder := len(val.PlaceholderFuncs)
		allFuncs := make([]TemplateFuncInfo, nStatic+nContextual+nPlaceholder)
		copy(allFuncs, val.StaticFuncs)
		copy(allFuncs[nStatic:], val.ContextualFuncs)
		copy(allFuncs[nStatic+nContextual:], val.PlaceholderFuncs)

		tpl, err := template.New("").Funcs(template.FuncMap{
			"strings": func() any { return stringsNS{} },
		}).Parse(funcs_docs)
		if !assert.NoError(t, err) {
			return
		}

		var out bytes.Buffer

		err = tpl.Execute(&out, allFuncs)
		if !assert.NoError(t, err) {
			return
		}

		err = os.WriteFile("../../docs/generated/template_funcs.md", out.Bytes(), 0644)
		assert.NoError(t, err)
	})
}

type TemplateFuncInfo struct {
	// UserCallHandle is the string with which user can call this func in template
	UserCallHandle string
	// CodeCallHandle is the code calling path
	CodeCallHandle string

	Ident string

	FuncType string
}

func collectTemplateFuncs(fms map[string]any) []TemplateFuncInfo {
	visited := make(map[string]struct{})

	var ret []TemplateFuncInfo
	for k, fn := range fms {
		tfs := parseTemplateFunc(k, fn)

		for _, tf := range tfs {
			if _, ok := visited[tf.Ident]; ok {
				continue
			}

			visited[tf.Ident] = struct{}{}
			ret = append(ret, tf)
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Ident < ret[j].Ident
	})

	return ret
}

var replacer = strings.NewReplacer(
	"templateutils.", "",
	"interface {}", "any",
	"reflect.Value", "any",
	"utils.LazyValue", "string",
)

func parseTemplateFunc(k string, fn any) []TemplateFuncInfo {
	ref, ok := fn.(FuncRef)
	if ok {
		// describes a func derived from namespace type

		m, ok := reflect.TypeOf(ref.nsType).MethodByName(ref.funcName)
		if !ok {
			panic(fmt.Errorf("method %q not found for %q", ref.funcName, k))
		}

		funcType, isExported := funcTypeOfMethod(m)
		if !isExported {
			panic(fmt.Errorf("invalid unexported method %q for %q", ref.funcName, k))
		}

		return []TemplateFuncInfo{{
			UserCallHandle: k,
			CodeCallHandle: "ns_" + ref.nsName + "." + ref.funcName,
			Ident:          k,
			FuncType:       replacer.Replace(funcType.String()),
		}}
	}

	// is a func, must only return a namespace struct

	typ := reflect.TypeOf(fn)

	if typ.NumIn() != 0 {
		panic(fmt.Errorf("invalid func %q with args: %q", k, typ.String()))
	}

	if typ.NumOut() != 1 {
		panic(fmt.Errorf("invalid func %q with multiple or none returns, expecting 1: %q", k, typ.String()))
	}

	ret0 := typ.Out(0)
	if ret0.Kind() != reflect.Struct {
		panic(fmt.Errorf("invalid return type of %q: %q", k, ret0.String()))
	}

	// namespace struct

	nMethods := ret0.NumMethod()
	ret := make([]TemplateFuncInfo, 0, nMethods+1)
	ret = append(ret, TemplateFuncInfo{
		UserCallHandle: k,
		CodeCallHandle: "get_ns_" + k,
		Ident:          k,
		FuncType:       "func() " + ret0.Name(),
	})

	for i := 0; i < nMethods; i++ {
		m := ret0.Method(i)
		funcType, isExported := funcTypeOfMethod(m)
		if !isExported {
			continue
		}

		ret = append(ret, TemplateFuncInfo{
			UserCallHandle: k + "." + m.Name,
			CodeCallHandle: "ns_" + k + "." + m.Name,
			Ident:          k + "_" + m.Name,
			FuncType:       replacer.Replace(funcType.String()),
		})
	}

	return ret
}

func funcTypeOfMethod(m reflect.Method) (ret reflect.Type, isExported bool) {
	var (
		fin  []reflect.Type
		fout []reflect.Type
	)

	funcType := m.Func.Type()

	// skip first (receiver)
	for i := 1; i < funcType.NumIn(); i++ {
		fin = append(fin, funcType.In(i))
	}

	for i := 0; i < funcType.NumOut(); i++ {
		fout = append(fout, funcType.Out(i))
	}

	ret = reflect.FuncOf(fin, fout, funcType.IsVariadic())
	isExported = len(m.PkgPath) == 0

	return
}

type FuncRef struct {
	nsName   string
	nsType   any
	funcName string
}

var staticFuncMaps = map[string]any{
	"archconv": func() archconvNS { return archconvNS{} },
	"path":     func() pathNS { return pathNS{} },
	"uuid":     func() uuidNS { return uuidNS{} },
	"re":       func() regexpNS { return regexpNS{} },

	// Math

	"math": func() mathNS { return mathNS{} },

	"seq": FuncRef{"math", mathNS{}, "Seq"},

	"min": FuncRef{"math", mathNS{}, "Min"},
	"max": FuncRef{"math", mathNS{}, "Max"},

	"mod": FuncRef{"math", mathNS{}, "Mod"},
	"add": FuncRef{"math", mathNS{}, "Add"},
	"sub": FuncRef{"math", mathNS{}, "Sub"},
	"mul": FuncRef{"math", mathNS{}, "Mul"},
	"div": FuncRef{"math", mathNS{}, "Div"},

	"add1":   FuncRef{"math", mathNS{}, "Add1"},
	"sub1":   FuncRef{"math", mathNS{}, "Sub1"},
	"double": FuncRef{"math", mathNS{}, "Double"},
	"half":   FuncRef{"math", mathNS{}, "Half"},

	// Collections

	"coll": func() collNS { return collNS{} },

	"list":       FuncRef{"coll", collNS{}, "List"},
	"stringList": FuncRef{"coll", collNS{}, "Strings"},
	"slice":      FuncRef{"coll", collNS{}, "Slice"},
	"index":      FuncRef{"coll", collNS{}, "Index"},
	"dict":       FuncRef{"coll", collNS{}, "MapStringAny"},
	"append":     FuncRef{"coll", collNS{}, "Append"},
	"prepend":    FuncRef{"coll", collNS{}, "Prepend"},
	"sort":       FuncRef{"coll", collNS{}, "Sort"},
	"has":        FuncRef{"coll", collNS{}, "HasAll"},
	"hasAny":     FuncRef{"coll", collNS{}, "HasAny"},
	"pick":       FuncRef{"coll", collNS{}, "Pick"},
	"omit":       FuncRef{"coll", collNS{}, "Omit"},
	"dup":        FuncRef{"coll", collNS{}, "Dup"},
	"uniq":       FuncRef{"coll", collNS{}, "Unique"},

	// Type conversion

	"type": func() typeNS { return typeNS{} },

	"close":    FuncRef{"type", typeNS{}, "Close"},
	"toString": FuncRef{"type", typeNS{}, "ToString"},
	"default":  FuncRef{"type", typeNS{}, "Default"},
	"all":      FuncRef{"type", typeNS{}, "AllTrue"},
	"any":      FuncRef{"type", typeNS{}, "AnyTrue"},

	// Network

	"dns":      func() dnsNS { return dnsNS{} },
	"sockaddr": func() sockaddrNS { return sockaddrNS{} },

	// Hashing and hmac

	"hash": func() hashNS { return hashNS{} },

	"md5":    FuncRef{"hash", hashNS{}, "MD5"},
	"sha1":   FuncRef{"hash", hashNS{}, "SHA1"},
	"sha256": FuncRef{"hash", hashNS{}, "SHA256"},
	"sha512": FuncRef{"hash", hashNS{}, "SHA512"},

	// Credentials

	"cred": func() credentialNS { return credentialNS{} },

	"totp": FuncRef{"cred", credentialNS{}, "Totp"},

	// Time

	"time": func() timeNS { return timeNS{} },

	"now": FuncRef{"time", timeNS{}, "Now"},

	// Encoding

	"enc": func() encNS { return encNS{} },

	"base64": FuncRef{"enc", encNS{}, "Base64"},
	"hex":    FuncRef{"enc", encNS{}, "Hex"},
	"toJson": FuncRef{"enc", encNS{}, "JSON"},
	"toYaml": FuncRef{"enc", encNS{}, "YAML"},

	// Strings

	"strings": func() stringsNS { return stringsNS{} },

	"replaceAll": FuncRef{"strings", stringsNS{}, "ReplaceAll"},
	"title":      FuncRef{"strings", stringsNS{}, "Title"},
	"upper":      FuncRef{"strings", stringsNS{}, "Upper"},
	"lower":      FuncRef{"strings", stringsNS{}, "Lower"},
	"indent":     FuncRef{"strings", stringsNS{}, "Indent"},
	"nindent":    FuncRef{"strings", stringsNS{}, "NIndent"},
	"quote":      FuncRef{"strings", stringsNS{}, "DoubleQuote"},
	"squote":     FuncRef{"strings", stringsNS{}, "SingleQuote"},
	"contains":   FuncRef{"strings", stringsNS{}, "Contains"},
	"split":      FuncRef{"strings", stringsNS{}, "Split"},
	"splitN":     FuncRef{"strings", stringsNS{}, "SplitN"},

	"trim":       FuncRef{"strings", stringsNS{}, "Trim"},
	"trimSpace":  FuncRef{"strings", stringsNS{}, "TrimSpace"},
	"trimPrefix": FuncRef{"strings", stringsNS{}, "TrimPrefix"},
	"trimSuffix": FuncRef{"strings", stringsNS{}, "TrimSuffix"},

	"hasPrefix": FuncRef{"strings", stringsNS{}, "HasPrefix"},
	"hasSuffix": FuncRef{"strings", stringsNS{}, "HasSuffix"},

	"addPrefix":    FuncRef{"strings", stringsNS{}, "AddPrefix"},
	"addSuffix":    FuncRef{"strings", stringsNS{}, "AddSuffix"},
	"removePrefix": FuncRef{"strings", stringsNS{}, "RemovePrefix"},
	"removeSuffix": FuncRef{"strings", stringsNS{}, "RemoveSuffix"},

	// golang built-in funcs (replaced `slice` and `index`)

	"call":     FuncRef{"golang", golangNS{}, "Call"},
	"html":     FuncRef{"golang", golangNS{}, "HTMLEscaper"},
	"js":       FuncRef{"golang", golangNS{}, "JSEscaper"},
	"len":      FuncRef{"golang", golangNS{}, "Length"},
	"and":      FuncRef{"golang", golangNS{}, "And"},
	"not":      FuncRef{"golang", golangNS{}, "Not"},
	"or":       FuncRef{"golang", golangNS{}, "Or"},
	"print":    FuncRef{"golang", golangNS{}, "Sprint"},
	"printf":   FuncRef{"golang", golangNS{}, "Sprintf"},
	"println":  FuncRef{"golang", golangNS{}, "Sprintln"},
	"urlquery": FuncRef{"golang", golangNS{}, "URLQueryEscaper"},
	"eq":       FuncRef{"golang", golangNS{}, "Eq"}, // ==
	"ge":       FuncRef{"golang", golangNS{}, "Ge"}, // >=
	"gt":       FuncRef{"golang", golangNS{}, "Gt"}, // >
	"le":       FuncRef{"golang", golangNS{}, "Le"}, // <=
	"lt":       FuncRef{"golang", golangNS{}, "Lt"}, // <
	"ne":       FuncRef{"golang", golangNS{}, "Ne"}, // !=
}

var contextualFuncs = map[string]any{
	// OS

	"os": func() osNS { return osNS{} },

	// Tagging

	"tag": func() tagNS { return tagNS{} },

	// Filesystem

	"fs": func() fsNS { return fsNS{} },

	"touch": FuncRef{"fs", fsNS{}, "Touch"},
	"write": FuncRef{"fs", fsNS{}, "WriteFile"},
	"mkdir": FuncRef{"fs", fsNS{}, "Mkdir"},
	"find":  FuncRef{"fs", fsNS{}, "Find"},

	// dukkha runtime values

	"dukkha": func() dukkhaNS { return dukkhaNS{} },

	"jq":    FuncRef{"dukkha", dukkhaNS{}, "JQ"},
	"yq":    FuncRef{"dukkha", dukkhaNS{}, "YQ"},
	"jqObj": FuncRef{"dukkha", dukkhaNS{}, "JQObj"},
	"yqObj": FuncRef{"dukkha", dukkhaNS{}, "YQObj"},

	"fromJson": FuncRef{"dukkha", dukkhaNS{}, "FromJson"},
	"fromYaml": FuncRef{"dukkha", dukkhaNS{}, "FromYaml"},

	// eval shell and template
	"eval": func() evalNS { return evalNS{} },

	// state for task execution
	"state": func() stateNS { return stateNS{} },

	"git":    FuncRef{"misc", miscNS{}, "Git"},  // git.{tag, branch ...}
	"host":   FuncRef{"misc", miscNS{}, "Host"}, // host.{arch, arch_simple, kernel ...}
	"env":    FuncRef{"misc", miscNS{}, "Env"},
	"values": FuncRef{"misc", miscNS{}, "Values"},
	"matrix": FuncRef{"misc", miscNS{}, "Matrix"},
	"VALUE":  FuncRef{"misc", miscNS{}, "VALUE"},
}

// placeholder functions to be overridden before template.Execute
var placeholderFuncMaps = map[string]any{
	"var":     FuncRef{"pending", pendingNS{}, "Var"},
	"include": FuncRef{"pending", pendingNS{}, "Include"},
}
