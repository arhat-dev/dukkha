package shell

import (
	"testing"

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
			name:      "Invalid No GetExecFunc",
			config:    &Config{},
			expectErr: true,
		},
		{
			name: "Valid",
			config: &Config{
				GetExecSpec: func(
					toExec []string, isFilePath bool,
				) (env []string, cmd []string, err error) {
					return
				},
			},
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
