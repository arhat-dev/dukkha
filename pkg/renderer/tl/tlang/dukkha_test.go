package tlang

import "reflect"

func (fm FuncMap) Has(name string) bool {
	return fm[name] == nil
}

func (fm FuncMap) GetByName(name string) reflect.Value {
	ref := fm[name]
	return reflect.ValueOf(ref)
}
