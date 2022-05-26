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

	"github.com/stretchr/testify/assert"

	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/templateutils"
)

func TestGenerateTemplateFuncDocs(t *testing.T) {
	t.Parallel()

	ctx := dukkha_test.NewTestContext(context.TODO())
	tpl, err := templateutils.CreateTemplate(ctx).ParseFiles("template_funcs.tpl")

	if !assert.NoError(t, err) {
		return
	}

	tfs := collectTemplateFuncs()
	var buf bytes.Buffer
	if !assert.NoError(t, tpl.ExecuteTemplate(&buf, "template_funcs.tpl", tfs)) {
		return
	}

	assert.NoError(t, os.WriteFile("./generated/template_funcs.md", buf.Bytes(), 0644))
}

type templateFunc struct {
	Namespace string
	Name      string
	FuncType  string
}

func collectTemplateFuncs() []templateFunc {
	ctx := dukkha_test.NewTestContext(context.TODO())

	tpl := templateutils.CreateTemplate(ctx)

	replacer := strings.NewReplacer(
		"templateutils.", "",
		"interface {}", "any",
	)

	var ret []templateFunc
	for k, v := range tpl.GetExecFuncs() {
		vt := v.Type()

		if vt.NumIn() != 0 {
			// is a func
			ret = append(ret, templateFunc{
				Name:     k,
				FuncType: replacer.Replace(vt.String()),
			})
			continue
		}

		// namespaced func
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

			ret = append(ret, templateFunc{
				Namespace: k,
				Name:      m.Name,
				FuncType:  replacer.Replace(reflect.FuncOf(fin, fout, ft.IsVariadic()).String()),
			})
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Namespace+"."+ret[i].Name < ret[j].Namespace+"."+ret[j].Name
	})

	return ret
}
