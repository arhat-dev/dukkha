package conf

import (
	"context"
	"io/fs"
	"testing"
	"testing/fstest"

	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
)

func TestRead(t *testing.T) {
	t.Parallel()

	var (
		testIncludeEmptyConfig = newConfig(func(c *Config) {
			c.Include = []*IncludeEntry{
				{
					Path: "empty.yaml",
				},
				{
					Path: "empty-dir",
				},
				{
					// include self is ok, should be ignored
					Path: "include.yaml",
				},
			}

			for _, inc := range c.Include {
				_ = rs.Init(inc, nil)
			}
		})
	)
	tests := []struct {
		name string

		configPaths        []string
		ignoreFileNotExist bool
		expected           *Config
		expectErr          bool
	}{
		{
			name:     "None",
			expected: NewConfig(),
		},

		// missing ok (ignoreFileNotExist=true)
		{
			name:               "Single File Missing OK",
			configPaths:        []string{"config-missing.yaml"},
			ignoreFileNotExist: true,
			expected:           NewConfig(),
		},
		{
			name:               "Single Dir Missing OK",
			configPaths:        []string{"dir-missing"},
			ignoreFileNotExist: true,
			expected:           NewConfig(),
		},
		{
			name:               "Multiple Missing OK",
			configPaths:        []string{"dir-missing", "config-missing.yaml"},
			ignoreFileNotExist: true,
			expected:           NewConfig(),
		},

		// missing not ok (ignoreFileNotExist=false)
		{
			name:               "Single File Missing",
			configPaths:        []string{"config-missing.yaml"},
			ignoreFileNotExist: false,
			expectErr:          true,
		},
		{
			name:               "Single Dir Missing OK",
			configPaths:        []string{"dir-missing"},
			ignoreFileNotExist: false,
			expectErr:          true,
		},
		{
			name:               "Multiple Missing OK",
			configPaths:        []string{"dir-missing", "config-missing.yaml"},
			ignoreFileNotExist: false,
			expectErr:          true,
		},

		// empty
		{
			name:               "Empty Single File",
			configPaths:        []string{"empty.yaml"},
			ignoreFileNotExist: false,
			expected:           NewConfig(),
		},
		{
			name:               "Empty Single Dir",
			configPaths:        []string{"empty-dir"},
			ignoreFileNotExist: false,
			expected:           NewConfig(),
		},
		{
			name:               "Empty Multiple Source",
			configPaths:        []string{"empty-dir", "empty.yaml"},
			ignoreFileNotExist: false,
			expected:           testIncludeEmptyConfig,
		},
	}

	testFS := fstest.MapFS{
		"empty.yaml": &fstest.MapFile{
			Data: nil,
		},
		"empty-dir": &fstest.MapFile{
			Data: nil,
			Mode: fs.ModeDir,
		},
		"include-empty.yaml": &fstest.MapFile{
			Data: configBytes(t, testIncludeEmptyConfig),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			visitedPaths := make(map[string]struct{})
			mergedConfig := NewConfig()

			rc := dukkha_test.NewTestContext(context.TODO(), t.TempDir())
			err := Read(
				rc,
				testFS,
				test.configPaths,
				test.ignoreFileNotExist,
				&visitedPaths,
				mergedConfig,
			)
			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expectedBytes, err := yaml.Marshal(test.expected)
			assert.NoError(t, err, "failed to marshal expected config")

			// field include is not merged
			mergedConfig.Include = test.expected.Include
			actualBytes, err := yaml.Marshal(mergedConfig)
			assert.NoError(t, err, "failed to marshal config loaded")

			assert.EqualValues(t, string(expectedBytes), string(actualBytes))
			t.Log(string(actualBytes))
		})
	}
}

func newConfig(update func(c *Config)) *Config {
	ret := NewConfig()
	if update != nil {
		update(ret)
	}

	return ret
}

func configBytes(t *testing.T, c *Config) []byte {
	data, err := yaml.Marshal(c)
	assert.NoError(t, err, "")
	return data
}
