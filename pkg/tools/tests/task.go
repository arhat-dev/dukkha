package tests

import (
	"context"
	"reflect"
	"testing"

	"arhat.dev/pkg/fshelper"
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
	for i := range tests {
		test := &tests[i]
		t.Run(test.Name, func(t *testing.T) {
			runTaskTest(taskCtx, test, t)
		})
	}
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

	tmp := t.TempDir()
	assert.NoError(t, test.Task.Init(fshelper.NewOSFS(false, func(op fshelper.Op, name string) (string, error) {
		return tmp, nil
	})))

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

type TaskTestCase[Tsk dukkha.Task] struct {
	rs.BaseField

	Env dukkha.NameValueList `yaml:"env"`

	// Tool dukkha.Tool `yaml:"tool"`
	Task Tsk `yaml:"task"`
}

type TaskCheckSpec[Exp rs.Field] struct {
	rs.BaseField

	ExpectErr bool `yaml:"expect_err"`
	Actual    Exp  `yaml:"actual"`
	Expected  Exp  `yaml:"expected"`
}

func TestTask[Tool dukkha.Tool, Task dukkha.Task, CheckSpec rs.Field](
	t *testing.T,
	dir string,
	tool Tool,
	newTask func() Task,
	newCheckSpec func() CheckSpec,
	check func(t *testing.T, exp, actual CheckSpec),
) {
	testhelper.TestFixtures(t, dir,
		func() *TaskTestCase[Task] {
			ret := &TaskTestCase[Task]{}
			ret.Task = newTask()

			rs.InitRecursively(reflect.ValueOf(ret), nil)
			return ret
		},
		func() *TaskCheckSpec[CheckSpec] {
			ret := &TaskCheckSpec[CheckSpec]{}
			ret.Expected = newCheckSpec()
			ret.Actual = newCheckSpec()

			rs.InitRecursively(reflect.ValueOf(ret), nil)
			return ret
		},
		func(t *testing.T, testCase *TaskTestCase[Task], exp *TaskCheckSpec[CheckSpec]) {
			ctx := dt.NewTestContext(context.TODO(), t.TempDir())
			ctx.AddRenderer("file", file.NewDefault("file"))
			ctx.AddRenderer("env", env.NewDefault("env"))
			ctx.AddRenderer("tmpl", tmpl.NewDefault("tmpl"))
			ctx.AddRenderer("shell", shell.NewDefault("shell"))

			afr := af.NewDefault("af")
			assert.NoError(t, afr.Init(ctx.RendererCacheFS("af")))
			ctx.AddRenderer("af", afr)

			if !assert.NoError(t, dukkha.ResolveAndAddEnv(ctx, testCase, "Env", "env")) {
				return
			}

			if !assert.NoError(t, testCase.ResolveFields(ctx, -1)) {
				return
			}

			rs.InitRecursively(reflect.ValueOf(tool), nil)

			assert.NoError(t, tool.Init(ctx.ToolCacheFS(tool)))
			ctx.AddTool(tool.Key(), tool)

			testCase.Task.Init(ctx.ToolCacheFS(tool))

			assert.NoError(t, tool.AddTasks([]dukkha.Task{testCase.Task}))

			err := tools.RunTask(&tools.TaskExecRequest{
				Context:     ctx,
				Tool:        tool,
				Task:        testCase.Task,
				IgnoreError: false,
			})

			if !assert.NoError(t, exp.ResolveFields(ctx, -1)) {
				return
			}

			if exp.ExpectErr {
				assert.Error(t, err)
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			check(t, exp.Expected, exp.Actual)
		},
	)
}
