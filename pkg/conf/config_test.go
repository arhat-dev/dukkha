package conf_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/utils"
)

func TestConfigUnmarshal(t *testing.T) {
	type testSpec struct {
		name       string
		configData []byte
	}

	const testDataDir = "testdata"

	entries, err := os.ReadDir(testDataDir)
	if !assert.NoError(t, err, "failed to read test config dir") {
		return
	}

	var tests []*testSpec
	for _, e := range entries {
		data, err := os.ReadFile(filepath.Join(testDataDir, e.Name()))
		if !assert.NoError(t, err, "failed to read test config") {
			return
		}

		tests = append(tests, &testSpec{
			name:       e.Name(),
			configData: data,
		})
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := &conf.Config{}

			err := utils.UnmarshalStrict(bytes.NewReader(test.configData), config)

			assert.NoError(t, err)
		})
	}
}
