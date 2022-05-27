package parse

import "reflect"

type TemplateFuncs interface {
	Has(name string) bool
	GetByName(name string) reflect.Value
}
