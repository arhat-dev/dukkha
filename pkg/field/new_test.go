package field

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Interface = (*testFieldStruct)(nil)

type testInnerFieldStruct struct {
	BaseField

	Bar string `yaml:"bar"`
}

type testFieldStruct struct {
	BaseField

	Foo   string   `yaml:"foo"`
	Other []string `yaml:"" dukkha:"other"`

	InnerStruct testInnerFieldStruct  `yaml:"innerStruct"`
	InnerPtr    *testInnerFieldStruct `yaml:"innerPtr"`
}

// should always panic when passed to New()
type testFieldPtr struct {
	*BaseField

	Foo string `yaml:"foo"`
}

func TestNewField(t *testing.T) {

	fStruct := &testFieldStruct{}
	fPtr1 := testFieldPtr{}
	fPtr2 := &testFieldPtr{}

	tests := []struct {
		name        string
		targetType  Interface
		panicOnInit bool

		getBaseFieldParentValue func() reflect.Value

		setDirectFoo          func(v string)
		getBaseFieldParentFoo func() string
	}{
		{
			name:       "Ptr BaseField",
			targetType: fStruct,
			getBaseFieldParentValue: func() reflect.Value {
				return fStruct.BaseField._parentValue
			},
			setDirectFoo: func(v string) {
				fStruct.Foo = v
			},
			getBaseFieldParentFoo: func() string {
				return fStruct.BaseField._parentValue.Interface().(*testFieldStruct).Foo
			},
		},
		{
			name:        "Struct *BaseField",
			targetType:  fPtr1,
			panicOnInit: true,
		},
		{
			name:       "Ptr *BaseField",
			targetType: fPtr2,
			getBaseFieldParentValue: func() reflect.Value {
				return fPtr2.BaseField._parentValue
			},
			setDirectFoo: func(v string) {
				fPtr2.Foo = v
			},
			getBaseFieldParentFoo: func() string {
				return fPtr2.BaseField._parentValue.Interface().(*testFieldPtr).Foo
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.panicOnInit {
				func() {
					defer func() {
						assert.NotNil(t, recover())
					}()

					_ = Init(test.targetType)
				}()

				return
			}

			_ = Init(test.targetType)

			if !assert.IsType(t, test.targetType, test.getBaseFieldParentValue().Interface()) {
				return
			}

			test.setDirectFoo("newValue")
			assert.Equal(t, "newValue", test.getBaseFieldParentFoo())
		})
	}
}
