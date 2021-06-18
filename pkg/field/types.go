package field

import (
	"bytes"
	"reflect"
	"unsafe"

	"arhat.dev/pkg/log"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/utils"
)

type _private struct{}

type Interface interface {
	Type() reflect.Type

	yaml.Unmarshaler

	requireBaseField(_private)
}

func New(f Interface) Interface {
	structType := f.Type()
	for structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	v := reflect.ValueOf(f)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.NumField() == 0 {
		panic("invalid empty field, BaseField is required")
	}

	firstField := v.Field(0)

	switch firstField.Type() {
	case baseFieldStructType:
	default:
		panic("invalid BaseField usage, must be first embedded struct")
	}

	var baseField *BaseField
	switch firstField.Kind() {
	case reflect.Struct:
		baseField = firstField.Addr().Interface().(*BaseField)
	default:
		panic("unexpected non struct")
	}

	baseField._parentType = structType
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

	log.Log.V("BaseField.UnmarshalYAML",
		log.String("type", f._parentType.String()),
		log.Any("node", n),
	)

	// TODO
	m := make(map[string]interface{})
	_ = utils.UnmarshalStrict(bytes.NewReader(dataBytes), &m)

	return nil
}
