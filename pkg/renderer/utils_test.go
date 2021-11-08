package renderer

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"arhat.dev/pkg/iohelper"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
)

func TestHandleRenderingRequestWithRemoteFetch(t *testing.T) {
	obj := cache.IdentifiableString("test")
	const cachedData = "test-data"
	fetchRemoteAlwaysOk := cache.RemoteCacheRefreshFunc(func(obj cache.IdentifiableObject) (io.ReadCloser, error) {
		return iohelper.CustomReadCloser(strings.NewReader("test-data"), nil), nil
	})

	t.Run("allow-expired", func(t *testing.T) {
		defer t.Cleanup(func() {})

		c := cache.NewTwoTierCache(t.TempDir(), -1, -1, -1)
		data, err := HandleRenderingRequestWithRemoteFetch(
			c, obj, fetchRemoteAlwaysOk, []dukkha.RendererAttribute{"allow-expired"},
		)
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))
	})

	t.Run("cached-file", func(t *testing.T) {
		defer t.Cleanup(func() {})

		dir := t.TempDir()
		c := cache.NewTwoTierCache(dir, -1, -1, -1)
		data, err := HandleRenderingRequestWithRemoteFetch(
			c, obj, fetchRemoteAlwaysOk, []dukkha.RendererAttribute{"cached-file"},
		)
		assert.NoError(t, err)
		assert.EqualValues(t, dir, filepath.Dir(string(data)))

		data, err = os.ReadFile(string(data))
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))
	})
}
