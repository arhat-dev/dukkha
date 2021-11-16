package env

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
)

func TestNewDriver(t *testing.T) {
	ret := NewDefault("")

	assert.NotNil(t, ret)
}

func TestDriver_Render(t *testing.T) {
	vFalse := false
	vTrue := true

	rc := dukkha_test.NewTestContext(context.TODO())
	rc.SetCacheDir(t.TempDir())
	rc.AddEnv(true, &dukkha.EnvEntry{
		Name:  "FOO",
		Value: "bar",
	})

	tests := []struct {
		name       string
		rawData    interface{}
		expected   string
		errStr     string
		enableExec *bool
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
			name:       "Valid Simple Shell Evaluation",
			rawData:    "foo $(echo hello)",
			expected:   "foo hello",
			enableExec: &vTrue,
		},
		{
			name:       "Ignored Shell Evaluation When Exec Disabled Explicitly",
			rawData:    "foo $(echo hello)",
			expected:   "foo $(echo hello)",
			enableExec: &vFalse,
		},
		{
			name:       "Ignored Shell Evaluation When Exec Disabled Implicitly",
			rawData:    "foo $(echo hello)",
			expected:   "foo $(echo hello)",
			enableExec: nil,
		},
		{
			name:       "Always ignored backquoted shell evaluation",
			rawData:    "foo `echo hello` $(echo hello)",
			expected:   "foo `echo hello` hello",
			enableExec: &vTrue,
		},
		{
			name:       "Valid Shell Evaluation With Round Brackets",
			rawData:    `foo $(echo hello)))`,
			expected:   "foo hello))",
			enableExec: &vTrue,
		},
		{
			name:       "Valid Multi Shell Evaluation With Round Brackets",
			rawData:    "foo $(echo hello) $(echo hello)",
			expected:   "foo hello hello",
			enableExec: &vTrue,
		},
		{
			name:       "Valid Multi Embedded Shell Evaluation",
			rawData:    "foo $(echo hello $(echo hello)) $(echo hello)",
			expected:   "foo hello hello hello",
			enableExec: &vTrue,
		},
		{
			name:       "Invalid Non-Terminated Shell Evaluation",
			rawData:    "some $(non-terminated evaluation",
			errStr:     "without matching ( with )",
			enableExec: &vTrue,
		},
		{
			name:    "Invalid Env Not Found",
			rawData: "${NO_SUCH_ENV}",
			errStr:  "unbound variable",
		},
		{
			name:       "Invalid Inner Non-Terminated Shell Evaluation",
			rawData:    "some $(inner(not-terminated)",
			errStr:     "must be followed by )",
			enableExec: &vTrue,
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
			d := NewDefault("").(*Driver)
			d.EnableExec = test.enableExec

			ret, err := d.RenderYaml(rc, test.rawData, nil)
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
