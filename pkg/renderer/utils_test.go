package renderer

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/iohelper"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
)

func TestHandleRenderingRequestWithRemoteFetch(t *testing.T) {
	t.Parallel()

	obj := cache.IdentifiableString("test")
	const cachedData = "test-data"
	fetchRemoteAlwaysOk := cache.RemoteCacheRefreshFunc(func(obj cache.IdentifiableObject) (io.ReadCloser, error) {
		return iohelper.CustomReadCloser(strings.NewReader("test-data"), nil), nil
	})

	t.Run("allow-expired", func(t *testing.T) {
		defer t.Cleanup(func() {})

		cacheDir := t.TempDir()
		c := cache.NewTwoTierCache(fshelper.NewOSFS(false, func(fshelper.Op) (string, error) {
			return cacheDir, nil
		}), -1, -1, -1)
		data, err := HandleRenderingRequestWithRemoteFetch(
			c, obj, fetchRemoteAlwaysOk, []dukkha.RendererAttribute{"allow-expired"},
		)
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))
	})

	t.Run("cached-file", func(t *testing.T) {
		defer t.Cleanup(func() {})

		cacheDir := t.TempDir()
		c := cache.NewTwoTierCache(fshelper.NewOSFS(false, func(fshelper.Op) (string, error) {
			return cacheDir, nil
		}), -1, -1, -1)
		data, err := HandleRenderingRequestWithRemoteFetch(
			c, obj, fetchRemoteAlwaysOk, []dukkha.RendererAttribute{"cached-file"},
		)
		assert.NoError(t, err)
		assert.EqualValues(t, cacheDir, filepath.Dir(string(data)))

		data, err = os.ReadFile(string(data))
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))
	})
}
