package docs_test

import (
	"bytes"
	"context"
	_ "embed"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"text/template"
	"text/template/parse"
	"unsafe"

	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/templateutils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTemplateFuncDocs(t *testing.T) {
	tpl, err := template.New("").ParseFiles("template_funcs.tpl")
	if !assert.NoError(t, err) {
		return
	}

	tfs := collectTemplateFuncs()
	buf := &bytes.Buffer{}
	if !assert.NoError(t, tpl.ExecuteTemplate(buf, "template_funcs.tpl", tfs)) {
		return
	}

	assert.NoError(t, os.WriteFile("./renderers/template_funcs.md", buf.Bytes(), 0644))
}

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

			ret = append(ret, &templateFunc{
				Name: k + "." + m.Name,
				Func: reflect.FuncOf(fin, fout, ft.IsVariadic()).String(),
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
