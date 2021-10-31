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
	type testInputSpec struct {
		rs.BaseField

		Env  dukkha.Env `yaml:"env"`
		Spec Action     `yaml:"spec"`
	}

	testhelper.TestFixtures(t, "./_fixtures/action",
		func() interface{} { return rs.Init(&testInputSpec{}, nil).(*testInputSpec) },
		func() interface{} { return rs.Init(&Action{}, nil).(*Action) },
		func(t *testing.T, in interface{}, exp interface{}) {
			actual := in.(*testInputSpec)
			expected := exp.(*Action)

			ctx := dukkha_test.NewTestContext(context.TODO())
			ctx.AddRenderer("env", env.NewDefault(""))
			ctx.AddEnv(true, actual.Env...)

			assert.NoError(t, actual.Spec.ResolveFields(ctx, -1))

			t.Log(actual)

			assert.EqualValues(t, expected.Cmd, actual.Spec.Cmd)
			assert.EqualValues(t, expected.ContinueOnError, actual.Spec.ContinueOnError)
			assert.EqualValues(t, expected.EmbeddedShell, actual.Spec.EmbeddedShell)
			assert.EqualValues(t, expected.ExternalShell, actual.Spec.ExternalShell)
			assert.EqualValues(t, expected.Task, actual.Spec.Task)
		},
	)
}
