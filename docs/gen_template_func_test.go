package docs_test

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"sync"
	"text/template"
	"text/template/parse"
	"unsafe"

	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/templateutils"
)

//go:linkname _Template text/template.Template
type _Template struct {
	name string
	*parse.Tree
	*_common
	leftDelim  string
	rightDelim string
}

//go:linkname _common text/template.common
type _common struct {
	tmpl   map[string]*_Template // Map from name to defined templates.
	muTmpl sync.RWMutex          // protects tmpl
	option _option
	// We use two maps, one for parsing and one for execution.
	// This separation makes the API cleaner since it doesn't
	// expose reflection to the client.
	muFuncs    sync.RWMutex // protects parseFuncs and execFuncs
	parseFuncs template.FuncMap
	execFuncs  map[string]reflect.Value
}

//go:linkname _missingKeyAction text/template.missingKeyAction
type _missingKeyAction int

//go:linkname _option text/template.option
type _option struct {
	missingKey _missingKeyAction
}

type templateFunc struct {
	Name string
	Func string
}

func collectTemplateFuncs() []*templateFunc {
	stdTpl := templateutils.CreateTemplate(dukkha_test.NewTestContext(context.TODO()))
	tpl := (*_Template)(unsafe.Pointer(stdTpl))

	var ret []*templateFunc
	for k, v := range tpl.execFuncs {
		vt := v.Type()

		if vt.NumIn() != 0 {
			// is a func
			ret = append(ret, &templateFunc{
				Name: k,
				Func: vt.String(),
			})
			continue
		}

		// using namespaced func
		ns := v.Call(nil)[0].Type()
		for i := 0; i < ns.NumMethod(); i++ {
			m := ns.Method(i)
			if len(m.PkgPath) != 0 {
				// unexported, ignore
				continue
			}

			ret = append(ret, &templateFunc{
				Name: k + "." + m.Name,
				Func: m.Func.Type().String(),
			})
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		ni, nj := ret[i].Name, ret[j].Name
		switch {
		case strings.Contains(ni, "."):
			if !strings.Contains(nj, ".") {
				return false
			}

			return ni < nj
		case strings.Contains(nj, "."):
			return true
		default:
			return ni < nj
		}
	})

	return ret
}
