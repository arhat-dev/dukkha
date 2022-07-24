package dukkha_internal

import (
	"context"
	"testing"

	"arhat.dev/tlang"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
)

func TestInternalTypes(t *testing.T) {
	t.Parallel()

	fakeGlobalEnv := &dukkha.GlobalEnvSet{}
	for i := range fakeGlobalEnv {
		fakeGlobalEnv[i] = tlang.ImmediateString("")
	}

	ctx := dukkha.NewConfigResolvingContext(context.TODO(), nil, fakeGlobalEnv)

	_, ok := ctx.(DefaultGitBranchOverrider)
	assert.True(t, ok)

	_, ok = ctx.(WorkDirOverrider)
	assert.True(t, ok)

	_, ok = ctx.(VALUEGetter)
	assert.True(t, ok)

	_, ok = ctx.(VALUESetter)
	assert.True(t, ok)
}
