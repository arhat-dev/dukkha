package field

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var _ RenderingHandler = (*testRenderingHandler)(nil)

type testRenderingHandler struct{}

func (h *testRenderingHandler) RenderYaml(r string, d interface{}) (result []byte, err error) {
	switch r {
	case "echo":
		switch t := d.(type) {
		case string:
			return []byte(t), nil
		case []byte:
			return t, nil
		}

		data, err := yaml.Marshal(d)
		return data, err
	case "err":
		return nil, fmt.Errorf("always error")
	case "empty":
		return []byte{}, nil
	}

	panic("unexpected renderer name")
}

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
				Foo: "bar",
			},
		},
		{
			name: "basic+renderer",
			yaml: `foo@a: echo bar`,
			expected: &testFieldStruct{
				BaseField: BaseField{
					unresolvedFields: map[unresolvedFieldKey]*unresolvedFieldValue{
						{
							yamlKey: "foo",
							suffix:  "a",
						}: {
							fieldName:  "Foo",
							fieldValue: reflect.Value{},
							rawDataList: []*alterInterface{
								{
									scalarData: "echo bar",
								},
							},
							renderers: []string{"a"},
						},
					},
				},
				Foo: "",
			},
		},
		{
			name: "catchAll+renderer",
			yaml: `{other_field_1@a: foo, other_field_2@b: bar }`,
			expected: &testFieldStruct{
				BaseField: BaseField{
					catchOtherFields: map[string]struct{}{
						"other_field_1": {},
						"other_field_2": {},
					},
					catchOtherCache: nil,
					unresolvedFields: map[unresolvedFieldKey]*unresolvedFieldValue{
						{
							yamlKey: "other_field_1",
							suffix:  "a",
						}: {
							fieldName:  "Other",
							fieldValue: reflect.Value{},
							rawDataList: []*alterInterface{
								{
									mapData: map[string]*alterInterface{
										"other_field_1": {scalarData: "foo"},
									},
								},
							},
							renderers:         []string{"a"},
							isCatchOtherField: true,
						},
						{
							yamlKey: "other_field_2",
							suffix:  "b",
						}: {
							fieldName:  "Other",
							fieldValue: reflect.Value{},
							rawDataList: []*alterInterface{
								{
									mapData: map[string]*alterInterface{
										"other_field_2": {scalarData: "bar"},
									},
								},
							},
							renderers:         []string{"b"},
							isCatchOtherField: true,
						},
					},
				},
				// `Other` field should NOT be initialized
				// it will be initialized during resolving
				Other: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := Init(&testFieldStruct{}, nil).(*testFieldStruct)
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

func TestBaseField_UnmarshalYAML_Init(t *testing.T) {
	type Inner struct {
		BaseField

		Foo string `yaml:"foo"`

		DeepInner struct {
			BaseField

			Bar string `yaml:"bar"`
		} `yaml:"deep"`
	}

	rh := &testRenderingHandler{}

	t.Run("struct", func(t *testing.T) {
		type T struct {
			BaseField

			Foo Inner `yaml:"foo"`
		}

		out := Init(&T{}, nil).(*T)

		assert.NoError(t, yaml.Unmarshal([]byte(`foo: { foo: bar }`), out))
		assert.Equal(t, "bar", out.Foo.Foo)
		assert.EqualValues(t, 1, out.Foo.BaseField._initialized)

		out = Init(&T{}, nil).(*T)

		assert.NoError(t, yaml.Unmarshal([]byte(`foo@echo: "{ foo: rendered-bar }"`), out))
		assert.Equal(t, "", out.Foo.Foo)
		assert.Len(t, out.BaseField.unresolvedFields, 1)
		assert.Len(t, out.Foo.BaseField.unresolvedFields, 0)
		assert.EqualValues(t, 1, out.Foo.BaseField._initialized)

		out.ResolveFields(rh, -1)

		assert.EqualValues(t, "rendered-bar", out.Foo.Foo)
	})

	t.Run("struct inline", func(t *testing.T) {
		type T struct {
			BaseField

			Foo Inner `yaml:",inline"`
		}

		out := Init(&T{}, nil).(*T)

		assert.NoError(t, yaml.Unmarshal([]byte(`foo: bar`), out))
		assert.Equal(t, "bar", out.Foo.Foo)
		assert.EqualValues(t, 1, out.Foo.BaseField._initialized)

		out = Init(&T{}, nil).(*T)

		assert.NoError(t, yaml.Unmarshal([]byte(`foo@echo: "{ foo: rendered-bar }"`), out))
		assert.Equal(t, "", out.Foo.Foo)
		assert.EqualValues(t, 1, out.Foo.BaseField._initialized)
		assert.Len(t, out.BaseField.unresolvedFields, 0)
		assert.Len(t, out.Foo.BaseField.unresolvedFields, 1)
	})

	t.Run("struct embedded ", func(t *testing.T) {
		// nolint:unused
		type T struct {
			BaseField

			Inner `yaml:"inner"`
		}

		// TODO
	})

	t.Run("struct embedded inline", func(t *testing.T) {
		// nolint:unused
		type T struct {
			BaseField

			Inner `yaml:",inline"`
		}

		// TODO
	})

	t.Run("ptr", func(t *testing.T) {
		// nolint:unused
		type T struct {
			BaseField

			Foo *Inner `yaml:"foo"`
		}

		// TODO
	})

	t.Run("ptr inline", func(t *testing.T) {
		// nolint:unused
		type T struct {
			BaseField

			Foo *Inner `yaml:",inline"`
		}

		// TODO
	})

	t.Run("ptr embedded ", func(t *testing.T) {
		// nolint:unused
		type T struct {
			BaseField

			*Inner `yaml:"inner"`
		}

		// TODO
	})

	t.Run("ptr embedded inline", func(t *testing.T) {
		// nolint:unused
		type T struct {
			BaseField

			*Inner `yaml:",inline"`
		}

		// TODO
	})
}
