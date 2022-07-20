package tools

import (
	"context"
	"fmt"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	di "arhat.dev/dukkha/internal"
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
		func() any { return rs.InitAny(&TestResolvable{}, nil) },
		func() any { return rs.InitAny(&CheckSpec{}, nil) },
		func(t *testing.T, spec, exp any) {
			in := spec.(*TestResolvable)
			cs := exp.(*CheckSpec)

			mCtx := dt.NewTestContext(context.TODO())
			mCtx.(di.CacheDirSetter).SetCacheDir(t.TempDir())
			mCtx.AddRenderer("tmpl", tmpl.NewDefault(""))
			mCtx.AddRenderer("file", file.NewDefault(""))
			mCtx.AddRenderer("echo", echo.NewDefault(""))

			jobs, err := ResolveActions(
				mCtx, in, "Actions", "actions",
			)

			assert.NoError(t, err)

			assertActionsVisibleFields := func(t *testing.T, ex, ac []*Action) {
				for i := range ex {
					expected, actual := ex[i], ac[i]

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

			i := 0
			for jobs != nil {
				t.Run(fmt.Sprint(i), func(t *testing.T) {
					assertActionsVisibleFields(t, cs.Steps[i].Actions, in.Actions)
					i++

					// call AlterExecFunc in second spec is like calling
					// next() to go to next step
					var ret any
					ret, err = jobs[1].AlterExecFunc(nil, nil, nil, nil)

					if cs.Steps[i].Error {
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

			assert.Equal(t, len(cs.Steps), i+1, "missing steps")
		},
	)
}
