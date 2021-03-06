package tools

import (
	"context"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
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
		func() *Action { return rs.Init(&Action{}, nil).(*Action) },
		func() *CheckSpec { return rs.Init(&CheckSpec{}, nil).(*CheckSpec) },
		func(t *testing.T, spec *Action, exp *CheckSpec) {
			ctx := dt.NewTestContext(context.TODO(), t.TempDir())
			ctx.AddRenderer("env", env.NewDefault(""))

			assert.NoError(t, spec.DoAfterFieldResolved(ctx, func(bool) error { return nil }))

			if !assertVisibleFields(t, &exp.Resolved, spec) {
				return
			}

			runReq, err := spec.GenSpecs(ctx, 0)
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

			if exp.Result.Failed {
				assert.EqualValues(t, dukkha.TaskExecFailed, ctx.State())
				if spec.ContinueOnError {
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
