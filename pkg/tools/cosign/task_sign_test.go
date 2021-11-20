package cosign

import (
	"context"
	"reflect"
	"testing"

	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/template"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"
)

func TestTaskSign(t *testing.T) {
	type Expected struct {
		rs.BaseField

		Signature string `yaml:"signature"`
	}

	testTask(t, "./fixtures/sign", &Tool{},
		func() dukkha.Task { return newTaskSign("test") },
		func() rs.Field { return &Expected{} },
		func(e, a rs.Field) {
			exp, actual := e.(*Expected), a.(*Expected)
			_, _ = exp, actual
			// assert.EqualValues(t, exp.Signature, actual.Signature)
		},
	)
}

func testTask(
	t *testing.T,
	dir string,
	tool dukkha.Tool,
	newTask func() dukkha.Task,
	newExpected func() rs.Field,
	check func(expected, actual rs.Field),
) {
	type TestCase struct {
		rs.BaseField

		// Tool dukkha.Tool `yaml:"tool"`
		Task dukkha.Task `yaml:"task"`
	}

	type CheckSpec struct {
		rs.BaseField

		ExpectErr bool     `yaml:"expect_err"`
		Actual    rs.Field `yaml:"actual"`
		Expected  rs.Field `yaml:"expected"`
	}

	testhelper.TestFixtures(t, dir,
		func() interface{} {
			return rs.Init(&TestCase{}, &rs.Options{
				InterfaceTypeHandler: rs.InterfaceTypeHandleFunc(
					func(typ reflect.Type, yamlKey string) (interface{}, error) {
						return rs.Init(newTask(), nil), nil
					},
				),
			})
		},
		func() interface{} {
			return rs.Init(&CheckSpec{}, &rs.Options{
				InterfaceTypeHandler: rs.InterfaceTypeHandleFunc(
					func(typ reflect.Type, yamlKey string) (interface{}, error) {
						return rs.Init(newExpected(), nil), nil
					},
				),
			})
		},
		func(t *testing.T, in, exp interface{}) {
			defer t.Cleanup(func() {

			})
			spec := in.(*TestCase)
			e := exp.(*CheckSpec)

			ctx := dukkha_test.NewTestContext(context.TODO())
			ctx.SetCacheDir(t.TempDir())
			ctx.AddRenderer("file", file.NewDefault("file"))
			ctx.AddRenderer("env", env.NewDefault("env"))
			ctx.AddRenderer("template", template.NewDefault("template"))
			ctx.AddRenderer("shell", shell.NewDefault("shell"))

			if !assert.NoError(t, spec.ResolveFields(ctx, -1)) {
				return
			}

			rs.Init(tool, nil)

			tool.Init("", ctx.CacheDir())
			ctx.AddTool(tool.Key(), tool)

			tool.AddTasks([]dukkha.Task{spec.Task})

			err := tools.RunTask(&tools.TaskExecRequest{
				Context:     ctx,
				Tool:        tool,
				Task:        spec.Task,
				IgnoreError: false,
			})

			if !assert.NoError(t, e.ResolveFields(ctx, -1)) {
				return
			}

			if e.ExpectErr {
				assert.Error(t, err)
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			check(e.Expected, e.Actual)
		},
	)
}
