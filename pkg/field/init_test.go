package field

import (
	"reflect"
	"testing"

	"arhat.dev/dukkha/pkg/types"
	"github.com/stretchr/testify/assert"
)

var _ types.Field = (*testFieldStruct)(nil)

type testFieldStruct struct {
	BaseField

	Foo   string   `yaml:"foo"`
	Other []string `yaml:"" dukkha:"other"`
}

// should always panic when passed to New()
type testFieldPtr struct {
	*BaseField

	Foo string `yaml:"foo"`
}

func TestInit(t *testing.T) {
	fStruct := &testFieldStruct{}
	fPtr1 := testFieldPtr{}
	fPtr2 := &testFieldPtr{}

	tests := []struct {
		name        string
		targetType  types.Field
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
