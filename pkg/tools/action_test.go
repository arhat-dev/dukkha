package tools

import (
	"context"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/env"

	_ "embed"
)

func TestActionFixtures(t *testing.T) {
	type CheckSpec struct {
		rs.BaseField

		Resolved Action `yaml:"resolved"`
		Result   struct {
			Failed bool `yaml:"failed"`
		} `yaml:"result"`
	}

	assertVisibleFields := func(t *testing.T, expected, actual *Action) bool {
		ok := assert.EqualValues(t, expected.Env, expected.Env) &&
			assert.EqualValues(t, expected.Name, expected.Name) &&
			assert.EqualValues(t, expected.Next, expected.Next) &&
			assert.EqualValues(t, expected.Cmd, actual.Cmd) &&
			assert.EqualValues(t, expected.Chdir, actual.Chdir) &&
			assert.EqualValues(t, expected.Idle, actual.Idle) &&
			assert.EqualValues(t, expected.ContinueOnError, actual.ContinueOnError) &&
			assert.EqualValues(t, expected.EmbeddedShell, actual.EmbeddedShell) &&
			assert.EqualValues(t, expected.ExternalShell, actual.ExternalShell) &&
			assert.EqualValues(t, expected.Task, actual.Task)

		return ok
	}

	testhelper.TestFixtures(t, "./_fixtures/action",
		func() interface{} { return rs.Init(&Action{}, nil).(*Action) },
		func() interface{} { return rs.Init(&CheckSpec{}, nil).(*CheckSpec) },
		func(t *testing.T, in interface{}, exp interface{}) {
			actual := in.(*Action)
			expected := exp.(*CheckSpec)

			ctx := dukkha_test.NewTestContext(context.TODO())
			ctx.SetCacheDir(t.TempDir())
			ctx.AddRenderer("env", env.NewDefault(""))

			assert.NoError(t, actual.DoAfterFieldResolved(ctx, func() error { return nil }))

			if !assertVisibleFields(t, &expected.Resolved, actual) {
				return
			}

			runReq, err := actual.GenSpecs(ctx, 0)
			if !assert.NoError(t, err) {
				return
			}

			switch rt := runReq.(type) {
			case []dukkha.TaskExecSpec:
				err = doRun(ctx, nil, rt, nil)
			case *TaskExecRequest:
				err = RunTask(rt)
			case nil:
				err = nil
			}

			if expected.Result.Failed {
				assert.EqualValues(t, dukkha.TaskExecFailed, ctx.State())
				if actual.ContinueOnError {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			} else {
				assert.EqualValues(t, dukkha.TaskExecSucceeded, ctx.State())
				assert.NoError(t, err)
			}
		},
	)
}
