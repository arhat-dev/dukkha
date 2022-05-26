package templateutils

import (
	"fmt"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/third_party/golang/text/template"
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
		nsTag    = createTagNS(rc)
		nsOS     = createOSNS(rc)
		nsDukkha = createDukkhaNS(rc)
		nsFS     = createFSNS(rc)
	)

	return template.New("tpl").
		Funcs(fm). // tool specific functions
		Funcs(map[string]any{
			"close": close,

			"archconv": func() archconvNS { return archconvNS{} },
			"path":     func() pathNS { return pathNS{} },
			"uuid":     func() uuidNS { return uuidNS{} },
			"re":       func() regexpNS { return regexpNS{} },

			// Math

			"math": func() mathNS { return mathNS{} },

			"seq": mathNS{}.Seq,

			"min": mathNS{}.Min,
			"max": mathNS{}.Max,

			"mod": mathNS{}.Mod,
			"add": mathNS{}.Add,
			"sub": mathNS{}.Sub,
			"mul": mathNS{}.Mul,
			"div": mathNS{}.Div,

			"add1":   mathNS{}.Add1,
			"sub1":   mathNS{}.Sub1,
			"double": mathNS{}.Double,
			"half":   mathNS{}.Half,

			// Collections

			"coll": func() collNS { return collNS{} },

			"last":       func(s Slice) any { return must(collNS{}.Index(-1, s)) },
			"first":      func(s Slice) any { return must(collNS{}.Index(0, s)) },
			"list":       collNS{}.List,
			"stringList": collNS{}.Strings,
			"dict":       collNS{}.MapStringAny,
			"append":     collNS{}.Append,
			"prepend":    collNS{}.Prepend,
			"sort":       collNS{}.Sort,
			"has":        collNS{}.HasAll,
			"hasAny":     collNS{}.HasAny,
			"pick":       collNS{}.Pick,
			"omit":       collNS{}.Omit,
			"dup":        collNS{}.Dup,
			"uniq":       collNS{}.Unique,

			// Type conversion

			"type": func() typeNS { return typeNS{} },

			"toString": typeNS{}.ToString,
			"default":  typeNS{}.Default,
			"all":      typeNS{}.AllTrue,
			"any":      typeNS{}.AnyTrue,

			// Network

			"dns":      func() dnsNS { return dnsNS{} },
			"sockaddr": func() sockaddrNS { return sockaddrNS{} },

			// Hashing and hmac

			"hash": func() hashNS { return hashNS{} },

			"md5":    hashNS{}.MD5,
			"sha1":   hashNS{}.SHA1,
			"sha256": hashNS{}.SHA256,
			"sha512": hashNS{}.SHA512,

			// Credentials

			"cred": func() credentialNS { return credentialNS{} },

			"totp": credentialNS{}.Totp,

			// Time

			"time": func() timeNS { return timeNS{} },

			"now": timeNS{}.Now,

			// Encoding

			"enc": func() encNS { return encNS{} },

			"base64": encNS{}.Base64,
			"hex":    encNS{}.Hex,
			"toJson": encNS{}.JSON,
			"toYaml": encNS{}.YAML,

			"fromJson": func(v String) (any, error) { return fromYaml(rc, v) },
			"fromYaml": func(v String) (any, error) { return fromYaml(rc, v) },

			// Strings

			"strings": func() stringsNS { return stringsNS{} },

			"replaceAll": stringsNS{}.ReplaceAll,
			"title":      stringsNS{}.Title,
			"upper":      stringsNS{}.Upper,
			"lower":      stringsNS{}.Lower,
			"trimSpace":  stringsNS{}.TrimSpace,
			"indent":     stringsNS{}.Indent,
			"nindent":    stringsNS{}.NIndent,
			"quote":      stringsNS{}.DoubleQuote,
			"squote":     stringsNS{}.SingleQuote,

			"contains":  stringsNS{}.Contains,
			"hasPrefix": stringsNS{}.HasPrefix,
			"hasSuffix": stringsNS{}.HasSuffix,
			"split":     stringsNS{}.Split,
			"splitN":    stringsNS{}.SplitN,

			"trim":       stringsNS{}.Trim,
			"trimPrefix": stringsNS{}.TrimPrefix,
			"trimSuffix": stringsNS{}.TrimSuffix,

			"addPrefix":    stringsNS{}.AddPrefix,
			"addSuffix":    stringsNS{}.AddSuffix,
			"removePrefix": stringsNS{}.RemovePrefix,
			"removeSuffix": stringsNS{}.RemoveSuffix,

			// contextual template functions

			// OS

			"os": func() osNS { return nsOS },

			// Tagging

			"tag": func() tagNS { return nsTag },

			// Filesystem

			"fs": func() fsNS { return nsFS },

			"touch": func(file String) (struct{}, error) { return nsFS.WriteFile(file) },
			"write": nsFS.WriteFile,
			"mkdir": nsFS.Mkdir,
			"find":  nsFS.Find,

			// dukkha specific

			"dukkha": func() dukkhaNS { return nsDukkha },

			"jq":    nsDukkha.JQ,
			"yq":    nsDukkha.YQ,
			"jqObj": nsDukkha.JQObj,
			"yqObj": nsDukkha.YQObj,

			"git":  rc.GitValues,  // git.{tag, branch ...}
			"host": rc.HostValues, // host.{arch, arch_simple, kernel ...}
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

			// placeholder functions to be overridden before template.Execute

			"var": func() map[string]any { return nil },
			// include like helm include
			"include": func(name string, data any) (string, error) {
				return "", fmt.Errorf("no implementation")
			},
		})
}
