package env

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
)

func TestNewDriver(t *testing.T) {
	ret := New(func(toExec []string, isFilePath bool) (env []string, cmd []string, err error) {
		return
	})

	assert.NotNil(t, ret)
}

func TestDriver_Render(t *testing.T) {
	cmdPrintHello := []string{"sh", "-c", "printf hello"}
	d := New(func(toExec []string, isFilePath bool) (env []string, cmd []string, err error) {
		return nil, cmdPrintHello, nil
	})

	rv := dukkha_test.NewTestContext(context.TODO())
	rv.AddEnv("FOO=bar")

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
			ret, err := d.RenderYaml(rv, test.rawData)
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
