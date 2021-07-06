package file

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"arhat.dev/dukkha/pkg/field"
	"github.com/stretchr/testify/assert"
)

func TestNewDriver(t *testing.T) {
	tests := []struct {
		name      string
		config    interface{}
		expectErr bool
	}{
		{
			name:      "Invalid Empty Config",
			config:    nil,
			expectErr: true,
		},
		{
			name:      "Invalid Unexpected Config",
			config:    "foo",
			expectErr: true,
		},
		{
			name:      "Valid",
			config:    &Config{},
			expectErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, err := NewDriver(test.config)

			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, DefaultName, d.Name())
		})
	}
}

func TestDriver_Render(t *testing.T) {
	d, err := NewDriver(&Config{})
	if !assert.NoError(t, err, "failed to create driver for test") {
		return
	}

	buf := make([]byte, 32)
	_, err = rand.Read(buf)
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

	rc := field.WithRenderingValues(context.TODO())

	t.Run("Valid File Exists", func(t *testing.T) {
		ret, err := d.Render(rc, tempFilePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedData, ret)
	})

	t.Run("Invalid Input Type", func(t *testing.T) {
		_, err := d.Render(rc, true)
		assert.Error(t, err)
	})

	t.Run("Invalid File Not Exists", func(t *testing.T) {
		_, err := d.Render(rc, randomData)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
}
