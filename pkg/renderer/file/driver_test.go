package file

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	tempFile, err := ioutil.TempFile(os.TempDir(), "dukkha-test-*")
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

	t.Run("Valid File Exists", func(t *testing.T) {
		ret, err := d.RenderYaml(nil, tempFilePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedData, string(ret))
	})

	t.Run("Invalid Input Type", func(t *testing.T) {
		_, err := d.RenderYaml(nil, true)
		assert.Error(t, err)
	})

	t.Run("Invalid File Not Exists", func(t *testing.T) {
		_, err := d.RenderYaml(nil, randomData)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
}
