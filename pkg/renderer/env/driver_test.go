package env

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/field"
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

func TestDriver_Render(t *testing.T) {
	cmdPrintHello := []string{"sh", "-c", "printf hello"}
	d, err := NewDriver(&Config{
		GetExecSpec: func(
			toExec []string, isFilePath bool,
		) (env []string, cmd []string, err error) {
			return []string{""}, cmdPrintHello, nil
		},
	})
	if !assert.NoError(t, err, "failed to create driver for test") {
		return
	}

	rc := field.WithRenderingValues(context.TODO(), []string{"FOO=bar"})

	tests := []struct {
		name     string
		rawData  interface{}
		expected string
		errStr   string
	}{
		{
			name:     "Valid Plain Text",
			rawData:  "No env ref",
			expected: "No env ref",
		},
		{
			name:     "Valid Plain Data Bytes",
			rawData:  []byte("No env ref"),
			expected: "No env ref",
		},
		{
			name:    "Valid String Sequence",
			rawData: []string{"No env ref", "No env ref"},
			expected: `- No env ref
- No env ref
`,
		},
		{
			name:     "Valid Simple Expansion $FOO",
			rawData:  "foo $FOO",
			expected: "foo bar",
		},
		{
			name:     "Valid Simple Expansion ${FOO}",
			rawData:  "foo ${FOO}",
			expected: "foo bar",
		},
		{
			name:     "Valid Simple Shell Evaluation",
			rawData:  "foo $(some command)",
			expected: "foo hello",
		},
		{
			name:     "Valid Shell Evaluation With Round Brackets",
			rawData:  "foo $(s())))",
			expected: "foo hello))",
		},
		{
			name:     "Valid Multi Shell Evaluation With Round Brackets",
			rawData:  "foo $(say-something() useful) $(say-something() useless)",
			expected: "foo hello hello",
		},
		{
			name:     "Valid Multi Embedded Shell Evaluation With Round Brackets",
			rawData:  "foo $(say-something() $(ok) useful) $(say-something() $(what) useless)",
			expected: "foo hello hello",
		},
		{
			name: "Valid Non-Terminated Shell Evaluation",
			// envhelper.Expand handling func will just ignore it
			rawData:  "some $(non-terminated evaluation ignored",
			expected: "some $(non-terminated evaluation ignored",
		},
		{
			name:    "Invalid Env Not Found",
			rawData: "${NO_SUCH_ENV}",
			errStr:  "not found",
		},
		{
			name:    "Invalid Inner Non-Terminated Shell Evaluation",
			rawData: "some $(inner(not-terminated)",
			errStr:  "non-terminated",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := d.Render(rc, test.rawData)
			if len(test.errStr) != 0 {
				if !assert.Error(t, err) {
					return
				}

				assert.Contains(t, err.Error(), test.errStr)
				return
			}

			assert.Equal(t, test.expected, ret)
		})

	}
}
