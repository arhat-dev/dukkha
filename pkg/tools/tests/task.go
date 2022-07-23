package tests

import (
	"context"
	"reflect"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/af"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/tmpl"
	"arhat.dev/dukkha/pkg/tools"
)

type ExecSpecGenerationTestCase struct {
	Name     string
	Prepare  func() error
	Finalize func()

	Options   dukkha.TaskMatrixExecOptions
	Task      dukkha.Task
	Expected  []dukkha.TaskExecSpec
	ExpectErr bool
}

func RunTaskExecSpecGenerationTests(
	t *testing.T,
	taskCtx dukkha.TaskExecContext,
	tests []ExecSpecGenerationTestCase,
) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runTaskTest(taskCtx, &test, t)
		})
	}
}

type baseTaskInitializer interface {
	InitBaseTask(
		k dukkha.ToolKind,
		n dukkha.ToolName,
		impl dukkha.Task,
	)
}

func runTaskTest(taskCtx dukkha.TaskExecContext, test *ExecSpecGenerationTestCase, t *testing.T) {
	if test.Finalize != nil {
		defer test.Finalize()
	}

	if test.Prepare != nil {
		if !assert.NoError(t, test.Prepare(), "preparing test environment") {
			return
		}
	}

	rs.InitRecursively(reflect.ValueOf(test.Task), nil)

	// nolint:gocritic
	switch t := test.Task.(type) {
	case baseTaskInitializer:
		t.InitBaseTask("test-tool", "test-tool-name", test.Task)
	}

	assert.NoError(t, test.Task.Init(taskCtx.(dukkha.ConfigResolvingContext).TaskCacheFS(test.Task)))

	if test.ExpectErr {
		_, err := test.Task.GetExecSpecs(taskCtx, test.Options)
		assert.Error(t, err)
		return
	}

	specs, err := test.Task.GetExecSpecs(taskCtx, test.Options)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, test.Expected, specs)
}

func TestTask(
	t *testing.T,
	dir string,
	tool dukkha.Tool,
	newTask func() dukkha.Task,
	newExpected func() rs.Field,
	check func(t *testing.T, expected, actual rs.Field),
) {
	type TestCase struct {
		rs.BaseField

		Env dukkha.Env `yaml:"env"`

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
		func() any {
			return rs.InitAny(&TestCase{}, &rs.Options{
				InterfaceTypeHandler: rs.InterfaceTypeHandleFunc(
					func(typ reflect.Type, yamlKey string) (any, error) {
						return rs.InitAny(newTask(), nil), nil
					},
				),
			})
		},
		func() any {
			return rs.InitAny(&CheckSpec{}, &rs.Options{
				InterfaceTypeHandler: rs.InterfaceTypeHandleFunc(
					func(typ reflect.Type, yamlKey string) (any, error) {
						return rs.InitAny(newExpected(), nil), nil
					},
				),
			})
		},
		func(t *testing.T, in, exp any) {
			defer t.Cleanup(func() {

			})
			spec := in.(*TestCase)
			e := exp.(*CheckSpec)

			ctx := dt.NewTestContext(context.TODO(), t.TempDir())
			ctx.AddRenderer("file", file.NewDefault("file"))
			ctx.AddRenderer("env", env.NewDefault("env"))
			ctx.AddRenderer("tmpl", tmpl.NewDefault("tmpl"))
			ctx.AddRenderer("shell", shell.NewDefault("shell"))

			afr := af.NewDefault("af")
			assert.NoError(t, afr.Init(ctx.RendererCacheFS("af")))
			ctx.AddRenderer("af", afr)

			if !assert.NoError(t, dukkha.ResolveEnv(ctx, spec, "Env", "env")) {
				return
			}

			if !assert.NoError(t, spec.ResolveFields(ctx, -1)) {
				return
			}

			rs.InitAny(tool, nil)

			assert.NoError(t, tool.Init(ctx.ToolCacheFS(tool)))
			ctx.AddTool(tool.Key(), tool)

			assert.NoError(t, tool.AddTasks([]dukkha.Task{spec.Task}))

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

			check(t, e.Expected, e.Actual)
		},
	)
}
