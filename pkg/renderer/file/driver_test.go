package file

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/sha256helper"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
)

var _ dukkha.Renderer = (*Driver)(nil)

func TestNewDriver(t *testing.T) {
	assert.NotNil(t, NewDefault(""))
}

func TestDriver_Render(t *testing.T) {
	d := NewDefault("")

	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if !assert.NoError(t, err, "failed to generate random bytes") {
		return
	}

	randomData := hex.EncodeToString(buf)
	expectedData := "Test DUKKHA File Renderer " + randomData

	tempFile, err := os.CreateTemp(os.TempDir(), "dukkha-test-*")
	if !assert.NoError(t, err, "failed to create temp file") {
		return
	}
	tempFilePath := tempFile.Name()
	defer func() {
		assert.NoError(t, os.Remove(tempFilePath), "failed to remove temp file")
	}()

	_, err = tempFile.Write([]byte(expectedData))
	_ = assert.NoError(t, tempFile.Close(), "failed to close temp file")
	if !assert.NoError(t, err, "failed to prepare test data") {
		return
	}

	rc := dukkha_test.NewTestContext(context.TODO())

	t.Run("Valid File Exists", func(t *testing.T) {
		ret, err := d.RenderYaml(rc, tempFilePath, nil)
		assert.NoError(t, err)
		assert.EqualValues(t, []byte(expectedData), ret)
	})

	t.Run("Invalid Input Type", func(t *testing.T) {
		_, err := d.RenderYaml(rc, true, nil)
		assert.Error(t, err)
	})

	t.Run("Invalid File Not Exists", func(t *testing.T) {
		_, err := d.RenderYaml(rc, randomData, nil)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
}

func TestDriver_readFile(t *testing.T) {

}

func TestDriver_cacheData(t *testing.T) {
	defer t.Cleanup(func() {})

	tmpdir := t.TempDir()

	d := NewDefault("").(*Driver)
	d.Init(fshelper.NewOSFS(false, func() (string, error) {
		return tmpdir, nil
	}))

	const testdata = "test-data"

	actual, err := d.cacheData([]byte(testdata))
	assert.NoError(t, err)
	assert.EqualValues(t, filepath.Join(tmpdir, hex.EncodeToString(sha256helper.Sum([]byte(testdata)))), actual)
	data, err := os.ReadFile(string(actual))
	assert.NoError(t, err)
	assert.EqualValues(t, testdata, string(data))
}
