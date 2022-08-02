package tools

import (
	"context"
	"fmt"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/echo"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/tmpl"
)

var _ dukkha.Resolvable = (*TestResolvable)(nil)

type TestResolvable struct {
	rs.BaseField

	Actions Actions `yaml:"actions"`
}

func (r *TestResolvable) DoAfterFieldsResolved(
	mCtx dukkha.RenderingContext, depth int, addEnv bool, do func() error, tagNames ...string,
) error {
	return do()
}

func TestResolveActions_steps(t *testing.T) {
	t.Parallel()

	type CheckSpec struct {
		rs.BaseField

		Steps []struct {
			rs.BaseField

			TestResolvable `yaml:",inline"`
			Error          bool `yaml:"error"`
		} `yaml:"steps"`
	}

	testhelper.TestFixtures(t, "./_fixtures/actions/resolve",
		func() *TestResolvable { return rs.Init(&TestResolvable{}, nil).(*TestResolvable) },
		func() *CheckSpec { return rs.Init(&CheckSpec{}, nil).(*CheckSpec) },
		func(t *testing.T, spec *TestResolvable, exp *CheckSpec) {
			t.Parallel()

			mCtx := dt.NewTestContext(context.TODO(), t.TempDir())
			mCtx.AddRenderer("tmpl", tmpl.NewDefault(""))
			mCtx.AddRenderer("file", file.NewDefault(""))
			mCtx.AddRenderer("echo", echo.NewDefault(""))

			jobs, err := ResolveActions(mCtx, spec, &spec.Actions, "actions")
			if !assert.NoError(t, err) {
				return
			}

			if jobs == nil {
				assertActionsVisibleFields(t, exp.Steps[0].Actions, spec.Actions)
			}

			i := 0
			for jobs != nil {
				t.Run(fmt.Sprint(i), func(t *testing.T) {
					assertActionsVisibleFields(t, exp.Steps[i].Actions, spec.Actions)
					i++

					// calling AlterExecFunc in second spec is like calling
					// next() to go to next step
					var ret any
					ret, err = jobs[1].AlterExecFunc(nil, nil, nil, nil)

					if exp.Steps[i].Error {
						assert.Error(t, err)
						jobs = nil
						return
					}

					if !assert.NoError(t, err) {
						return
					}

					jobs = ret.([]dukkha.TaskExecSpec)
				})
			}

			assert.Equal(t, len(exp.Steps), i+1, "missing steps")
		},
	)
}

func assertActionsVisibleFields(t *testing.T, ex, ac []*Action) {
	for i := range ex {
		expected, actual := ex[i], ac[i]

		assert.EqualValues(t, expected.Run, actual.Run)
		assert.EqualValues(t, expected.Name, actual.Name)
		assert.EqualValues(t, expected.Env, actual.Env)
		assert.EqualValues(t, expected.Chdir, actual.Chdir)
		assert.EqualValues(t, expected.ContinueOnError, actual.ContinueOnError)
		assert.EqualValues(t, expected.Cmd, actual.Cmd)
		assert.EqualValues(t, expected.Idle, actual.Idle)
		assert.EqualValues(t, expected.Task, actual.Task)
		assert.EqualValues(t, expected.EmbeddedShell, actual.EmbeddedShell)
		assert.EqualValues(t, expected.ExternalShell, actual.ExternalShell)
	}
}
