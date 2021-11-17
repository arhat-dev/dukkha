package docs_test

import (
	"bytes"
	"context"
	_ "embed"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/templateutils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTemplateFuncDocs(t *testing.T) {
	ctx := dukkha_test.NewTestContext(context.TODO())
	ctx.SetCacheDir(t.TempDir())
	tpl, err := templateutils.CreateTemplate(ctx).ParseFiles("template_funcs.tpl")

	if !assert.NoError(t, err) {
		return
	}

	tfs := collectTemplateFuncs(t)
	buf := &bytes.Buffer{}
	if !assert.NoError(t, tpl.ExecuteTemplate(buf, "template_funcs.tpl", tfs)) {
		return
	}

	assert.NoError(t, os.WriteFile("./generated/template_funcs.md", buf.Bytes(), 0644))
}

type templateFunc struct {
	Name string
	Func string
}

func collectTemplateFuncs(t *testing.T) []*templateFunc {
	ctx := dukkha_test.NewTestContext(context.TODO())
	ctx.SetCacheDir(t.TempDir())

	tpl := templateutils.CreateTemplate(ctx)

	var ret []*templateFunc
	for k, v := range tpl.GetExecFuncs() {
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
