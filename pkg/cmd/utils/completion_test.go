package utils

import (
	"context"
	"testing"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// editorconfig-checker-disable
const testDukkhaConfig = `
tools:
  workflow:
  - name: local

workflow:run:
- name: test
  matrix:
    a: [a1,a2]
    b: [b,c]
`

// editorconfig-checker-enable

func newCompletionContext(t *testing.T) dukkha.Context {
	ctx := dukkha_test.NewTestContext(context.TODO())

	config := conf.NewConfig()
	err := yaml.Unmarshal([]byte(testDukkhaConfig), config)
	if !assert.NoError(t, err) {
		panic(err)
	}

	err = config.Resolve(ctx, true)
	if !assert.NoError(t, err) {
		panic(err)
	}

	return ctx
}
