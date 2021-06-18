package field

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Interface = (*testFieldStruct)(nil)

type testFieldStruct struct {
	BaseField

	Foo string `yaml:"foo"`
}

func (f *testFieldStruct) Type() reflect.Type {
	return reflect.TypeOf(*f)
}

// should always panic when passed to New()
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

		getBaseFieldParentValue func(in Interface) reflect.Value

		setDirectFoo          func(in Interface, v string)
		getBaseFieldParentFoo func(in Interface) string
	}{
		{
			name:       "struct",
			targetType: &testFieldStruct{},
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

			assert.Equal(t, test.targetType.Type(), test.getBaseFieldParentValue(foo).Type().Elem())

			if !assert.IsType(t, test.targetType, test.getBaseFieldParentValue(foo).Interface()) {
				return
			}

			test.setDirectFoo(foo, "newValue")
			assert.Equal(t, "newValue", test.getBaseFieldParentFoo(foo))
		})
	}
}
