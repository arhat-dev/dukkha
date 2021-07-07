package conf

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBootstrapResolve(t *testing.T) {
	getAbsPath := func(dir string) string {
		p, err := filepath.Abs(dir)
		assert.NoError(t, err, "failed to get absolute path of %q", dir)
		return p
	}

	tests := []struct {
		name string

		config   BootstrapConfig
		expected BootstrapConfig

		expectError bool
	}{
		{
			name:   "Empty",
			config: BootstrapConfig{},
			expected: BootstrapConfig{
				CacheDir:  getAbsPath(".dukkha/cache"),
				ScriptCmd: []string{"sh"},
			},
		},
		{
			name: "Custom Env",
			config: BootstrapConfig{
				Env: []string{
					"FOO=foo",
					"BAR=${FOO}",
				},
			},
			expected: BootstrapConfig{
				Env: []string{
					"FOO=foo",
					"BAR=foo",
				},
				CacheDir:  getAbsPath(".dukkha/cache"),
				ScriptCmd: []string{"sh"},
			},
		},
		{
			name: "Invalid Ref Env Not Set",
			config: BootstrapConfig{
				Env: []string{
					"FOO=${SOME_NON_EXISTING_ENV}",
				},
			},
			expectError: true,
		},
		{
			name: "Invalid Ref Env Is Shell Evaluation",
			config: BootstrapConfig{
				Env: []string{
					"FOO=$(some shell evaluation)",
				},
			},
			expectError: true,
		},
		{
			name: "Custom Script Cmd",
			config: BootstrapConfig{
				ScriptCmd: []string{
					"bash",
				},
			},
			expected: BootstrapConfig{
				CacheDir:  getAbsPath(".dukkha/cache"),
				ScriptCmd: []string{"bash"},
			},
		},
		{
			name: "Invalid Script Cmd With Env Not Set",
			config: BootstrapConfig{
				ScriptCmd: []string{
					"${SOME_ENV_NOT_SET}",
				},
			},
			expectError: true,
		},
		{
			name: "Invalid Script Cmd With Env Is Shell Evaluation",
			config: BootstrapConfig{
				ScriptCmd: []string{
					"$(some shell evaluation)",
				},
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			emptyGlobalEnv := make(map[string]string)
			if test.expectError {
				assert.Error(t, test.config.Resolve(&emptyGlobalEnv))
				return
			}

			if !assert.NoError(t, test.config.Resolve(&emptyGlobalEnv)) {
				return
			}

			assert.EqualValues(t, test.expected, test.config)
		})
	}
}
