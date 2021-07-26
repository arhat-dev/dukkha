package buildah

import (
	"testing"

	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"github.com/stretchr/testify/assert"
)

func TestTaskPush_ManifestHandling(t *testing.T) {
	const (
		manifestName = "foo:latest"
	)

	task := &TaskPush{
		manifestCache: make(map[manifestCacheKey]manifestCacheValue),
	}

	opts := dukkha_test.CreateTaskMatrixExecOptions([]string{"buildah"})
	task.cacheManifestPushSpec(0, opts, manifestName)

	assert.Len(t, task.manifestCache, 1)
	for k, v := range task.manifestCache {
		assert.Equal(t, opts.ID(), k.execID)
		assert.Equal(t, 0, v.subIndex)
		assert.Equal(t, opts, v.opts)
		assert.Equal(t, manifestName, v.name)
	}

	task.cacheManifestPushSpec(9, opts, manifestName)
	assert.Len(t, task.manifestCache, 1)

	assert.EqualValues(t, []dukkha.TaskExecSpec{{
		Command: []string{
			"buildah", "manifest", "push", "--all",
			getLocalManifestName(manifestName),
			"docker://" + manifestName,
		},
	}}, task.createManifestPushSpecsFromCache(opts.ID()))
}
