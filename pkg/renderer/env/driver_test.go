package env

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
)

func TestNewDriver(t *testing.T) {
	ret := NewDefault()

	assert.NotNil(t, ret)
}

func TestDriver_Render(t *testing.T) {
	cmdPrintHello := []string{"sh", "-c", "printf hello"}
	d := NewDefault().(*driver)
	d.getExecSpec = func(toExec []string, isFilePath bool) (env dukkha.Env, cmd []string, err error) {
		return nil, cmdPrintHello, nil
	}

	rv := dukkha_test.NewTestContext(context.TODO())
	rv.AddEnv(dukkha.EnvEntry{
		Name:  "FOO",
		Value: "bar",
	})

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
			rawData:  `foo $(s)))`,
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
			name:    "Invalid Non-Terminated Shell Evaluation",
			rawData: "some $(non-terminated evaluation ignored",
			errStr:  "without matching ( with )",
		},
		{
			name:    "Invalid Env Not Found",
			rawData: "${NO_SUCH_ENV}",
			errStr:  "unbound variable",
		},
		{
			name:    "Invalid Inner Non-Terminated Shell Evaluation",
			rawData: "some $(inner(not-terminated)",
			errStr:  "must be followed by )",
		},
		{
			name:     "Valid Keep Reference Untouched",
			rawData:  `some \${ESCAPED_DATA}`,
			expected: `some ${ESCAPED_DATA}`,
		},
		{
			name:     "Valid Multiple Keep Reference Untouched",
			rawData:  `some \${ESCAPED_DATA} \${FOO} not expanded`,
			expected: `some ${ESCAPED_DATA} ${FOO} not expanded`,
		},
		{
			name:     "Valid Mixed Reference",
			rawData:  `{{- \$base_image_name := "${FOO}" -}}`,
			expected: `{{- $base_image_name := "bar" -}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := d.RenderYaml(rv, test.rawData)
			if len(test.errStr) != 0 {
				if !assert.Error(t, err, ret) {
					return
				}

				assert.Contains(t, err.Error(), test.errStr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expected, string(ret))
		})

	}
}
