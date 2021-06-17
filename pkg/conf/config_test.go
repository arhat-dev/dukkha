package conf_test

import (
	"os"
	"path/filepath"
	"testing"

	"arhat.dev/dukkha/pkg/conf"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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

			err := yaml.Unmarshal(test.configData, config)

			assert.NoError(t, err)
		})
	}
}
