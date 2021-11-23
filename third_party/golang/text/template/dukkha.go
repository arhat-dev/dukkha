package template

import "reflect"

// GetExecFuncs returns all template funcs accessible when executing
// the template
//
// *ATTENTION: THIS METHOD IS NOT PART OF STD IMPLEMENTATION*
func (t *Template) GetExecFuncs() map[string]reflect.Value {
	return t.execFuncs
}
