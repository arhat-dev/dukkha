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
			yaml: `foo@a: echo bar`,
			expected: &testFieldStruct{
				BaseField: BaseField{
					unresolvedFields: map[unresolvedFieldKey]*unresolvedFieldValue{
						{
							fieldName: "Foo",
							renderer:  "a",
						}: {
							fieldValue:    reflect.Value{},
							yamlFieldName: "foo",
							rawData:       []string{"echo bar"},
						},
					},
				},
				Foo: "",
			},
		},
		{
			name: "catchAll+renderer",
			yaml: `{other_field_1@a: foo, other_field_2@b: bar}`,
			expected: &testFieldStruct{
				BaseField: BaseField{
					unresolvedFields: map[unresolvedFieldKey]*unresolvedFieldValue{
						{
							fieldName: "Other",
							renderer:  "a",
						}: {
							fieldValue:    reflect.Value{},
							yamlFieldName: "other_field_1",
							rawData:       []string{"foo"},
						},
						{
							fieldName: "Other",
							renderer:  "b",
						}: {
							fieldValue:    reflect.Value{},
							yamlFieldName: "other_field_2",
							rawData:       []string{"bar"},
						},
					},
				},
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
