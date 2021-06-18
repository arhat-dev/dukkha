package field

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/utils"
)

var _ Interface = (*testFieldStruct)(nil)

type testFieldStruct struct {
	BaseField

	Foo string `yaml:"foo"`
}

func (f *testFieldStruct) Type() reflect.Type {
	return reflect.TypeOf(*f)
}

// should always panic when passed to NewField()
type testFieldPtr struct {
	*BaseField

	Foo string `yaml:"foo"`
}

func (f testFieldPtr) Type() reflect.Type {
	return reflect.TypeOf(f)
}

func TestNewField(t *testing.T) {
	tests := []struct {
		name       string
		targetType Interface
		willPanic  bool

		getBaseFieldParentType  func(in Interface) reflect.Type
		getBaseFieldParentValue func(in Interface) reflect.Value

		setDirectFoo          func(in Interface, v string)
		getBaseFieldParentFoo func(in Interface) string
	}{
		{
			name:       "struct",
			targetType: &testFieldStruct{},
			getBaseFieldParentType: func(in Interface) reflect.Type {
				return in.(*testFieldStruct).BaseField._parentType
			},
			getBaseFieldParentValue: func(in Interface) reflect.Value {
				return in.(*testFieldStruct).BaseField._parentValue
			},
			setDirectFoo: func(in Interface, v string) {
				in.(*testFieldStruct).Foo = v
			},
			getBaseFieldParentFoo: func(in Interface) string {
				return in.(*testFieldStruct).BaseField._parentValue.Interface().(*testFieldStruct).Foo
			},
		},
		{
			name:       "pointer",
			targetType: testFieldPtr{},
			willPanic:  true,
		},
		{
			name:       "pointer2",
			targetType: &testFieldPtr{},
			willPanic:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.willPanic {
				defer func() {
					assert.NotNil(t, recover())
				}()
			}

			foo := New(test.targetType)

			assert.Equal(t, test.targetType.Type(), test.getBaseFieldParentType(foo))

			if !assert.IsType(t, test.targetType, test.getBaseFieldParentValue(foo).Interface()) {
				return
			}

			test.setDirectFoo(foo, "newValue")
			assert.Equal(t, "newValue", test.getBaseFieldParentFoo(foo))
		})
	}
}

func TestBaseField_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected interface{}
	}{
		{
			name:     "basic",
			yaml:     `foo: bar`,
			expected: &testFieldStruct{Foo: "bar"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := New(&testFieldStruct{}).(*testFieldStruct)

			if !assert.NoError(t, utils.UnmarshalStrict(strings.NewReader(test.yaml), out)) {
				return
			}

			out._parentType = nil
			out._parentValue = reflect.Value{}

			assert.EqualValues(t, test.expected, out)
		})
	}
}
