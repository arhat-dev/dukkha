package cache

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestTwoTierCache(t *testing.T) {
	t.Parallel()

	const (
		cachedData = "test-data"

		cacheID = "foo"
		// cacheFilenamePrefix is calculated from cacheID
		cacheFilenamePrefix = "acbd18db4cc2f85cedef654fccc4a4d8"
	)

	calledOk := 0
	fetchRemoteAlwaysOk := RemoteCacheRefreshFunc(func(_ IdentifiableObject) (io.ReadCloser, error) {
		calledOk++
		return ioutil.NopCloser(strings.NewReader(cachedData)), nil
	})

	calledFail := 0
	fetchRemoteError := fmt.Errorf("test error")
	fetchRemoteAlwaysFail := RemoteCacheRefreshFunc(func(_ IdentifiableObject) (io.ReadCloser, error) {
		calledFail++
		return nil, fetchRemoteError
	})

	obj := IdentifiableString("foo")

	t.Run("All Zero", func(t *testing.T) {
		defer t.Cleanup(func() {
			calledFail = 0
			calledOk = 0
		})

		cacheDir := t.TempDir()
		cache := NewTwoTierCache(fshelper.NewOSFS(false, func(fshelper.Op, string) (string, error) {
			return cacheDir, nil
		}), 0, 0, 0)

		data, expired, err := cache.Get(obj, 1111111110, true, fetchRemoteAlwaysOk)
		assert.EqualValues(t, 1, calledOk)
		assert.NoError(t, err)
		assert.False(t, expired)
		assert.EqualValues(t, cachedData, string(data))

		data, ok := cache.memcache.Get(cacheID)
		assert.False(t, ok)
		assert.Nil(t, data)
		assert.Zero(t, cache.memcache.Size())

		t.Run("Not Fetching Remote When Not Expired", func(t *testing.T) {
			// should not call fetch remote at the exact same time
			path, expired, err := cache.GetPath(obj, 1111111110, true, fetchRemoteAlwaysFail)
			assert.EqualValues(t, 1, calledOk)
			assert.EqualValues(t, 0, calledFail)
			assert.EqualValues(t, cacheFilenamePrefix+"-00000000001111111110", filepath.Base(path))
			assert.EqualValues(t, cacheDir, filepath.Dir(path))
			assert.NoError(t, err)
			assert.False(t, expired)
			data, err = os.ReadFile(path)
			assert.NoError(t, err)
			assert.EqualValues(t, cachedData, string(data))

			data, ok = cache.memcache.Get(cacheID)
			assert.False(t, ok)
			assert.Nil(t, data)
			assert.Zero(t, cache.memcache.Size())
		})

		t.Run("Always Fetch Remote On Expired", func(t *testing.T) {
			// should always call fetch remote when expired
			path, expired, err := cache.GetPath(obj, 1111111111, true, fetchRemoteAlwaysOk)
			assert.EqualValues(t, 2, calledOk)
			assert.EqualValues(t, cacheFilenamePrefix+"-00000000001111111111", filepath.Base(path))
			assert.EqualValues(t, cacheDir, filepath.Dir(path))
			assert.NoError(t, err)
			assert.False(t, expired)
			data, err = os.ReadFile(path)
			assert.NoError(t, err)
			assert.EqualValues(t, cachedData, string(data))

			data, ok = cache.memcache.Get(cacheID)
			assert.False(t, ok)
			assert.Nil(t, data)
			assert.Zero(t, cache.memcache.Size())
		})

		t.Run("Use Expired When Fetch Remote Failed", func(t *testing.T) {
			data, expired, err := cache.Get(obj, 1111111112, true, fetchRemoteAlwaysFail)
			assert.ErrorIs(t, err, fetchRemoteError)
			assert.True(t, expired)

			assert.EqualValues(t, 1, calledFail)
			assert.EqualValues(t, cachedData, string(data))

			path, expired, err := cache.GetPath(obj, 1111111112, true, fetchRemoteAlwaysFail)
			assert.ErrorIs(t, err, fetchRemoteError)
			assert.True(t, expired)

			assert.EqualValues(t, 2, calledFail)
			assert.EqualValues(t, cacheFilenamePrefix+"-00000000001111111111", filepath.Base(path))
			assert.EqualValues(t, cacheDir, filepath.Dir(path))
		})
	})

	t.Run("Max Age 100s", func(t *testing.T) {
		defer t.Cleanup(func() {
			calledFail = 0
			calledOk = 0
		})

		cacheDir := t.TempDir()
		cache := NewTwoTierCache(fshelper.NewOSFS(false, func(fshelper.Op, string) (string, error) {
			return cacheDir, nil
		}), -1, -1, 100)

		// use expired
		data, expired, err := cache.Get(obj, 1111111111, true, fetchRemoteAlwaysFail)
		assert.EqualValues(t, 1, calledFail)
		assert.Error(t, err)

		assert.Nil(t, data)
		assert.False(t, expired)

		data, expired, err = cache.Get(obj, 1111111111, true, fetchRemoteAlwaysOk)
		assert.EqualValues(t, 1, calledOk)
		assert.False(t, expired)
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))

		data, ok := cache.memcache.Get(cacheID)
		assert.True(t, ok)
		assert.Equal(t, cachedData, string(data))

		data, expired, err = cache.Get(obj, 1111111111, true, fetchRemoteAlwaysOk)
		assert.EqualValues(t, 1, calledOk)
		assert.False(t, expired)
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))
	})

	t.Run("Max Size Too Small", func(t *testing.T) {
		defer t.Cleanup(func() {
			calledFail = 0
			calledOk = 0
		})

		cacheDir := t.TempDir()
		cache := NewTwoTierCache(fshelper.NewOSFS(false, func(fshelper.Op, string) (string, error) {
			return cacheDir, nil
		}), 1, 1024, 100)
		data, expired, err := cache.Get(obj, 1111111111, true, fetchRemoteAlwaysOk)
		assert.False(t, expired)
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))

		data, ok := cache.memcache.Get(cacheID)
		assert.False(t, ok)
		assert.Nil(t, data)

		data, expired, err = cache.Get(obj, 1111111111, true, fetchRemoteAlwaysOk)
		assert.False(t, expired)
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))
	})

	t.Run("Max Item Size Too Small", func(t *testing.T) {
		defer t.Cleanup(func() {
			calledFail = 0
			calledOk = 0
		})

		cacheDir := t.TempDir()
		cache := NewTwoTierCache(fshelper.NewOSFS(false, func(fshelper.Op, string) (string, error) {
			return cacheDir, nil
		}), 0, 1, 100)
		data, expired, err := cache.Get(obj, 1111111111, true, fetchRemoteAlwaysOk)
		assert.False(t, expired)
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))

		data, ok := cache.memcache.Get(cacheID)
		assert.False(t, ok)
		assert.Nil(t, data)

		data, expired, err = cache.Get(obj, 1111111111, true, fetchRemoteAlwaysOk)
		assert.False(t, expired)
		assert.NoError(t, err)
		assert.EqualValues(t, cachedData, string(data))
	})
}

func TestFormatCacheFilenamePrefix(t *testing.T) {
	for _, test := range []struct {
		id       string
		expected string
	}{
		{"foo", "acbd18db4cc2f85cedef654fccc4a4d8"},
		{"bar", "37b51d194a7513e45b56f6524f2d51f2"},
	} {
		t.Run(test.id, func(t *testing.T) {
			assert.EqualValues(t, test.expected, formatCacheFilenamePrefix(test.id))
		})
	}
}

func TestFormatLocalCacheFilename(t *testing.T) {
	for _, test := range []struct {
		now      int64
		prefix   string
		expected string
	}{
		{123, "a", "a-" + strings.Repeat("0", 17) + "123"},
		{1234567890, "b", "b-" + strings.Repeat("0", 10) + "1234567890"},
	} {
		t.Run(fmt.Sprint(test.now), func(t *testing.T) {
			assert.EqualValues(t, test.expected, formatLocalCacheFilename(test.prefix, "", test.now))
		})
	}
}

func TestStoreLocalCache(t *testing.T) {

	t.Run("Invalid Path", func(t *testing.T) {
		defer t.Cleanup(func() {})

		tmpdir := t.TempDir()

		ofs := fshelper.NewOSFS(true, func(fshelper.Op, string) (string, error) {
			return tmpdir, nil
		})

		size, content, err := storeLocalCache(ofs, "invalid/non-existing", strings.NewReader("NOT USED"), true)
		assert.ErrorIs(t, err, fs.ErrNotExist)
		assert.Nil(t, content)
		assert.Zero(t, size)
	})

	t.Run("Reader Error", func(t *testing.T) {
		defer t.Cleanup(func() {})

		tmpdir := t.TempDir()

		ofs := fshelper.NewOSFS(true, func(fshelper.Op, string) (string, error) {
			return tmpdir, nil
		})

		size, content, err := storeLocalCache(ofs, "test", testhelper.NewAlwaysFailReader(io.ErrClosedPipe), true)
		assert.ErrorIs(t, err, io.ErrClosedPipe)
		assert.Nil(t, content)
		assert.Zero(t, size)
	})

	for _, test := range []struct {
		name       string
		data       string
		retContent bool
	}{
		{"No Content", "test-data", false},
		{"With Content", "test-data", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer t.Cleanup(func() {})

			tmpdir := t.TempDir()

			ofs := fshelper.NewOSFS(true, func(fshelper.Op, string) (string, error) {
				return tmpdir, nil
			})

			size, content, err := storeLocalCache(ofs, "test",
				strings.NewReader(test.data),
				test.retContent,
			)
			if !assert.NoError(t, err) {
				return
			}

			assert.EqualValues(t, len(test.data), size)
			if test.retContent {
				assert.EqualValues(t, test.data, string(content))
			}

			content, err = ofs.ReadFile("test")
			if !assert.NoError(t, err, "failed to read file just wrote") {
				return
			}

			assert.EqualValues(t, test.data, string(content))
		})
	}
}

func TestLookupLocalCache_fs(t *testing.T) {
	for _, test := range []struct {
		name   string
		fs     fs.FS
		expErr error
	}{
		{"Valid ErrNotExist", testhelper.NewAlwaysErrFS(fs.ErrNotExist), nil},
		{"Invalid ErrPerm", testhelper.NewAlwaysErrFS(fs.ErrPermission), fs.ErrPermission},
		{"Valid Empty FS", fstest.MapFS{}, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, _, _, err := lookupLocalCache(test.fs, "", "", 0)
			assert.ErrorIs(t, err, test.expErr)
		})
	}
}

// nolint:revive
func TestLookupLocalCache_entries(t *testing.T) {
	foo_1 := formatLocalCacheFilename("foo", "", 1)
	foo_100 := formatLocalCacheFilename("foo", "", 100)
	foo_200 := formatLocalCacheFilename("foo", "", 200)
	foo_invalid_timestamp := "foo-0x1b"
	fs := fstest.MapFS{
		foo_1:   &fstest.MapFile{},
		foo_100: &fstest.MapFile{},
		foo_200: &fstest.MapFile{},

		foo_invalid_timestamp: &fstest.MapFile{},

		"foo": &fstest.MapFile{},
		"gee": &fstest.MapFile{},
	}

	for _, test := range []struct {
		name      string
		notBefore int64
		prefix    string

		active, expired, invalid []string
	}{
		{
			name:      "Valid No Active",
			notBefore: 201,
			prefix:    "foo",
			active:    nil,
			expired:   []string{foo_1, foo_100, foo_200},
			invalid:   []string{"foo", foo_invalid_timestamp},
		},
		{
			name:      "Valid One Active",
			notBefore: 199,
			prefix:    "foo",
			active:    []string{foo_200},
			expired:   []string{foo_1, foo_100},
			invalid:   []string{"foo", foo_invalid_timestamp},
		},
		{
			name:      "Valid Two Active",
			notBefore: 99,
			prefix:    "foo",
			active:    []string{foo_100, foo_200},
			expired:   []string{foo_1},
			invalid:   []string{"foo", foo_invalid_timestamp},
		},
		{
			name:      "Valid No Cache",
			notBefore: 0,
			prefix:    "bar",
			active:    nil,
			expired:   nil,
			invalid:   nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			active, expired, invalid, err := lookupLocalCache(fs, test.prefix, "", test.notBefore)
			assert.EqualValues(t, test.active, active)
			assert.EqualValues(t, test.expired, expired)
			assert.EqualValues(t, test.invalid, invalid)
			assert.NoError(t, err)
		})
	}
}
