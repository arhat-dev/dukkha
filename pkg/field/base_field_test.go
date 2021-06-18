package field

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestBaseField_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected interface{}
	}{
		{
			name: "basic",
			yaml: `foo: bar`,
			expected: &testFieldStruct{
				BaseField: BaseField{
					unresolvedFields: nil,
				},
				// TODO: add back after supported
				// Foo: "bar",
			},
		},
		{
			name: "basic+renderer",
			yaml: `foo@hi: echo bar`,
			expected: &testFieldStruct{
				BaseField: BaseField{
					unresolvedFields: map[unresolvedFieldKey]*unresolvedFieldValue{
						{
							fieldName: "Foo",
						}: {
							fieldValue: reflect.Value{},
							renderer:   "hi",
							rawData:    "echo bar",
						},
					},
				},
				Foo: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := New(&testFieldStruct{}).(*testFieldStruct)
			assert.EqualValues(t, 1, out._initialized)

			if !assert.NoError(t, yaml.Unmarshal([]byte(test.yaml), out)) {
				return
			}

			out._initialized = 0
			out._parentValue = reflect.Value{}
			for k := range out.unresolvedFields {
				out.unresolvedFields[k].fieldValue = reflect.Value{}
			}

			assert.EqualValues(t, test.expected, out)
		})
	}
}
