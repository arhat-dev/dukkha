package field

import (
	"reflect"
	"unsafe"

	"gopkg.in/yaml.v3"
)

type _private struct{}

type Interface interface {
	Type() reflect.Type

	yaml.Unmarshaler

	requireBaseField(_private)
}

func New(f Interface) Interface {
	v := reflect.ValueOf(f)

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.NumField() == 0 {
		return nil
	}

	firstField := v.Field(0)

	switch firstField.Type() {
	case baseFieldStructType:
	default:
		panic("invalid BaseField usage, must be first struct")
	}

	var baseField *BaseField
	switch firstField.Kind() {
	case reflect.Struct:
		baseField = firstField.Addr().Interface().(*BaseField)
	default:
		panic("unexpected non struct")
	}

	baseField._parentType = f.Type()
	baseField._parentValue = reflect.NewAt(
		baseField._parentType,
		unsafe.Pointer(firstField.UnsafeAddr()),
	)

	return f
}

var (
	baseFieldStructType = reflect.TypeOf(BaseField{})
)

type BaseField struct {
	_parentType  reflect.Type  `yaml:"-"`
	_parentValue reflect.Value `yaml:"-"`
}

func (f *BaseField) requireBaseField(_private) {}

func (f *BaseField) UnmarshalYAML(n *yaml.Node) error {
	dataBytes, err := yaml.Marshal(n)
	if err != nil {
		return err
	}

	// TODO
	m := make(map[string]interface{})
	_ = yaml.Unmarshal(dataBytes, &m)

	return nil
}
