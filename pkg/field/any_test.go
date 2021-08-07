package field

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func createExpectedJSONValue(t *testing.T, i interface{}) string {
	data, err := json.Marshal(i)
	if !assert.NoError(t, err) {
		t.FailNow()
		return ""
	}

	return string(data)
}

func TestAnyObject(t *testing.T) {
	tests := []struct {
		name  string
		input string

		expectedUnmarshaled *AnyObject
		expectedResolved    *AnyObject

		expectedEquivalent interface{}
	}{
		{
			name:  "Basic Map",
			input: `foo: bar`,

			expectedUnmarshaled: &AnyObject{
				mapData: &mapData{
					Data: map[string]*AnyObject{
						"foo": {scalarData: "bar"},
					},
				},
			},
			expectedResolved: &AnyObject{
				mapData: &mapData{
					Data: map[string]*AnyObject{
						"foo": {scalarData: "bar"},
					},
				},
			},
			expectedEquivalent: map[string]string{
				"foo": "bar",
			},
		},
		{
			name:  "Basic List",
			input: `[foo, bar]`,

			expectedUnmarshaled: &AnyObject{
				arrayData: []*AnyObject{
					{scalarData: "foo"},
					{scalarData: "bar"},
				},
			},
			expectedResolved: &AnyObject{
				arrayData: []*AnyObject{
					{scalarData: "foo"},
					{scalarData: "bar"},
				},
			},
			expectedEquivalent: []string{"foo", "bar"},
		},
		{
			name:  "Simple Rendering Suffix",
			input: `foo@echo: bar`,

			expectedUnmarshaled: &AnyObject{
				mapData: &mapData{
					Data: nil,
				},
			},
			expectedResolved: &AnyObject{
				mapData: &mapData{
					Data: map[string]*AnyObject{
						"foo": {scalarData: "bar"},
					},
				},
			},
			expectedEquivalent: map[string]string{
				"foo": "bar",
			},
		},

		{
			name:  "Combined",
			input: `[{foo@echo: [a,b]}, {bar@echo: [c,d]}]`,

			expectedUnmarshaled: &AnyObject{
				arrayData: []*AnyObject{
					{mapData: &mapData{Data: nil}},
					{mapData: &mapData{Data: nil}},
				},
			},
			expectedResolved: &AnyObject{
				arrayData: []*AnyObject{
					{mapData: &mapData{Data: map[string]*AnyObject{
						"foo": {
							arrayData: []*AnyObject{
								{scalarData: "a"},
								{scalarData: "b"},
							},
						},
					}}},
					{mapData: &mapData{Data: map[string]*AnyObject{
						"bar": {
							arrayData: []*AnyObject{
								{scalarData: "c"},
								{scalarData: "d"},
							},
						},
					}}},
				},
			},
			expectedEquivalent: []interface{}{
				map[string]interface{}{
					"foo": []string{"a", "b"},
				},
				map[string]interface{}{
					"bar": []string{"c", "d"},
				},
			},
		},
	}

	rc := &testRenderingHandler{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			t.Run("unmarshal", func(t *testing.T) {
				obj := &AnyObject{}

				assert.NoError(t, yaml.Unmarshal([]byte(test.input), obj))
				unsetAnyObjectBaseField(obj)
				assert.EqualValues(t, test.expectedUnmarshaled, obj)
			})

			t.Run("resolve+marshal", func(t *testing.T) {
				obj := &AnyObject{}

				assert.NoError(t, yaml.Unmarshal([]byte(test.input), obj))
				assert.NoError(t, obj.ResolveFields(rc, -1))
				unsetAnyObjectBaseField(obj)

				assert.EqualValues(t, test.expectedResolved, obj)

				yamlResult, err := yaml.Marshal(obj)
				if !assert.NoError(t, err) {
					return
				}

				assert.Equal(t, createExpectedYamlValue(t, test.expectedEquivalent), string(yamlResult))

				jsonResult, err := json.Marshal(obj)
				if !assert.NoError(t, err) {
					return
				}

				assert.Equal(t, createExpectedJSONValue(t, test.expectedEquivalent), string(jsonResult))
			})
		})
	}
}

func unsetAnyObjectBaseField(obj *AnyObject) {
	if obj.mapData != nil {
		obj.mapData.BaseField = BaseField{}
		for _, v := range obj.mapData.Data {
			unsetAnyObjectBaseField(v)
		}
	}

	if obj.arrayData != nil {
		for _, v := range obj.arrayData {
			unsetAnyObjectBaseField(v)
		}
	}
}
