package dukkha_internal

import (
	"context"
	"testing"

	"arhat.dev/dukkha/pkg/dukkha"
	"github.com/stretchr/testify/assert"
)

func TestInternalTypes(t *testing.T) {

	ctx := dukkha.NewConfigResolvingContext(context.TODO(), nil, nil)

	_, ok := ctx.(DefaultGitBranchOverrider)
	assert.True(t, ok)

	_, ok = ctx.(WorkingDirOverrider)
	assert.True(t, ok)

	_, ok = ctx.(CacheDirSetter)
	assert.True(t, ok)

	_, ok = ctx.(VALUEGetter)
	assert.True(t, ok)

	_, ok = ctx.(VALUESetter)
	assert.True(t, ok)
}
